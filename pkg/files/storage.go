package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

func SimpanUpload(dir string, f *multipart.FileHeader) (string, error) {
	src, err := f.Open(); if err != nil { return "", err }
	defer src.Close()

	path := fmt.Sprintf("%s/%d_%s", dir, time.Now().Unix(), f.Filename)
	dst, err := os.Create(path); if err != nil { return "", err }
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil { return "", err }
	return path, nil
}
