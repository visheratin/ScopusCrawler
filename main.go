package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/visheratin/scopus-crawler/config"
	"github.com/visheratin/scopus-crawler/crawler"
	"github.com/visheratin/scopus-crawler/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	config.InitKeys("keys.txt")
	manager := crawler.Manager{}
	manager.Init("data-sources.json", 4)
	req, err := readRequest("request.json")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	err = manager.StartCrawling(req)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	fmt.Scanln()
	// address := "http://api.elsevier.com/content/serial/title"
	// params := map[string]string{}
	// params["httpAccept"] = "application/json"
	// params["apiKey"] = "bd9ddf64bbcc7ed6d09ddcc16d607d75"
	// params["start"] = "50"
	// data, err := query.MakeQuery(address, params)
	// if err != nil {
	// 	logger.Error.Println(err)
	// }
	// fmt.Println(data["serial-metadata-response"])
}

func readRequest(requestPath string) (crawler.SearchRequest, error) {
	var req crawler.SearchRequest
	file, err := os.Open(requestPath)
	if err != nil {
		logger.Error.Println("Unable to open request file.")
		return req, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&req)
	if err != nil {
		logger.Error.Println(err)
		return req, err
	}
	return req, nil
}
