package query

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/visheratin/scopus-crawler/config"
	"github.com/visheratin/scopus-crawler/storage"
)

func MakeQuery(address string, id string, params map[string]string, timeoutSec int,
	storage storage.GenericStorage, config config.Configuration) (map[string]interface{}, error) {
	requestPath := address
	if id != "" {
		requestPath = strings.Replace(requestPath, "{_id_}", id, -1)
		//requestPath = requestPath + "/" + id
	}
	for key, value := range params {
		requestPath += key + "=" + value + "&"
	}
	var data map[string]interface{}
	var body []byte
	finishedRequest, _ := storage.GetFinishedRequest(requestPath)
	if finishedRequest == "" {
		request := requestPath
		authKey := config.GetKey()
		requestPath = requestPath + "apiKey=" + authKey
		fmt.Println(requestPath)
		req, err := http.NewRequest("GET", requestPath, nil)
		if err != nil {
			return nil, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			config.RemoveKey(authKey)
			return nil, err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		storage.CreateFinishedRequest(request, string(body))
	} else {
		body = []byte(finishedRequest)
	}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	duration := time.Duration(timeoutSec) * time.Second
	time.Sleep(duration)
	return data, nil
}
