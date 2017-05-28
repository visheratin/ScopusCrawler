package crawler

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Manager struct {
	DataSources []DataSource
	Queue       []SearchRequest
	Workers     map[string]Worker
}

func (manager Manager) Init(dataSourcesPath string) error {
	ds, err := manager.readDataSources(dataSourcesPath)
	if err != nil {
		return err
	}
	manager.DataSources = ds
	return nil
}

func (manager Manager) readDataSources(path string) ([]DataSource, error) {
	var ds []DataSource
	if path == "" {
		return ds, errors.New("path to data source was not specified")
	}
	file, err := os.Open(path)
	if err != nil {
		return ds, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ds)
	if err != nil {
		return ds, err
	}
	return ds, nil
}

func (manager Manager) StartCrawling(req SearchRequest) error {
	fieldsPart := map[string][]string{}
	var dataSource DataSource
	for _, ds := range manager.DataSources {
		if ds.Name == req.SourceName {
			dataSource = ds
			break
		}
	}
	if dataSource.Name == "" {
		return errors.New("incorrect data source name specified")
	}
	for key, value := range req.Fields {
		checkDs := false
		for _, dsField := range dataSource.Keys {
			if key == dsField {
				checkDs = true
				break
			}
		}
		if !checkDs {
			return errors.New("key " + key + " was not found in data source " + dataSource.Name)
		}
		fieldsPart[key] = []string{}
		setParts := strings.Split(value, ",")
		if len(setParts) > 1 {
			fieldsPart[key] = setParts
		} else {
			rangeParts := strings.Split(value, "-")
			if len(rangeParts) != 2 {
				return errors.New("search range for key " + key + " was specified incorrectly: " + value)
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return err
			}
			finish, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return err
			}
			if start > finish {
				return errors.New("range error for key " + key + ": start value must be less or equal than finish value")
			}
			rangeSlice := make([]string, finish-start+1)
			for i := range rangeSlice {
				rangeSlice[i] = strconv.Itoa(start + i)
			}
			fieldsPart[key] = rangeSlice
		}
	}
}
