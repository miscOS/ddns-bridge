package services

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/miscOS/ddns-bridge/models"
)

type Service interface {
	GetID() string
	Init(config map[string]interface{}) error
	Update(value *models.UpdateRequest) ([]models.UpdateResult, error)
}

type ServiceRegistry struct {
	services map[string]Service
	mu       sync.RWMutex
}

var registry *ServiceRegistry = &ServiceRegistry{
	services: make(map[string]Service),
	mu:       sync.RWMutex{},
}

func Register(service Service) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.services[service.GetID()] = service
}

func Get(name string) (Service, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	service, ok := registry.services[name]
	if !ok {
		return nil, fmt.Errorf("service '%s' not found", name)
	}
	return service, nil
}

func GetServiceInformation() []map[string]interface{} {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	infos := make([]map[string]interface{}, 0)

	for name, service := range registry.services {
		info := retrieveServiceTag(service)
		info["id"] = name
		infos = append(infos, info)
	}
	return infos
}

func retrieveServiceTag(service Service) map[string]interface{} {
	serviceElem := reflect.ValueOf(service).Elem()

	info := make(map[string]interface{})

	for i := 0; i < serviceElem.NumField(); i++ {
		field := serviceElem.Type().Field(i)
		fieldValue := serviceElem.Field(i)

		serviceTag := field.Tag.Get("service")

		if fieldValue.Kind() == reflect.Slice && serviceTag == "parameters" {
			paramType := fieldValue.Type().Elem()
			params := extractServiceParamsTag(paramType)
			info["parameters"] = params
		} else {
			info[serviceTag] = fieldValue.String()
		}
	}
	return info
}

func extractServiceParamsTag(paramType reflect.Type) []map[string]interface{} {
	var result []map[string]interface{}

	for i := 0; i < paramType.NumField(); i++ {
		paramMap := make(map[string]interface{})

		field := paramType.Field(i)
		if field.Tag.Get("name") != "" {
			paramMap["name"] = field.Tag.Get("name")
		}
		if field.Tag.Get("type") != "" {
			paramMap["type"] = field.Tag.Get("type")
		}
		if field.Tag.Get("req") != "" {
			paramMap["required"] = field.Tag.Get("req") == "true"
		}
		if field.Tag.Get("desc") != "" {
			paramMap["description"] = field.Tag.Get("desc")
		}
		result = append(result, paramMap)
	}
	return result
}
