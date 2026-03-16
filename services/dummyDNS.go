package services

import (
	"fmt"

	"github.com/miscOS/ddns-bridge/models"
)

type exampleDNS struct {
	ID          string        `service:"name"`
	Description string        `service:"description"`
	Params      exampleParams `service:"parameters"`
}

type exampleParams struct {
	APIKey    string `name:"api_key" type:"string" req:"true" desc:"API Key for authentication."`
	APISecret string `name:"api_secret" type:"string" req:"false" desc:"API Secret for authentication."`
}

func NewExampleDNS() *exampleDNS {
	return &exampleDNS{
		ID:          "example",
		Description: "An example DNS service for testing purposes.",
	}
}

func (d *exampleDNS) GetID() string {
	return d.ID
}

func (d *exampleDNS) Init(params map[string]interface{}) error {

	fmt.Println("ExampleDNS Setup")
	fmt.Println(params)

	return nil
}

func (d *exampleDNS) Update(v *models.UpdateRequest) (result []models.UpdateResult, err error) {

	result = append(result, models.UpdateResult{Success: true, Message: "exampleDNS Update successful"})
	return result, err
}
