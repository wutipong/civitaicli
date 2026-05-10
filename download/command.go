package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/civitaicli/cache"
)

func Command() *cli.Command {
	urlStr := "https://civitai.com/api/download/models/12345"
	output := "."

	return &cli.Command{
		Name:  "download",
		Usage: "Download a file from Civitai",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "url",
				Config: cli.StringConfig{
					TrimSpace: true,
				},
				Destination: &urlStr,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "Output file path",
				Value:       "",
				Destination: &output,
				Config: cli.StringConfig{
					TrimSpace: true,
				},
			},
		},
		Action: func(ctx context.Context, args *cli.Command) error {
			apiKey := os.Getenv("CIVITAI_API_KEY")
			if apiKey == "" {
				return fmt.Errorf("CIVITAI_API_KEY environment variable is not set")
			}

			u, err := url.Parse(urlStr)
			if err != nil {
				return fmt.Errorf("url parse failed: %w", err)
			}
			fmt.Printf("Downloading from URL: %s\n", u.String())

			cacheFilePath, err := doDownload(ctx, u, apiKey)
			if err != nil {
				return fmt.Errorf("failed to download: %w", err)
			}

			fmt.Printf("cached file location: %s\n", cacheFilePath)
			if output == "" {

				return nil
			}

			filename := filepath.Base(cacheFilePath)
			outputFilePath := filepath.Join(output, filename)

			err = CopyFile(outputFilePath, cacheFilePath)
			if err != nil {
				return fmt.Errorf("unable to copy file to the destination: %w", err)
			}

			fmt.Printf("File saved to: %s\n", outputFilePath)

			return nil
		},
	}
}

func doDownload(ctx context.Context, u *url.URL, apiKey string) (filePath string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("download failed: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed with status: %s", resp.Status)
		return
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		err = fmt.Errorf("failed to get content disposition")
		return
	}

	disposition, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		err = fmt.Errorf("failed to parse content disposition: %w", err)
		return
	}
	if disposition != "attachment" {
		err = fmt.Errorf("unexpected content disposition: %s", disposition)
		return
	}

	filename, ok := params["filename"]
	if !ok {
		err = fmt.Errorf("filename not found in content disposition")
		return
	}

	cacheDir, err := cache.EnsureCacheLocation()
	if err != nil {
		err = fmt.Errorf("unable to get cache location:%w", err)
		return
	}

	filePath = filepath.Join(cacheDir, u.EscapedPath(), filename)
	stat, err := os.Stat(filePath)
	if !errors.Is(err, os.ErrNotExist) && resp.ContentLength == stat.Size() {
		return
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		err = fmt.Errorf("failed to create output file: %w", err)
		return
	}
	defer func() {
		_ = os.Remove(f.Name())
	}()
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	err = os.MkdirAll(filepath.Dir(filePath), 0750)
	if err != nil {
		err = fmt.Errorf("unable to create directory for cached content: %w", err)
		return
	}

	err = os.Rename(f.Name(), filePath)
	if err != nil {
		err = fmt.Errorf("unable to move temp file to cache: %w", err)
		return
	}

	return
}

func CopyFile(dstPath, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("unable to open source file: %w", err)
	}

	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}

	defer dst.Close()

	statSrc, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("unable to get source file stat: %w", err)
	}

	statDst, err := os.Stat(dstPath)
	if err == nil {
		if statSrc.Size() == statDst.Size() {
			fmt.Println("destination file exists. skipped.")
			return nil
		}
	}

	bar := progressbar.DefaultBytes(
		statSrc.Size(),
		"copying",
	)

	_, err = io.Copy(io.MultiWriter(dst, bar), src)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}
