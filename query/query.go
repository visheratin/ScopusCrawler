package query

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func MakeQuery(address string, params map[string]string) (map[string]interface{}, error) {
	requestPath := address + "?"
	for key, value := range params {
		requestPath += key + "=" + value + "&"
	}
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
