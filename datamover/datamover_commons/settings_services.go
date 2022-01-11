package datamover_commons

type SettingsNet struct {
	Enabled       bool                      `json:"enabled"`
	Authorization *SettingsNetAuthorization `json:"authorization"`
	Services      []*SettingsNetService     `json:"services"`
}

type SettingsNetAuthorization struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SettingsNetService struct {
	Enabled               bool                   `json:"enabled"`
	Name                  string                 `json:"name"`
	Protocol              string                 `json:"protocol"`
	ProtocolConfiguration map[string]interface{} `json:"protocol_configuration"` // each protocol has its own configuration
}

type SettingsNetProtocolNio struct {
	Port int `json:"port"`
}
