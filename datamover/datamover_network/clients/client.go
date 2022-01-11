package clients

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
	"net/url"
)

type ClientNetwork interface {
	Send(payload string) (interface{}, error)
}

func BuildNetworkClient(uri *url.URL, auth *datamover_commons.SettingsNetAuthorization) (ClientNetwork, error) {
	protocol := uri.Scheme
	switch protocol {
	case "nio":
		host := uri.Hostname()
		port := gg.Convert.ToInt(uri.Port())
		return NewClientNio(host, port, auth)
	}
	return nil, gg.Errors.Prefix(datamover_commons.PanicSystemError, fmt.Sprintf("Protocol Not Supported '%s': ", protocol))
}
