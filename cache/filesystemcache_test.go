package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	jpegImageFile = "../images/lisbon-tram.jpg"
	pngImageFile  = "../images/star.png"
)

func TestFileSystemCache_Write(t *testing.T) {
	cache := NewFileSystemCache()
	cont := &CacheContent{
		Content: []byte("fakeImage"),
	}
	err := cache.Write("./tmpimages/image.jpg", cont)

	assert.Nil(t, err)
	assert.DirExists(t, "./tmpimages")
	assert.FileExists(t, "./tmpimages/image.jpg")

	os.Remove("./tmpimages/image.jpg")
	os.RemoveAll("./tmpimages")
}

func TestFileSystemCache_createDirectoriesForKey(t *testing.T) {

	err := createDirectoriesForKey("./hello/kitty/image.jpg")

	assert.Nil(t, err)
	assert.DirExists(t, "./hello/kitty")

	os.RemoveAll("./hello")
}

func TestFileSystemCache_writeContentToFile(t *testing.T) {

	cont := &CacheContent{
		Content: []byte("fakeImage"),
	}

	err := writeContentToFile("./image.jpg", cont)

	assert.Nil(t, err)
	assert.FileExists(t, "./image.jpg")

	os.Remove("./image.jpg")
}

func TestFileSystemCache_getFileContentType(t *testing.T) {

	contentType, _ := getFileContentType(jpegImageFile)

	assert.EqualValues(t, "image/jpeg", contentType)


	contentType, _ = getFileContentType(pngImageFile)

	assert.EqualValues(t, "image/png", contentType)
}
