package models

import "net/netip"

type UpdateRequest struct {
	IPv4 netip.Addr `json:"ipv4"`
	IPv6 netip.Addr `json:"ipv6"`
}

type UpdateResult struct {
	Success bool   `json:"success"`
	Task    string `json:"task"`
	Message string `json:"message"`
}
