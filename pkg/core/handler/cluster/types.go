package cluster

type UploadData struct {
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	Content  string `json:"content"`
}
