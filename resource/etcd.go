package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func newEtcd(cfg Config) (provider, error) {
	// Since we are making plain HTTP requests atm, there isn't much going on here
	// If the etcd client stabilizes, more will go here
	return &etcdProvider{}, nil
}

type etcdProvider struct {
}

func (etcd *etcdProvider) Save(key string, value string) error {
	// This will need to be replaced with the docker service name
	url := "http://127.0.0.1:2379/v3/kv/put"
	postBody, err := json.Marshal(map[string]string{
		"key":   key,
		"value": value,
	})
	if err != nil {
		return err
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("etcd response not %d: %v", http.StatusOK, resp.Status)
	}
	return nil
}

func (etcd *etcdProvider) Query(key string) (string, error) {
	// This will need to be replaced with the docker service name
	url := "http://127.0.0.1:2379/v3/kv/range"
	postBody, _ := json.Marshal(map[string]string{
		"key": key,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("etcd response not %d: %v", http.StatusOK, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var keys etcdResponse
	err = json.Unmarshal(bodyBytes, &keys)
	if err != nil {
		return "", err
	}
	resourceString := string(bodyBytes)

	return resourceString, nil
}

type etcdResponse struct {
	Keys []key `json:"kvs"`
}

type key struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
