package clients

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
	"net/url"
	"strings"
)

type ClientNetwork interface {
	Send(payload string) (interface{}, error)
}

func BuildNetworkClient(endpoint string, settings *datamover_commons.DataMoverNetworkSettings) (ClientNetwork, error) {
	if len(endpoint) > 0 {
		if strings.HasPrefix(endpoint, "nio:") {
			uri, err := url.Parse(endpoint)
			if nil != err {
				return nil, err
			}
			host := uri.Hostname()
			port := gg.Convert.ToInt(uri.Port())
			return NewClientNio(host, port, settings)
		} else {
			return NewClientHttp(endpoint, settings)
		}
	}
	return nil, gg.Errors.Prefix(datamover_commons.PanicSystemError, fmt.Sprintf("Protocol Not Supported '%s': ", endpoint))
}
