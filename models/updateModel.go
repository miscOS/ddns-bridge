package models

import "net/netip"

type UpdaetValue struct {
	IPv4 netip.Addr `json:"ipv4"`
	IPv6 netip.Addr `json:"ipv6"`
}

type UpdateResult struct {
	Success bool   `json:"success"`
	Domain  string `json:"domain"`
	Record  string `json:"record"`
}
