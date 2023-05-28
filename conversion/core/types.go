package core

type PipelineOptions struct {
	FileID     string `json:"fileId"`
	SnapshotID string `json:"snapshotId"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
}

type PipelineResponse struct {
	Options   PipelineOptions `json:"options,omitempty"`
	Original  *S3Object       `json:"original,omitempty"`
	Preview   *S3Object       `json:"preview,omitempty"`
	Text      *S3Object       `json:"text,omitempty"`
	OCR       *S3Object       `json:"ocr,omitempty"`
	Thumbnail *Thumbnail      `json:"thumbnail,omitempty"`
}

type S3Object struct {
	Bucket string      `json:"bucket"`
	Key    string      `json:"key"`
	Size   int64       `json:"size"`
	Image  *ImageProps `json:"image,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Thumbnail struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
