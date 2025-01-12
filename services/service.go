package services

import (
	"net/netip"
	"reflect"

	"github.com/pkg/errors"
)

type DNSValues struct {
	IPv4 netip.Addr
	IPv6 netip.Addr
}

type DNSResult struct {
	Success bool   `json:"success"`
	Domain  string `json:"domain"`
	Record  string `json:"record"`
}

type DNSService interface {
	Setup(configObject string) error
	Update(params *DNSValues) ([]DNSResult, error)
}

// GetDNSService returns a DNSProvider instance based on the provided name.
// If the name does not match any known DNS provider, an error is returned.
//
// Parameters:
//   - name: The name of the DNS provider to retrieve.
//
// Returns:
//   - DNSProvider: The DNS provider instance corresponding to the given name.
//   - error: An error if the DNS provider name is unknown.
func GetDNSService(name string) (DNSService, error) {
	if p := getDNSServiceMap()[name]; p == nil {
		return nil, errors.Errorf("Unknown dns provider: \"%s\"", name)
	} else {
		return reflect.New(p).Interface().(DNSService), nil
	}
}

// returns a map where the keys are the names of supported
// DNS providers and the values are the corresponding reflect.Type of
// the provider's struct. This map is used to dynamically create instances of
// the DNS provider structs based on the provider name.
func getDNSServiceMap() map[string]reflect.Type {
	return map[string]reflect.Type{
		"dummy":      reflect.TypeOf((*dummyDNS)(nil)).Elem(),
		"cloudflare": reflect.TypeOf((*cloudflareDNS)(nil)).Elem(),
	}
}
