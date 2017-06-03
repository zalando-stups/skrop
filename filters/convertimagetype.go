package filters

import (
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/zalando-incubator/skrop/parse"
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
	ctx.Response().Header.Set("Content-Type", "image/"+bimg.ImageTypeName(c.imageType))
	HandleImageResponse(ctx, c)
}