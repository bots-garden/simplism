package registry

import "time"

type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Created time.Time `json:"created"`
	Size    int64     `json:"size"`
	//FileType    string     `json:"fileType"`
}
