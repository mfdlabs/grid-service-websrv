package clientsettings

type importApplicationSettingsRequest struct {
	ApplicationName                    string                 `json:"applicationName"`
	ApplicationSettings                map[string]interface{} `json:"applicationSettings"`
	Dependendies                       []string               `json:"dependencies"`
	Reference                          *string                `json:"reference"`
	IsAllowedFromClientSettingsService *bool                  `json:"isAllowedFromClientSettingsService"`
}

type getApplicationSettingsResponse struct {
	ApplicationSettings map[string]interface{} `json:"applicationSettings"`
}

type getApplicationSettingResponse struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type setApplicationSettingRequest struct {
	ApplicationName string      `json:"applicationName"`
	SettingName     string      `json:"settingName"`
	Value           interface{} `json:"value"`
}

type setApplicationSettingResponse struct {
	DidUpdate bool                           `json:"didUpdate"`
	OldValue  *getApplicationSettingResponse `json:"oldValue"`
	NewValue  *getApplicationSettingResponse `json:"newValue"`
}
