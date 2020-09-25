package security

type Security struct {
	RestrictedAccess RestrictedAccess
}

type RestrictedAccess struct {
	Enabled string
	Subnets string
}
