package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	addressFlag    = "address"
	verboseFlag    = "verbose"
	defaultAddress = ":9090"
	routesFileFlag = "routes-file"
)

const (
	usageHeader = `
skrop – Skipper based media service using the vips library.

https://github.com/zalando-incubator/skrop

`
	addressUsage    = "network address that skoap should listen on"
	verboseUsage    = "enable verbose logging"
	routesFileUsage = `alternatively to the target address, it is possible to use a full
	eskip route configuration, and specify the auth() and authTeam()
	filters for the routes individually.
	See also: https://godoc.org/github.com/zalando/skipper/eskip`
)

var fs *flag.FlagSet

var (
	address    string
	verbose    bool
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
	fs.BoolVar(&verbose, verboseFlag, false, verboseUsage)
	fs.StringVar(&routesFile, routesFileFlag, "", routesFileUsage)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(-1)
	}
}

func main() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if routesFile == "" {
		logUsage("A routes file needs to be specified.")
	}
	log.Debug(fmt.Sprintf("Using routes-file %s", routesFile))
}
