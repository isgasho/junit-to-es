package idxtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ResultsConfig configures the results.
type ResultsConfig struct {
	Client         *elasticsearch.Client
	IndexSuiteName string
	IndexTestName  string
}

// Results allows to index test results.
type Results struct {
	es           *elasticsearch.Client
	idxSuiteName string
	idxTestName  string
}

// NewResults returns a new instance of the test results.
func NewResults(c ResultsConfig) (*Results, error) {
	idxSuiteName := c.IndexSuiteName
	if idxSuiteName == "" {
		idxSuiteName = "test-suites"
	}

	idxTestName := c.IndexTestName
	if idxTestName == "" {
		idxTestName = "test-cases"
	}

	r := Results{es: c.Client, idxSuiteName: idxSuiteName, idxTestName: idxTestName}
	return &r, nil
}

// CreateIndexes creates a new index with mapping
func (r *Results) CreateIndexes(suiteMapping, testMapping string) error {
	res, err := r.es.Indices.Create(r.idxSuiteName, r.es.Indices.Create.WithBody(strings.NewReader((suiteMapping))))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	res, err = r.es.Indices.Create(r.idxTestName, r.es.Indices.Create.WithBody(strings.NewReader(testMapping)))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}

// CreateTest indexes a new test result into associated index.
func (r *Results) CreateTest(item *Test) error {
	payload, err := json.Marshal(item)
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index: r.idxTestName,
		Body:  bytes.NewReader(payload),
	}.Do(ctx, r.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return nil
}

// CreateSuite indexes a new suite result into associated index.
func (r *Results) CreateSuite(item *Suite) error {
	payload, err := json.Marshal(item)
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index: r.idxSuiteName,
		Body:  bytes.NewReader(payload),
	}.Do(ctx, r.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return nil
}
