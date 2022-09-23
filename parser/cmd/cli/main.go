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

type LocationRaw struct {
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

type LocationJson struct {
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

type Page struct {
	Dev1  []LocationJson `json:"dev1"`
	Dev2  []LocationJson `json:"dev2"`
	Dev3  []LocationJson `json:"dev3"`
	Dev4  []LocationJson `json:"dev4"`
	Dev5  []LocationJson `json:"dev5"`
	Dev6  []LocationJson `json:"dev6"`
	Dev7  []LocationJson `json:"dev7"`
	Dev8  []LocationJson `json:"dev8"`
	Dev9  []LocationJson `json:"dev9"`
	Dev10 []LocationJson `json:"dev10"`
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

	devs := make([][]LocationRaw, 10)

	for i, val := range results {
		var rows []LocationRaw
		for _, v := range val {
			result, _ := json.Marshal(v)

			var row LocationRaw

			if err := json.Unmarshal(result, &row); err != nil {
				panic("error unmarshaling")
			}

			rows = append(rows, row)
		}
		devs[i] = rows
	}

	pages := Page{
		Dev1:  make([]LocationJson, 0),
		Dev2:  make([]LocationJson, 0),
		Dev3:  make([]LocationJson, 0),
		Dev4:  make([]LocationJson, 0),
		Dev5:  make([]LocationJson, 0),
		Dev6:  make([]LocationJson, 0),
		Dev7:  make([]LocationJson, 0),
		Dev8:  make([]LocationJson, 0),
		Dev9:  make([]LocationJson, 0),
		Dev10: make([]LocationJson, 0),
	}
	wg := &sync.WaitGroup{}
	wg.Add(10)
	go helper(devs[0], &pages.Dev1, wg)
	go helper(devs[1], &pages.Dev2, wg)
	go helper(devs[2], &pages.Dev3, wg)
	go helper(devs[3], &pages.Dev4, wg)
	go helper(devs[4], &pages.Dev5, wg)
	go helper(devs[5], &pages.Dev6, wg)
	go helper(devs[6], &pages.Dev7, wg)
	go helper(devs[7], &pages.Dev8, wg)
	go helper(devs[8], &pages.Dev9, wg)
	go helper(devs[9], &pages.Dev10, wg)
	wg.Wait()

	bres, err := json.Marshal(&pages)
	if err != nil {
		panic("hi it is panic, i love panics")
	}

	if err := os.WriteFile("pages.json", bres, 0644); err != nil {
		fmt.Println("errored due to: ", err.Error())
	}
}

func helper(d []LocationRaw, dev *[]LocationJson, wg *sync.WaitGroup) {
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

		temp := LocationJson{
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
