package alertmanager

type Alertmanager struct {
	Address  string
	LogLevel string
	Storage  AlertmanagerStorage
	Version  string
}

type AlertmanagerStorage struct {
	CreatePVC string
	Size      string
}
