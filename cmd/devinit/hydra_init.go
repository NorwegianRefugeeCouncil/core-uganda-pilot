package devinit

import (
	"encoding/json"
	"os"
	"path"
)

type HydraClients struct {
	Clients []ClientConfig `json:"clients"`
}

func (c *Config) makeHydraInit() error {

	jsonBytes, err := json.MarshalIndent(HydraClients{
		Clients: c.hydraClients,
	}, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(HydraCredsDir, "clients.json"), jsonBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}
