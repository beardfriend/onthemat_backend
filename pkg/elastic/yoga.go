package elastic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func InitYoga(ela *elasticsearch.Client) error {
	res, err := ela.Indices.GetFieldMapping([]string{"yoga"})
	if err != nil {
		return err
	}

	if res.StatusCode != 404 {
		return errors.New("이미 존재")
	}

	settings := `
	{
		"settings": {
		  "analysis": {
			"filter": {
			  "suggest_filter": {
				"type": "edge_ngram",
				"min_gram": 1,
				"max_gram": 50
			  }
			},
			"analyzer": {
			  "autocomplete": {
				"tokenizer": "autocomplete",
				"filter": [
				  "lowercase"
				]
			  },
			  "autocomplete_search": {
				"tokenizer": "lowercase"
			  },
			  "suggest_search_analyzer": {
				"type": "custom",
				"tokenizer": "jaso_tokenizer"
			  },
			  "suggest_index_analyzer": {
				"type": "custom",
				"tokenizer": "jaso_tokenizer",
				"filter": [
				  "suggest_filter"
				]
			  }
			},
			"tokenizer": {
			  "autocomplete": {
				"type": "edge_ngram",
				"min_gram": 1,
				"max_gram": 50,
				"token_chars": [
				  "letter"
				]
			  }
			}
		  }
		}
	}
	`
	sbuf := strings.NewReader(settings)
	req := esapi.IndicesCreateRequest{
		Index: "yoga",
		Body:  sbuf,
	}
	resp, err := req.Do(context.Background(), ela)

	fmt.Println(resp)
	fmt.Println(err)

	data := `
	{
		"properties": {
			"name": {
			  "type": "text",
			  "store": true,
			  "analyzer": "suggest_index_analyzer",
			  "search_analyzer": "suggest_search_analyzer"
			}
		  }
	}
	`
	pbuf := strings.NewReader(data)
	resp, err = ela.Indices.PutMapping([]string{"yoga"}, pbuf)

	fmt.Println(resp)
	return err
}
