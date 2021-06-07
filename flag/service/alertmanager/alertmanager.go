package alertmanager

type Alertmanager struct {
	Address    string
	BaseDomain string
	LogLevel   string
	Storage    AlertmanagerStorage
	Version    string
}

type AlertmanagerStorage struct {
	CreatePVC string
	Size      string
}
