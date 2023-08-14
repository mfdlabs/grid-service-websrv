package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"

	"github.com/mfdlabs/grid-service-websrv/flags"
	"github.com/mfdlabs/grid-service-websrv/http"
)

var applicationName string
var buildMode string
var commitSha string

// Pre-setup, runs before main.
func init() {
	flags.SetupFlags(applicationName, buildMode, commitSha)
}

// Main entrypoint.
func main() {
	defer glog.Flush()

	if *flags.HelpFlag {
		flag.Usage()

		return
	}

	go http.Start()

	// Wait for a signal to quit
	s := make(chan os.Signal, 1)

	// We want to catch ALL signals to quit
	signal.Notify(s, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		sig := <-s

		glog.Warningf("Received signal %s, exiting\n", sig)

		os.Exit(0)
	}()
}
