package DNSProvider

import "log"

type dummyDNS struct {
	Name string
}

func (d *dummyDNS) Configure(config interface{}) error {
	d.Name = config.(string)
	log.Printf("Configuring dummy provider with name: %s", d.Name)
	return nil
}

func (d *dummyDNS) Update() error {
	log.Printf("Updating DNS for dummy provider: %s", d.Name)
	return nil
}
