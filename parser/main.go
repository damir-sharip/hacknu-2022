package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type DataString struct {
	Latitude                     string `json:"latitude"`
	Longitude                    string `json:"longitude"`
	Altitude                     string `json:"altitude"`
	Identifier                   string `json:"identifier"`
	Timestamp                    string `json:"timestamp"`
	FloorLabel                   string `json:"floor_label"`
	HorizontalAccuracy           string `json:"horizontal_accuracy"`
	VerticalAccuracy             string `json:"vertical_accuracy"`
	ConfidenceInLocationAccuracy string `json:"confidence_in_location_accuracy"`
	Activity                     string `json:"activity"`
}

type DataTyped struct {
	Latitude                     float64 `json:"latitude"`
	Longitude                    float64 `json:"longitude"`
	Altitude                     float64 `json:"altitude"`
	Identifier                   *string `json:"identifier"`
	Timestamp                    int64   `json:"timestamp"`
	FloorLabel                   *string `json:"floor_label"`
	HorizontalAccuracy           float64 `json:"horizontal_accuracy"`
	VerticalAccuracy             float64 `json:"vertical_accuracy"`
	ConfidenceInLocationAccuracy float64 `json:"confidence_in_location_accuracy"`
	Activity                     *string `json:"activity"`
}

type Data struct {
	Dev1  []DataString `json:"dev1"`
	Dev2  []DataString `json:"dev2"`
	Dev3  []DataString `json:"dev3"`
	Dev4  []DataString `json:"dev4"`
	Dev5  []DataString `json:"dev5"`
	Dev6  []DataString `json:"dev6"`
	Dev7  []DataString `json:"dev7"`
	Dev8  []DataString `json:"dev8"`
	Dev9  []DataString `json:"dev9"`
	Dev10 []DataString `json:"dev10"`
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
	bdata, err := os.ReadFile("data.json")
	if err != nil {
		panic("err no file: " + err.Error())
	}
	var d Data

	if err := json.Unmarshal(bdata, &d); err != nil {
		panic("error unmarshaling")
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
	go helper(d.Dev1, &res.Dev1, wg)
	go helper(d.Dev2, &res.Dev2, wg)
	go helper(d.Dev3, &res.Dev3, wg)
	go helper(d.Dev4, &res.Dev4, wg)
	go helper(d.Dev5, &res.Dev5, wg)
	go helper(d.Dev6, &res.Dev6, wg)
	go helper(d.Dev7, &res.Dev7, wg)
	go helper(d.Dev8, &res.Dev8, wg)
	go helper(d.Dev9, &res.Dev9, wg)
	go helper(d.Dev10, &res.Dev10, wg)
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
	for _, v := range d {
		latitude, _ := strconv.ParseFloat(v.Latitude, 64)
		longitude, _ := strconv.ParseFloat(v.Longitude, 64)
		altitude, _ := strconv.ParseFloat(v.Altitude, 64)
		ts, _ := strconv.ParseInt(v.Timestamp, 10, 64)
		ha, _ := strconv.ParseFloat(v.HorizontalAccuracy, 64)
		va, _ := strconv.ParseFloat(v.VerticalAccuracy, 64)

		ca, _ := strconv.ParseFloat(v.ConfidenceInLocationAccuracy, 64)

		var identifier *string
		if v.Identifier != "null" {
			identifier = &v.Identifier
		}

		var floorLabel *string
		if v.FloorLabel != "null" {
			floorLabel = &v.FloorLabel
		}

		var activity *string
		if v.Activity != "null" {
			activity = &v.Activity
		}

		temp := DataTyped{
			Latitude:                     latitude,
			Longitude:                    longitude,
			Altitude:                     altitude,
			Identifier:                   identifier,
			Timestamp:                    ts,
			FloorLabel:                   floorLabel,
			HorizontalAccuracy:           ha,
			VerticalAccuracy:             va,
			ConfidenceInLocationAccuracy: ca,
			Activity:                     activity,
		}

		*dev = append(*dev, temp)
	}
}
