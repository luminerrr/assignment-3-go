package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"log"
	"math/rand"
)

type Data struct {
	Status Status `json:"status"`
	Condition string `json:"condition"`
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func updateData() {
	for {
		var data = Data{}
		
		maxVal := 30

		data.Status.Water = rand.Intn(maxVal)
		data.Status.Wind = rand.Intn(maxVal)
		if data.Status.Water <5 || data.Status.Wind <6 {
			data.Condition = "Aman"
		} else if data.Status.Water < 9 || data.Status.Water > 5 || data.Status.Wind < 16 || data.Status.Wind > 6{
			data.Condition = "Siaga"
		} else if  data.Status.Water > 8 || data.Status.Wind > 15 {
			data.Condition = "Bahaya"
		}
		

		b, err := json.MarshalIndent(&data, "", " ")
		if err != nil {
			log.Fatal("error occurred when marshalling data",err)
		}

		err = ioutil.WriteFile("data.json", b, 0644)
		if err != nil {
			log.Fatal("error occurred when write file to data.json", err)
		}

		time.Sleep(time.Second * 15)
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal("error while reading html file", err)
		}

		data := Data{Status: Status{}}
		b,err := ioutil.ReadFile("data.json")
		if err != nil {
			log.Fatal("error while loading data on middleware", err)
		}

		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Fatal("error while unmarshalling", err)
		}

		tpl.ExecuteTemplate(w, "index.html", data)

	})

	http.ListenAndServe(":8080", nil)
}