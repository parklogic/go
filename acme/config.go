package acme

import (
	"time"
)

type Configuration struct {
	CreateNewAccount bool          `description:"Allow creating new ACME accounts"`
	Directory        string        `description:"ACME directory URL"`
	KeyPath          string        `description:"Path to the ACME account key"`
	MaxLabels        int           `description:"Maximum number of labels the directory accepts in the FQDN"`
	PollInterval     time.Duration `description:"How often to poll the ACME server for updates"`
	PollTimeout      time.Duration `description:"How long to wait for the order and authorizations to complete"`
	UserAgent        string        `description:"User agent to send to the ACME server"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Directory:    "https://acme-v02.api.letsencrypt.org/directory",
		KeyPath:      "account.key",
		MaxLabels:    10, // todo: implement server-side filtering
		PollInterval: 250 * time.Millisecond,
		PollTimeout:  5 * time.Minute,
		UserAgent:    "parklogic",
	}
}
