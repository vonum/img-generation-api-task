package image

type Response struct {
  Error   string `json:"error,omitempty"`
  ImageID string `json:"image_id,omitempty"`
}
