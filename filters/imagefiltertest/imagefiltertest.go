package imagefiltertest

import "gopkg.in/h2non/bimg.v1"

type FakeImageFilter bimg.Options

func (h *FakeImageFilter) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	options := bimg.Options(*h)
	return &options, nil
}
