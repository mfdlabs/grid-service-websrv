package flags

import (
	"flag"
)

var (
	// BindAddressIpv4 is the address to bind the prometheus metrics server to.
	BindAddressIpv4 = flag.String("bind-address-ipv4", ":8080", "Address to bind the HTTP server to.")

	// HelpFlag prints the usage.
	HelpFlag = flag.Bool("help", false, "Print usage.")
)

const FlagsUsageString string = `
	[-h|--help] [--bind-address-ipv4=[:8080]]`
