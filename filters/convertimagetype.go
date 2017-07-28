package filters

import (
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/zalando-incubator/skrop/parse"
	"strings"
	"fmt"
)

const ConvertImageType  = "convertImageType"

type convertImageType struct{
	imageType bimg.ImageType
}

func NewConvertImageType() filters.Spec {
	return &convertImageType{}
}

func (c *convertImageType) Name() string {
	return ConvertImageType
}

func (c *convertImageType) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for convert image type", c)

	return &bimg.Options{
		Type: c.imageType,
	}, nil
}

func (r *convertImageType) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	zero := bimg.GaussianBlur{}

	//it can be merged if the background was not set (in options or in self) or if they are set to the same value
	return other.GaussianBlur == zero || self.GaussianBlur == zero || other.Background == self.Background
}

func (r *convertImageType) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Background = self.Background
	return other
}

func (c *convertImageType) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error
	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &convertImageType{}

	imgType, err := parse.EskipStringArg(args[0]);

	if err != nil || !bimg.IsTypeNameSupported(imgType) {
		return nil, filters.ErrInvalidFilterParameters
	}

	for ImageType, value := range bimg.ImageTypes {
		if value == imgType {
			f.imageType = ImageType
			break
		}
	}

	return f, err
}


func (c *convertImageType) Request(ctx filters.FilterContext) {}

func (c *convertImageType) Response(ctx filters.FilterContext) {

	HandleImageResponse(ctx, c)

	resp := ctx.Response()

	fileType := bimg.ImageTypeName(c.imageType)

	contentType := fmt.Sprintf("image/%s",fileType)
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