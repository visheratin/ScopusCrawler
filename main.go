package main

import (
	"fmt"

	"github.com/visheratin/scopus-crawler/config"
	"github.com/visheratin/scopus-crawler/logger"
	"github.com/visheratin/scopus-crawler/query"
)

func main() {
	logger.InitLog()
	config.InitKeys("key.txt")
	address := "http://api.elsevier.com/content/serial/title"
	params := map[string]string{}
	params["httpAccept"] = "application/json"
	params["apiKey"] = "bd9ddf64bbcc7ed6d09ddcc16d607d75"
	params["start"] = "50"
	data, err := query.MakeQuery(address, params)
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Println(data["serial-metadata-response"])
}
