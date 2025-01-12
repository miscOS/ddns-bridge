package services

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/imroc/req/v3"
)

type cloudflareDNS struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	Token     string `json:"token"`
	client    *req.Client
	url       string
}

type cf_Body struct {
	Content string `json:"content"`
}

type cf_DNS struct {
	ID           string `json:"id"`
	ZoneID       string `json:"zone_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Content      string `json:"content"`
	LastModified string `json:"modified_on"`
}

type cf_DNS_Response struct {
	Success bool     `json:"success"`
	Result  []cf_DNS `json:"result"`
}

type cf_DNS_Update_Response struct {
	Success bool   `json:"success"`
	Result  cf_DNS `json:"result"`
}

type cf_Zone struct {
	ID string `json:"id"`
}

type cf_Zone_Response struct {
	Success bool      `json:"success"`
	Result  []cf_Zone `json:"result"`
}

type cf_Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cf_Error_Response struct {
	Success bool       `json:"success"`
	Errors  []cf_Error `json:"errors"`
}

func (s *cloudflareDNS) Setup(config string) error {

	if err := json.Unmarshal([]byte(config), &s); err != nil {
		log.Printf("Error unmarshalling config: %s", err)
		return err
	}

	s.client = req.C().SetUserAgent("miscOS/ddnsbridge-cloudflare")
	s.url = "https://api.cloudflare.com/client/v4"
	return nil
}

func (s *cloudflareDNS) Update(params *DNSValues) (result []DNSResult, err error) {

	zones, err := s.getZone()
	if err != nil {
		return result, err
	}

	dns, err := s.getDNS(zones)
	if err != nil {
		return result, err
	}

	updates, err := s.updateDNS(dns, params)
	if err != nil {
		return result, err
	}

	for _, r := range updates {
		result = append(result, DNSResult{Success: r.Success, Domain: r.Result.Name, Record: r.Result.Type})
	}

	return result, nil
}

func (s *cloudflareDNS) getZone() (response cf_Zone_Response, err error) {

	var errorResponse cf_Error_Response

	_, err = s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBearerAuthToken(s.Token).
		SetPathParam("name", s.Domain).
		SetSuccessResult(&response).
		SetErrorResult(&errorResponse).
		Get(s.url + "/zones?name={name}&status=active")

	if err != nil {
		return response, err
	}

	if response.Success {
		return response, nil
	} else {
		return response, errors.New(errorResponse.Errors[0].Message)
	}
}

func (s *cloudflareDNS) getDNS(zones cf_Zone_Response) (response cf_DNS_Response, err error) {

	var errorResponse cf_Error_Response
	var name string

	if s.Subdomain == "" {
		name = s.Domain
	} else {
		name = s.Subdomain + "." + s.Domain
	}

	for _, zone := range zones.Result {

		_, err = s.client.R().
			SetHeader("Content-Type", "application/json").
			SetBearerAuthToken(s.Token).
			SetPathParam("zone_id", zone.ID).
			SetPathParam("name", name).
			SetSuccessResult(&response).
			SetErrorResult(&errorResponse).
			Get(s.url + "/zones/{zone_id}/dns_records?name={name}")

		if err != nil {
			return response, err
		}

		if response.Success {
			return response, nil
		} else {
			return response, errors.New(errorResponse.Errors[0].Message)
		}
	}

	return response, errors.New("no DNS records found")
}

func (s *cloudflareDNS) updateDNS(dns cf_DNS_Response, v *DNSValues) (responses []cf_DNS_Update_Response, err error) {

	for _, record := range dns.Result {

		var response cf_DNS_Update_Response
		var errorResponse cf_Error_Response
		var body cf_Body

		if record.Type == "A" && v.IPv4.Is4() {
			body.Content = v.IPv4.String()
		} else if record.Type == "AAAA" && v.IPv6.Is6() {
			body.Content = v.IPv6.String()
		} else {
			continue
		}

		_, err = s.client.R().
			SetHeader("Content-Type", "application/json").
			SetBearerAuthToken(s.Token).
			SetBody(body).
			SetPathParam("zone_id", record.ZoneID).
			SetPathParam("dns_id", record.ID).
			SetSuccessResult(&response).
			SetErrorResult(&errorResponse).
			Patch(s.url + "/zones/{zone_id}/dns_records/{dns_id}")

		if err != nil {
			return responses, err
		}

		if response.Success {
			responses = append(responses, response)
		} else {
			return responses, errors.New(errorResponse.Errors[0].Message)
		}
	}
	return responses, nil
}
