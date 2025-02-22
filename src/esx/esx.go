package esx

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/trancecho/open-sdk/config"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	Es *elasticsearch.Client
)

func InitElastic() {
	addresses := config.GetConfig().Elasticsearch.Addresses
	//addresses := []string{"http://localhost:9200"}
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	var err error
	Es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	log.Printf("Connected to Elasticsearch at %s", cfg.Addresses)
}

func CreateIndex(index string) (exist bool) {
	// 检查索引是否存在
	exists, err := Es.Indices.Exists([]string{index})
	if err != nil {
		log.Fatalf("Error checking if index exists: %s", err)
	}
	if exists.StatusCode == 200 {
		log.Printf("Index %s already exists", index)
		return true
	}

	// 创建索引
	_, err = Es.Indices.Create(index)

	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	log.Printf("Index %s created", index)
	return false
}

func IndexDocument(index, body string) error {
	req := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(body),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), Es)
	if err != nil {
		log.Printf("Error indexing document to index %s: %s", index, err)
		return err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %s", err)
		}
	}()

	if res.IsError() {
		log.Printf("Error response from server: %s", res.String())
		return fmt.Errorf("error indexing document to index %s: %s", index, res.String())
	}

	log.Printf("Document indexed in index %s: %s", index, res.String())
	return nil
}

func GetDocument(index, id string) (result any, err error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}
	res, err := req.Do(context.Background(), Es)
	if err != nil {
		log.Fatalf("Error getting document: %s", err)
	}
	defer res.Body.Close()
	var r map[string]interface{}
	if res.IsError() {
		log.Printf("Document %s not found in index %s", id, index)
		return nil, fmt.Errorf("document not found in index %s", index)
	} else {

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %w", err)
		}
		log.Printf("Document %s retrieved from index %s: %s", id, index, r["_source"])
	}
	return r["_source"], nil
}

func SearchDocuments(index, query string) (interface{}, error) {
	var buf strings.Builder
	// Construct the search query as JSON with fuzziness
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": map[string]interface{}{
					"query":     query,
					"fuzziness": "AUTO", // 自动处理模糊度
				},
			},
		},
	}
	// Encode the query to JSON
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("error encoding search query: %w", err)
	}

	// Perform the search request
	res, err := Es.Search(
		Es.Search.WithContext(context.Background()),
		Es.Search.WithIndex(index),
		Es.Search.WithBody(strings.NewReader(buf.String())),
		Es.Search.WithTrackTotalHits(true),
		Es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %w", err)
	}
	defer res.Body.Close()

	// Check if the response has an error
	if res.IsError() {
		var esErr map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&esErr); err != nil {
			return nil, fmt.Errorf("error parsing error response: %w", err)
		}
		errorType := esErr["error"].(map[string]interface{})["type"]
		errorReason := esErr["error"].(map[string]interface{})["reason"]
		return nil, fmt.Errorf("search error: %s, reason: %s", errorType, errorReason)
	}

	// Decode the response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	// Return search hits
	return r["hits"], nil
}

func UpdateDocument(index, id, body string) (bool, error) {
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(body), // 使用 strings.NewReader
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), Es)
	if err != nil {
		log.Fatalf("Error updating document: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return false, fmt.Errorf("error updating document: %s", res.String())
	}
	log.Printf("Document %s updated in index %s", id, index)
	return true, nil
}

func DeleteDocument(index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), Es)
	if err != nil {
		return fmt.Errorf("error deleting document: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}
	log.Printf("Document %s deleted from index %s", id, index)
	return nil
}

func main() {
	InitElastic()
	indexName := "mama"

	// 创建索引
	CreateIndex(indexName)
	//
	//// 索引文档
	//docID := "1"
	//docBody := `{"title": "Test Document", "content": "This is a test document."}`
	//// 添加行
	//IndexDocument(indexName, docID, docBody)
	//
	//// 获取文档
	//GetDocument(indexName, docID)
	////
	//// 搜索文档
	//SearchDocuments(indexName, "test")
	//
	//// 更新文档
	//updateBody := `{"doc": {"content": "This is an updated test document."}}`
	//UpdateDocument(indexName, docID, updateBody)
	//
	//// 获取更新后的文档
	//GetDocument(indexName, docID)
	//
	//// 删除文档
	//DeleteDocument(indexName, docID)
	//
	//// 尝试获取已删除的文档
	//GetDocument(indexName, docID)
}
