package download

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
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
				Value:       ".",
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

			tempFilePath, filename, err := doDownload(ctx, u, apiKey)
			if err != nil {
				return fmt.Errorf("failed to download: %w", err)
			}

			outputFilePath := filepath.Join(output, filename)
			err = os.Rename(tempFilePath, outputFilePath)
			if err != nil {
				return fmt.Errorf("unable to move file to the destination: %w", err)
			}

			fmt.Printf("File saved to: %s\n", outputFilePath)
			return nil
		},
	}
}

func doDownload(ctx context.Context, u *url.URL, apiKey string) (tempFilePath string, filename string, err error) {
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

	f, err := os.CreateTemp("", "")

	if err != nil {
		err = fmt.Errorf("failed to create output file: %w", err)
		return
	}

	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	tempFilePath = f.Name()

	return
}
