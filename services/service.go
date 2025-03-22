package services

import (
	"reflect"

	"github.com/miscOS/ddns-bridge/models"
)

type Service interface {
	Setup(config map[string]interface{}) error
	Update(value *models.UpdaetValue) ([]models.UpdateResult, error)
}

func GetService(service string) Service {
	if p := getMappedServices()[service]; p == nil {
		return nil
	} else {
		return reflect.New(p).Interface().(Service)
	}
}

func getMappedServices() map[string]reflect.Type {
	return map[string]reflect.Type{
		"cloudflare": reflect.TypeOf((*CloudflareDNS)(nil)).Elem(),
		"dummy":      reflect.TypeOf((*DummyDNS)(nil)).Elem(),
	}
}
