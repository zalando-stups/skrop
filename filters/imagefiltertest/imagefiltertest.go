package imagefiltertest

import "gopkg.in/h2non/bimg.v1"

type FakeImageFilter bimg.Options

func (h *FakeImageFilter) CreateOptions() *bimg.Options {
	options := bimg.Options(*h)
	return &options
}
