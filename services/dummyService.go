package services

import (
	"encoding/json"

	"github.com/miscOS/ddns-bridge/models"
)

type DummyDNS struct {
}

func (d *DummyDNS) Setup(config string) error {

	err := json.Unmarshal([]byte(config), &d)
	return err
}

func (d *DummyDNS) Update(v *models.UpdaetValue) (result []models.UpdateResult, err error) {

	result = append(result, models.UpdateResult{Success: true, Domain: "dummy", Record: "A"})
	return result, err
}
