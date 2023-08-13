package flags

func applyEnvironmentVariableFlags() {
	getEnvironmentVariableOrFlag("BIND_ADDRESS_IPv4", BindAddressIpv4)
}
