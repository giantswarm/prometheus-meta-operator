package resource

type Security struct {
	LetsEncryptEnabled bool
	Whitelisting       SecurityWhitelisting
}

type SecurityWhitelisting struct {
	Enabled   bool
	SourceIPs string
}
