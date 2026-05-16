package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ModelVersion struct {
	ID                int           `json:"id"`
	ModelID           int           `json:"modelId"`
	Name              string        `json:"name"`
	Description       interface{}   `json:"description"`
	BaseModel         string        `json:"baseModel"`
	BaseModelType     string        `json:"baseModelType"`
	Air               string        `json:"air"`
	Status            string        `json:"status"`
	Availability      string        `json:"availability"`
	NsfwLevel         int           `json:"nsfwLevel"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
	PublishedAt       time.Time     `json:"publishedAt"`
	UploadType        string        `json:"uploadType"`
	UsageControl      string        `json:"usageControl"`
	TrainedWords      []interface{} `json:"trainedWords"`
	EarlyAccessConfig interface{}   `json:"earlyAccessConfig"`
	EarlyAccessEndsAt interface{}   `json:"earlyAccessEndsAt"`
	TrainingStatus    interface{}   `json:"trainingStatus"`
	TrainingDetails   interface{}   `json:"trainingDetails"`
	Stats             struct {
		DownloadCount int `json:"downloadCount"`
		ThumbsUpCount int `json:"thumbsUpCount"`
	} `json:"stats"`
	Model struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Nsfw bool   `json:"nsfw"`
		Poi  bool   `json:"poi"`
	} `json:"model"`
	Files       []ModelFile   `json:"files"`
	Images      []interface{} `json:"images"`
	DownloadURL string        `json:"downloadUrl"`
}

type ModelFile struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	SizeKB   float64 `json:"sizeKB"`
	Metadata struct {
		Format string `json:"format"`
		Size   string `json:"size"`
		Fp     string `json:"fp"`
	} `json:"metadata"`
	PickleScanResult string `json:"pickleScanResult"`
	VirusScanResult  string `json:"virusScanResult"`
	Hashes           struct {
		AutoV1 string `json:"AutoV1"`
		AutoV2 string `json:"AutoV2"`
		SHA256 string `json:"SHA256"`
		CRC32  string `json:"CRC32"`
		BLAKE3 string `json:"BLAKE3"`
		AutoV3 string `json:"AutoV3"`
	} `json:"hashes"`
	DownloadURL string `json:"downloadUrl"`
	Primary     bool   `json:"primary"`
}

func GetModelVersionInfo(ctx context.Context, id string, apiKey string) (m ModelVersion, err error) {
	u, err := url.JoinPath(BASE_URL, "/model-versions/", id)
	if err != nil {
		err = fmt.Errorf("invalid url: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("unable to read respose body: %w", err)
		return
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		err = fmt.Errorf("unable to parse json data: %w", err)
		return
	}

	return

}
