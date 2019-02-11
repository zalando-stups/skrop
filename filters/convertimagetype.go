package filters

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
	"strings"
)

// ConvertImageType is the name of the filter
const ConvertImageType = "convertImageType"

type convertImageType struct {
	imageType bimg.ImageType
}

// NewConvertImageType creates a new filter of this type
func NewConvertImageType() filters.Spec {
	return &convertImageType{}
}

func (f *convertImageType) Name() string {
	return ConvertImageType
}

func (f *convertImageType) CreateOptions(_ *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for convert image type", f)

	return &bimg.Options{
		Type: f.imageType,
	}, nil
}

func (f *convertImageType) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	var zero bimg.ImageType

	//it can be merged if the background was not set (in options or in self) or if they are set to the same value
	return other.Type == zero || other.Type == self.Type
}

func (f *convertImageType) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Type = self.Type
	return other
}

func (f *convertImageType) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error
	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &convertImageType{}

	imgType, err := parse.EskipStringArg(args[0])

	if err != nil || !bimg.IsTypeNameSupported(imgType) {
		return nil, filters.ErrInvalidFilterParameters
	}

	for ImageType, value := range bimg.ImageTypes {
		if value == imgType {
			c.imageType = ImageType
			break
		}
	}

	return c, err
}

func (f *convertImageType) Request(ctx filters.FilterContext) {}

func (f *convertImageType) Response(ctx filters.FilterContext) {

	err := HandleImageResponse(ctx, f)

	if err != nil {
		return
	}

	resp := ctx.Response()

	fileType := bimg.ImageTypeName(f.imageType)

	contentType := fmt.Sprintf("image/%s", fileType)
	contentDisp := fmt.Sprintf("inline;filename=%s.%s", extractFileName(ctx), fileType)

	resp.Header.Set("Content-Type", contentType)
	resp.Header.Set("Content-Disposition", contentDisp)

}

func extractFileName(ctx filters.FilterContext) string {
	var fileName string
	uriParts := strings.Split(ctx.Request().RequestURI, "/")

	if len(uriParts) > 0 {
		fileName = uriParts[len(uriParts)-1]
		if len(strings.Split(fileName, ".")) == 2 {
			fileName = strings.Split(fileName, ".")[0]
		}
	}

	return fileName
}
