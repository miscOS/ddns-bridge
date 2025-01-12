package services

import (
	"encoding/json"
	"log"
)

type dummyDNS struct {
	Name string
}

func (d *dummyDNS) Setup(config string) error {

	err := json.Unmarshal([]byte(config), &d)
	return err
}

func (d *dummyDNS) Update(params *DNSParams) error {
	log.Printf("Updating DNS for dummy provider: %s", d.Name)

	log.Printf("IPv4: %s", params.IPv4)
	log.Printf("IPv6: %s", params.IPv6)

	return nil
}
