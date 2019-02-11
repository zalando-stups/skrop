package dataclient

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando/skipper/eskip"
	"github.com/zalando/skipper/eskipfile"
	"github.com/zalando/skipper/routing"
)

type skropDataClient struct {
	fileName string
	prepend  *eskip.Filter
}

// NewSkropDataClient creates a new dataclient specific for skrop
func NewSkropDataClient(eskipFile string) routing.DataClient {

	return skropDataClient{
		fileName: eskipFile,
	}
}

func (s skropDataClient) LoadAll() ([]*eskip.Route, error) {
	f, err := eskipfile.Open(s.fileName)
	if err != nil {
		log.Error("error while opening eskip file", err)
		return nil, err
	}

	routes, err := f.LoadAll()
	if err != nil {
		log.Error("error while loading eskip routes", err)
		return nil, err
	}

	if s.prepend == nil {
		return routes, nil
	}

	for _, route := range routes {
		route.Filters = append([]*eskip.Filter{s.prepend}, route.Filters...)
	}

	return routes, nil
}

func (s skropDataClient) LoadUpdate() ([]*eskip.Route, []string, error) {
	return nil, nil, nil
}
