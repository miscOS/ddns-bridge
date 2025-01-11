package DNSProvider

import (
	"net/netip"
	"reflect"

	"github.com/pkg/errors"
)

type DNSParams struct {
	IPv4 netip.Addr
	IPv6 netip.Addr
}

type DNSProvider interface {
	Setup(interface{}) error
	Update(params *DNSParams) error
}

// GetDNSProvider returns a DNSProvider instance based on the provided name.
// If the name does not match any known DNS provider, an error is returned.
//
// Parameters:
//   - name: The name of the DNS provider to retrieve.
//
// Returns:
//   - DNSProvider: The DNS provider instance corresponding to the given name.
//   - error: An error if the DNS provider name is unknown.
func GetDNSProvider(name string) (DNSProvider, error) {
	if p := getDNSProviderMap()[name]; p == nil {
		return nil, errors.Errorf("Unknown dns provider: \"%s\"", name)
	} else {
		return reflect.New(p).Interface().(DNSProvider), nil
	}
}

// returns a map where the keys are the names of supported
// DNS providers and the values are the corresponding reflect.Type of
// the provider's struct. This map is used to dynamically create instances of
// the DNS provider structs based on the provider name.
func getDNSProviderMap() map[string]reflect.Type {
	return map[string]reflect.Type{
		"dummy": reflect.TypeOf((*dummyDNS)(nil)).Elem(),
	}
}
