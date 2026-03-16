package services

import (
	"errors"

	"github.com/imroc/req/v3"
	"github.com/miscOS/ddns-bridge/models"
	"github.com/mitchellh/mapstructure"
)

type cloudflareDNS struct {
	ID          string           `service:"name"`
	Description string           `service:"description"`
	Params      cloudflareParams `service:"parameters"`
	client      *req.Client
	url         string
}

type cloudflareParams struct {
	Domain    string `name:"domain" type:"string" req:"true" desc:"The domain to update."`
	Subdomain string `name:"subdomain" type:"string" req:"false" desc:"The subdomain to update. If not provided, the root domain will be updated."`
	Token     string `name:"token" type:"string" req:"true" desc:"API Token with permissions to edit DNS records."`
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

func NewCloudflareDNS() *cloudflareDNS {
	return &cloudflareDNS{
		ID:          "cloudflare",
		Description: "A Cloudflare DNS service for updating A and AAAA records.",
		client:      req.C().SetUserAgent("miscOS/ddnsbridge-cloudflare"),
		url:         "https://api.cloudflare.com/client/v4",
	}
}

func (d *cloudflareDNS) GetID() string {
	return d.ID
}

func (d *cloudflareDNS) Init(params map[string]interface{}) error {

	err := mapstructure.Decode(params, &d.Params)
	return err
}

func (d *cloudflareDNS) Update(value *models.UpdateRequest) (result []models.UpdateResult, err error) {

	zones, err := d.getZone()
	if err != nil {
		return result, err
	}

	dns, err := d.getDNS(zones)
	if err != nil {
		return result, err
	}

	updates, err := d.updateDNS(zones, dns, value)
	if err != nil {
		return result, err
	}

	for _, u := range updates {
		result = append(result, models.UpdateResult{Success: u.Success, Task: u.Result.Name, Message: u.Result.Type})
	}

	return result, nil
}

func (d *cloudflareDNS) getZone() (response cf_Zone_Response, err error) {

	var errorResponse cf_Error_Response

	_, err = d.client.R().
		SetHeader("Content-Type", "application/json").
		SetBearerAuthToken(d.Params.Token).
		SetPathParam("name", d.Params.Domain).
		SetSuccessResult(&response).
		SetErrorResult(&errorResponse).
		Get(d.url + "/zones?name={name}&status=active")

	if err != nil {
		return response, err
	}

	if response.Success {
		return response, nil
	} else {
		return response, errors.New(errorResponse.Errors[0].Message)
	}
}

func (d *cloudflareDNS) getDNS(zones cf_Zone_Response) (response cf_DNS_Response, err error) {

	var errorResponse cf_Error_Response
	var name string

	if d.Params.Subdomain == "" {
		name = d.Params.Domain
	} else {
		name = d.Params.Subdomain + "." + d.Params.Domain
	}

	for _, zone := range zones.Result {

		_, err = d.client.R().
			SetHeader("Content-Type", "application/json").
			SetBearerAuthToken(d.Params.Token).
			SetPathParam("zone_id", zone.ID).
			SetPathParam("name", name).
			SetSuccessResult(&response).
			SetErrorResult(&errorResponse).
			Get(d.url + "/zones/{zone_id}/dns_records?name={name}")

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

func (d *cloudflareDNS) updateDNS(zones cf_Zone_Response, dns cf_DNS_Response, r *models.UpdateRequest) (responses []cf_DNS_Update_Response, err error) {

	zoneID := zones.Result[0].ID // Quick and dirty fix for now, since dns records dont have zone id anymore

	for _, record := range dns.Result {

		var response cf_DNS_Update_Response
		var errorResponse cf_Error_Response
		var body cf_Body

		if record.Type == "A" && r.IPv4.Is4() {
			body.Content = r.IPv4.String()
		} else if record.Type == "AAAA" && r.IPv6.Is6() {
			body.Content = r.IPv6.String()
		} else {
			continue
		}

		_, err = d.client.R().
			SetHeader("Content-Type", "application/json").
			SetBearerAuthToken(d.Params.Token).
			SetBody(body).
			SetPathParam("zone_id", zoneID).
			SetPathParam("dns_id", record.ID).
			SetSuccessResult(&response).
			SetErrorResult(&errorResponse).
			Patch(d.url + "/zones/{zone_id}/dns_records/{dns_id}")

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
