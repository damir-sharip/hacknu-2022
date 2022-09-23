package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/FerdinaKusumah/excel2json"
)

type DataString struct {
	Latitude                     string `json:"latitude"`
	Longitude                    string `json:"longitude"`
	Altitude                     string `json:"altitude"`
	Identifier                   string `json:"identifier"`
	Timestamp                    string `json:"timestamp"`
	FloorLabel                   string `json:"floorlabel"`
	HorizontalAccuracy           string `json:"horizontalaccuracy"`
	VerticalAccuracy             string `json:"verticalaccuracy"`
	ConfidenceInLocationAccuracy string `json:"confidenceinlocationaccuracy"`
	Activity                     string `json:"activity"`
}

type DataTyped struct {
	Latitude                     float64 `json:"latitude"`
	Longitude                    float64 `json:"longitude"`
	Altitude                     float64 `json:"altitude"`
	Identifier                   *string `json:"identifier"`
	Timestamp                    float64 `json:"timestamp"`
	FloorLabel                   *string `json:"floorLabel"`
	HorizontalAccuracy           float64 `json:"horizontalAccuracy"`
	VerticalAccuracy             float64 `json:"verticalAccuracy"`
	ConfidenceInLocationAccuracy float64 `json:"confidenceInLocationAccuracy"`
	Activity                     *string `json:"activity"`
}

type DataT struct {
	Dev1  []DataTyped `json:"dev1"`
	Dev2  []DataTyped `json:"dev2"`
	Dev3  []DataTyped `json:"dev3"`
	Dev4  []DataTyped `json:"dev4"`
	Dev5  []DataTyped `json:"dev5"`
	Dev6  []DataTyped `json:"dev6"`
	Dev7  []DataTyped `json:"dev7"`
	Dev8  []DataTyped `json:"dev8"`
	Dev9  []DataTyped `json:"dev9"`
	Dev10 []DataTyped `json:"dev10"`
}

func main() {
	var (
		result    []*map[string]interface{}
		err       error
		path      = "./hacknu-dev-data.xlsx"
		sheetName = []string{"dev1", "dev2", "dev3", "dev4", "dev5", "dev6", "dev7", "dev8", "dev9", "dev10"}
		headers   []string
	)

	var results [][]*map[string]interface{}
	for i := 0; i < 10; i++ {
		if result, err = excel2json.GetExcelFilePath(path, sheetName[i], headers); err != nil {
			log.Fatalf(`unable to parse file, error: %s`, err)
		}
		results = append(results, result)
	}

	rawses := make([][]DataString, 10)

	for i, val := range results {
		var raws []DataString
		for _, v := range val {
			result, _ := json.Marshal(v)

			var data DataString

			if err := json.Unmarshal(result, &data); err != nil {
				panic("error unmarshaling")
			}

			raws = append(raws, data)
		}
		rawses[i] = raws
	}

	res := DataT{
		Dev1:  make([]DataTyped, 0),
		Dev2:  make([]DataTyped, 0),
		Dev3:  make([]DataTyped, 0),
		Dev4:  make([]DataTyped, 0),
		Dev5:  make([]DataTyped, 0),
		Dev6:  make([]DataTyped, 0),
		Dev7:  make([]DataTyped, 0),
		Dev8:  make([]DataTyped, 0),
		Dev9:  make([]DataTyped, 0),
		Dev10: make([]DataTyped, 0),
	}
	wg := &sync.WaitGroup{}
	wg.Add(10)
	go helper(rawses[0], &res.Dev1, wg)
	go helper(rawses[1], &res.Dev2, wg)
	go helper(rawses[2], &res.Dev3, wg)
	go helper(rawses[3], &res.Dev4, wg)
	go helper(rawses[4], &res.Dev5, wg)
	go helper(rawses[5], &res.Dev6, wg)
	go helper(rawses[6], &res.Dev7, wg)
	go helper(rawses[7], &res.Dev8, wg)
	go helper(rawses[8], &res.Dev9, wg)
	go helper(rawses[9], &res.Dev10, wg)
	wg.Wait()

	bres, err := json.Marshal(&res)
	if err != nil {
		panic("hi it is panic, i love panics")
	}

	if err := os.WriteFile("res.json", bres, 0644); err != nil {
		fmt.Println("errored due to: ", err.Error())
	}
}

func helper(d []DataString, dev *[]DataTyped, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range d {
		latitude, _ := strconv.ParseFloat(d[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(d[i].Longitude, 64)
		altitude, _ := strconv.ParseFloat(d[i].Altitude, 64)
		ts, _ := strconv.ParseFloat(d[i].Timestamp, 64)
		ha, _ := strconv.ParseFloat(d[i].HorizontalAccuracy, 64)
		va, _ := strconv.ParseFloat(d[i].VerticalAccuracy, 64)

		ca, _ := strconv.ParseFloat(d[i].ConfidenceInLocationAccuracy, 64)

		var identifier *string
		if d[i].Identifier != "null" {
			identifier = &d[i].Identifier
		}

		var activity *string
		if d[i].Activity != "null" {
			activity = &d[i].Activity
		}

		var floorlable *string
		if d[i].FloorLabel != "null" {
			floorlable = &d[i].FloorLabel
		}

		temp := DataTyped{
			Latitude:                     latitude,
			Longitude:                    longitude,
			Altitude:                     altitude,
			Identifier:                   identifier,
			FloorLabel:                   floorlable,
			Timestamp:                    ts,
			HorizontalAccuracy:           ha,
			VerticalAccuracy:             va,
			ConfidenceInLocationAccuracy: ca,
			Activity:                     activity,
		}

		*dev = append(*dev, temp)
	}
}
