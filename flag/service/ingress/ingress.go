package ingress

type Ingress struct {
	ExternalDNS ExternalDNS
}

type ExternalDNS struct {
	Enabled string
}
