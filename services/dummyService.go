package services

import (
	"github.com/miscOS/ddns-bridge/models"
	"github.com/mitchellh/mapstructure"
)

type DummyDNS struct {
}

func (d *DummyDNS) Setup(params map[string]interface{}) error {

	if err := mapstructure.Decode(params, d); err != nil {
		return err
	}
	return nil
}

func (d *DummyDNS) Update(v *models.UpdaetValue) (result []models.UpdateResult, err error) {

	result = append(result, models.UpdateResult{Success: true, Domain: "dummy", Record: "A"})
	return result, err
}
