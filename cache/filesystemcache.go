package cache

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type fileSystemCache struct {
}

func NewFileSystemCache() Cache {
	return &fileSystemCache{}
}

func (f *fileSystemCache) Write(cacheKey string, content *CacheContent) error {
	log.Infof("writing key:%q\n", cacheKey)

	err := createDirectoriesForKey(cacheKey)

	if err != nil {
		return err
	}

	err = writeContentToFile(cacheKey, content)

	if err != nil {
		return err
	}

	return nil
}

func writeContentToFile(filePath string, content *CacheContent) error {

	return ioutil.WriteFile(filePath, content.Content, 0644)
}

func (f *fileSystemCache) Read(cacheKey string) (*CacheContent, error) {

	fileContent, err := ioutil.ReadFile(cacheKey)

	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(cacheKey)

	if err != nil {
		return nil, err
	}

	contentType, err := getFileContentType(cacheKey)

	cacheContent := &CacheContent{
		Content:      fileContent,
		LastModified: fileInfo.ModTime(),
		ContentType:  contentType}

	return cacheContent, nil
}

func createDirectoriesForKey(cacheKey string) error {

	dirPath := filepath.Dir(cacheKey)
	return os.MkdirAll(dirPath, os.ModePerm)
}

func getFileContentType(fileName string) (string, error) {

	out, err := os.Open(fileName)
	defer out.Close()

	if err != nil {
		return "", err
	}

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
