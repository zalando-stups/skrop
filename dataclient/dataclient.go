package dataclient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando-incubator/skrop/filters"
	"github.com/zalando/skipper/eskip"
	"github.com/zalando/skipper/eskipfile"
	"github.com/zalando/skipper/routing"
)

type skropDataClient struct {
	fileName string
	prepend  *eskip.Filter
	append   *eskip.Filter
}

func NewSkropDataClient(eskipFile string) routing.DataClient {

	emptyArgs := make([]interface{}, 0)

	pre := &eskip.Filter{
		Name: filters.NewFinalizeResponse().Name(),
		Args: emptyArgs,
	}
	post := &eskip.Filter{
		Name: filters.NewSetupResponse().Name(),
		Args: emptyArgs,
	}
	return skropDataClient{
		fileName: eskipFile,
		prepend:  pre,
		append:   post,
	}
}

func (s skropDataClient) LoadAll() ([]*eskip.Route, error) {
	f, err := eskipfile.Open(s.fileName)
	if err != nil {
		log.Error("Error while opening eskip file", err)
		return nil, err
	}

	routes, err := f.LoadAll()
	if err != nil {
		log.Error("Error while loading eskip routes", err)
		return nil, err
	}

	for _, route := range routes {
		if route.BackendType != eskip.ShuntBackend {
			route.Filters = append(append([]*eskip.Filter{s.prepend}, route.Filters...), s.append)
		}
	}

	return routes, nil
}

func (s skropDataClient) LoadUpdate() ([]*eskip.Route, []string, error) {
	return nil, nil, nil
}
