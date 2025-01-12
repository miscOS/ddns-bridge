package services

import (
	"encoding/json"
)

type dummyDNS struct {
}

func (d *dummyDNS) Setup(config string) error {

	err := json.Unmarshal([]byte(config), &d)
	return err
}

func (d *dummyDNS) Update(v *DNSValues) (result []DNSResult, err error) {

	result = append(result, DNSResult{Success: true, Domain: "dummy", Record: "A"})
	return result, err
}
