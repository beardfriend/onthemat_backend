package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearch() {
	cert, er := ioutil.ReadFile("../../../elastic.crt")
	fmt.Println(er)
	cnf := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		CACert:   cert,
		Username: "elastic",
		Password: "asd1234",
	}
	es, err := elasticsearch.NewClient(cnf)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())

	dataas := map[string]interface{}{
		"asd": "hello",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(dataas)

	asd, e := es.Create("index_t", "asd1", &buf)
	fmt.Println(asd, e)
	ff, _ := es.Get("index_t", "1")
	fmt.Println(ff)
}
