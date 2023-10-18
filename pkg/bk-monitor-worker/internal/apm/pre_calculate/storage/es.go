package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"time"
)

const (
	SavePrecalculateResult Action = "savePrecalculateResult"
)

type Converter func(io.ReadCloser) (any, error)

var (
	BytesConverter Converter = func(body io.ReadCloser) (any, error) {
		var resInstance EsQueryResult
		if err := json.NewDecoder(body).Decode(&resInstance); err != nil {
			return nil, err
		}

		var resMap []map[string]any
		for _, item := range resInstance.Hits.Hits {
			resMap = append(resMap, item.Source)
		}
		r, _ := json.Marshal(resMap)
		return r, nil
	}

	AggsCountConvert Converter = func(body io.ReadCloser) (any, error) {
		var resInstance EsQueryResult
		if err := json.NewDecoder(body).Decode(&resInstance); err != nil {
			return nil, err
		}

		res := make(map[string]int)
		for k, item := range resInstance.Aggregations {
			v, vE := item["value"]
			if vE {
				res[k] = v
			}
		}

		return res, nil
	}
)

type EsStorageData struct {
	DocumentId string
	Value      []byte
}

type EsQueryData struct {
	Body      map[string]any
	Converter Converter
}

type EsQueryResult struct {
	Hits         EsQueryHit                `json:"hits"`
	Aggregations map[string]map[string]int `json:"aggregations"`
}

type EsQueryHit struct {
	Total EsQueryHitTotal  `json:"total"`
	Hits  []EsQueryHitItem `json:"hits"`
}

type EsQueryHitItem struct {
	Source map[string]any `json:"_source"`
}

type EsQueryHitTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type EsOption func(options *EsOptions)

type EsOptions struct {
	indexName string
	host      string
	username  string
	password  string
}

func EsIndexName(n string) EsOption {
	return func(options *EsOptions) {
		options.indexName = n
	}
}

func EsHost(h string) EsOption {
	return func(options *EsOptions) {
		options.host = h
	}
}

func EsUsername(u string) EsOption {
	return func(options *EsOptions) {
		options.username = u
	}
}

func EsPassword(p string) EsOption {
	return func(options *EsOptions) {
		options.password = p
	}
}

type esStorage struct {
	indexName string
	client    *elasticsearch.Client
}

func (e *esStorage) Save(data EsStorageData) error {

	buf := bytes.NewBuffer(data.Value)
	req := esapi.IndexRequest{Index: e.getSaveIndexName(e.indexName), DocumentID: data.DocumentId, Body: buf}
	res, err := req.Do(context.Background(), e.client)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warnf("failed to close the body")
		}
	}(res.Body)

	if err != nil {
		return err
	}

	return nil
}

func (e *esStorage) SaveBatch(items []EsStorageData) error {
	var buf bytes.Buffer
	for _, item := range items {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, item.DocumentId, "\n"))
		data := item.Value
		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	req := esapi.BulkRequest{Index: e.getSaveIndexName(e.indexName), Body: &buf}
	response, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warnf("failed to close the body")
		}
	}(response.Body)

	if response.IsError() {
		return fmt.Errorf("bulk insert returned an abnormal status code： %d", response.StatusCode)
	}

	return nil
}

func (e *esStorage) Query(data any) (any, error) {
	body := data.(EsQueryData)
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(body.Body); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(e.indexName),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Warnf("failed to close the body")
		}
	}(res.Body)

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	return body.Converter(res.Body)
}

func (e *esStorage) getSaveIndexName(indexName string) string {
	return fmt.Sprintf("write_%s_%s", time.Now().Format("20060102"), indexName)
}

func newEsStorage(options EsOptions) (*esStorage, error) {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{options.host},
		Username:  options.username,
		Password:  options.password,
	})
	if err != nil {
		return nil, err
	}
	logger.Infof("create ES Client")
	return &esStorage{client: c, indexName: options.indexName}, nil
}
