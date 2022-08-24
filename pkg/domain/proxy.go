package domain

import "strings"

type ProxyConfiguration struct {
	HTTPProxy  string
	HTTPSProxy string

	/* no_proxy is a comma- or space-separated list of machine
	   or domain names, with optional :port part.  If no :port
	   part is present, it applies to all ports on that domain.
	*/
	NoProxy string
}

func (p ProxyConfiguration) GetURLForEndpoint(endpoint string) string {
	if !strings.Contains(p.NoProxy, endpoint) {
		if len(p.HTTPSProxy) > 0 {
			return p.HTTPSProxy
		} else if len(p.HTTPProxy) > 0 {
			return p.HTTPProxy
		}
	}
	return ""
}
