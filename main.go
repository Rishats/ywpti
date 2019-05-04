package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func apiData() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("YW_API_KEY")
	apiUri := os.Getenv("YW_API_URI") + "?lat=" + os.Getenv("YW_LAT") + "&lon=" + os.Getenv("YW_LON") + "&lang=" + os.Getenv("YW_LANG")

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiUri, nil)

	req.Header.Add("X-Yandex-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body)
}

func conditionTranslate(condition string) string {
	conditions := map[string]string{
		"clear":                            "ясно",
		"partly-cloudy":                    "малооблачно",
		"cloudy":                           "облачно с прояснениями",
		"overcast":                         "пасмурно",
		"partly-cloudy-and-light-rain":     "небольшой дождь",
		"partly-cloudy-and-rain":           "дождь",
		"overcast-and-rain":                "сильный дождь",
		"overcast-thunderstorms-with-rain": "сильный дождь, гроза",
		"cloudy-and-light-rain":            "небольшой дождь",
		"overcast-and-light-rain":          "небольшой дождь",
		"cloudy-and-rain":                  "дождь",
		"overcast-and-wet-snow":            "дождь со снегом",
		"partly-cloudy-and-light-snow":     "небольшой снег",
		"partly-cloudy-and-snow":           "снег",
		"overcast-and-snow":                "снегопад",
		"cloudy-and-light-snow":            "небольшой снег",
		"overcast-and-light-snow":          "небольшой снег",
		"cloudy-and-snow":                  "снег"}

	return conditions[condition]
}

func dayForecastShow() {
	fmt.Println("I will show forecast data.")
	apiData()
}

func tasks() {
	gocron.Every(1).Monday().At("9:00").Do(dayForecastShow)
	gocron.Every(1).Tuesday().At("9:00").Do(dayForecastShow)
	gocron.Every(1).Wednesday().At("9:00").Do(dayForecastShow)
	gocron.Every(1).Thursday().At("9:00").Do(dayForecastShow)
	gocron.Every(1).Friday().At("9:00").Do(dayForecastShow)

	// remove, clear and next_run
	_, time := gocron.NextRun()
	fmt.Println(time)

	// function Start start all the pending jobs
	<-gocron.Start()
}

func main() {
	tasks()
}
