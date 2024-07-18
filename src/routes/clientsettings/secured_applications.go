package clientsettings

import (
	"strings"

	"github.com/mfdlabs/grid-service-websrv/flags"
)

func isApplicationSecured(applicationName string) bool {
	securedApplications := strings.Split(*flags.ClientSettingsSecuredSettingsList, ",")

	if len(securedApplications) == 1 && securedApplications[0] == "" {
		return true
	}

	if applicationName == "" {
		return false
	}

	for _, securedApplication := range securedApplications {
		if applicationName == securedApplication {
			return true
		}
	}

	return false
}
