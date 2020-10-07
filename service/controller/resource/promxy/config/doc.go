package proxyconfig

// Config defines the content of promxy configuration file.
// This structure is copied over from https://raw.githubusercontent.com/jacksontj/promxy/v0.0.60/pkg/config/config.go

// As jacksontj/promxy is relying on prometheus 1.8 and we are using prometheus 2.20 which implies some changes in prometheus types
// and makes the defautl promxy struct incompatible with our codebase.
// TODO, move tojacksontj/promxy once https://github.com/jacksontj/promxy/issues/352 is resolved
