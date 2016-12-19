package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"flag"
	"os"
)

const (
	addressFlag = "address"
	defaultAddress = ":9090"
	routesFileFlag = "routes-file"
)

const (
	usageHeader = `
skrop - Skipper based media service based on the vips library.

https://github.com/zalando-incubator/skrop
`
	addressUsage = `network address that skoap should listen on`

	routesFileUsage = `alternatively to the target address, it is possible to use a full eskip route
configuration, and specify the auth() and authTeam() filters for the routes individually. See also:
https://godoc.org/github.com/zalando/skipper/eskip`
)

var fs *flag.FlagSet

var (
	address string
	verbose bool
	routesFile string
)

func usage() {
	fmt.Fprint(os.Stderr, usageHeader)
	fs.PrintDefaults()
}

func logUsage(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(-1)
}

func init() {
	fs = flag.NewFlagSet("flags", flag.ContinueOnError)
	fs.Usage = usage

	fs.StringVar(&address, addressFlag, defaultAddress, addressUsage)
	fs.StringVar(&routesFile, routesFileFlag, "", routesFileUsage)
}

func main() {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	if routesFile == "" {
		logUsage("a routes file needs to be specified")
	}
}
