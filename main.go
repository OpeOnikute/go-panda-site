package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/opeonikute/panda"
)

type HomePageData struct {
	PandaFound bool
	ImageURL   string
	Source     string
	FileName   string
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := HomePageData{
			PandaFound: false,
		}

		panda, err := getPandaOfTheDay()

		if err != nil {
			fmt.Print(err)
		} else {
			data.PandaFound = true
			data.ImageURL = panda.URL
			data.Source = panda.Source
			data.FileName = panda.FileName
		}

		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.Execute(w, data)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Server listening on port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func getPandaOfTheDay() (panda.Entry, error) {
	var en panda.Entry
	goPanda := panda.GoPanda{
		Config: panda.Settings{
			MongoURL: os.Getenv("MONGO_URL"),
			MongoDB:  os.Getenv("MONGO_DATABASE"),
		},
	}

	tm := time.Now()
	res, err := goPanda.GetPOD(tm)

	if err != nil {
		return en, err
	}

	return res, nil
}
