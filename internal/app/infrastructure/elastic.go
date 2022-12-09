package infrastructure

import (
	"fmt"
	"io/ioutil"
	"log"

	"onthemat/internal/app/config"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearch(c *config.Config) *elasticsearch.Client {
	cert, err := ioutil.ReadFile("../../configs/elastic.crt")
	if err != nil {
		panic(err)
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
		panic(err)
	}
	resp, err := client.Info()
	if err != nil {
		fmt.Println("health check error")
	}
	fmt.Printf("helath check %s \n", resp)

	return client
}
