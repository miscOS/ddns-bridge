package helpers

import (
	"reflect"

	"github.com/miscOS/ddns-bridge/models"
	"github.com/miscOS/ddns-bridge/services"
)

type Service interface {
	Setup(config string) error
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
		"cloudflare": reflect.TypeOf((*services.CloudflareDNS)(nil)).Elem(),
		"dummy":      reflect.TypeOf((*services.DummyDNS)(nil)).Elem(),
	}
}
