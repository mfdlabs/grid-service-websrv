package clientsettings

type importApplicationSettingsRequest struct {
	ApplicationName     string                 `json:"applicationName"`
	ApplicationSettings map[string]interface{} `json:"applicationSettings"`
	Dependendies        []string               `json:"dependencies"`
}

type getApplicationSettingsResponse struct {
	ApplicationSettings map[string]interface{} `json:"applicationSettings"`
}
