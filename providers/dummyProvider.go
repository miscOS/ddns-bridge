package DNSProvider

import "log"

type dummyDNS struct {
	Name string
}

func (d *dummyDNS) Setup(config interface{}) error {
	d.Name = config.(string)
	log.Printf("Configuring dummy provider with name: %s", d.Name)
	return nil
}

func (d *dummyDNS) Update(params *DNSParams) error {
	log.Printf("Updating DNS for dummy provider: %s", d.Name)

	log.Printf("IPv4: %s", params.IPv4)
	log.Printf("IPv6: %s", params.IPv6)

	return nil
}
