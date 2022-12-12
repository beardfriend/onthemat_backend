package infrastructure

import (
	"fmt"
	"io/ioutil"
	"log"

	"onthemat/internal/app/config"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearch(c *config.Config, certUrl string) *elasticsearch.Client {
	cert, err := ioutil.ReadFile(certUrl)
	if err != nil {
		log.Fatalf("fail %v", err)
	}

	address := fmt.Sprintf("https://%s:%d", c.Elastic.Host, c.Elastic.Port)
	cnf := elasticsearch.Config{
		Addresses: []string{
			address,
		},
		CACert:   cert,
		Username: c.Elastic.User,
		Password: c.Elastic.Password,
	}
	client, err := elasticsearch.NewClient(cnf)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	resp, err := client.Info()
	if err != nil {
		log.Fatalf("helath check error")
	}

	fmt.Printf("helath check %s \n", resp)

	return client
}
