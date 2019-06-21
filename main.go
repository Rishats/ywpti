package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ivahaev/russian-time"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func getTemplate(fileName string, funcmap template.FuncMap, data interface{}) (result string, err error) {
	template, err := template.New(fileName).Funcs(funcmap).ParseFiles("templates/" + fileName)
	if err != nil {
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		panic(err)
	}

	result = tpl.String()

	return
}

func apiData() string {
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

func sendToHorn(text string) {
	m := map[string]interface{}{
		"text": text,
	}
	mJson, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJson)
	req, _ := http.NewRequest("POST", os.Getenv("INTEGRAM_WEBHOOK_URI"), contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)

	fmt.Println(resp)
}

func conditionTranslate(condition string) string {
	conditions := map[string]string{
		"clear":                            "Ясно ☀️️",
		"partly-cloudy":                    "Малооблачно ⛅",
		"cloudy":                           "Облачно с прояснениями ⛅️",
		"overcast":                         "Пасмурно 🌁",
		"partly-cloudy-and-light-rain":     "Небольшой дождь ☂️",
		"partly-cloudy-and-rain":           "Дождь ☔",
		"overcast-and-rain":                "Сильный дождь 🌧️",
		"overcast-thunderstorms-with-rain": "Сильный дождь, гроза ⛈️",
		"cloudy-and-light-rain":            "Небольшой дождь 🌧",
		"overcast-and-light-rain":          "Небольшой дождь 🌧",
		"cloudy-and-rain":                  "Дождь 🌧️",
		"overcast-and-wet-snow":            "Дождь со снегом 🌧❄",
		"partly-cloudy-and-light-snow":     "Небольшой снег ❄",
		"partly-cloudy-and-snow":           "Снег ❄️",
		"overcast-and-snow":                "Снегопад 🌨️ ❄️❄️❄️",
		"cloudy-and-light-snow":            "Небольшой снег 🌨️",
		"overcast-and-light-snow":          "Небольшой снег 🌨️",
		"cloudy-and-snow":                  "Снег ❄️"}

	return conditions[condition]
}

func windDirTranslate(windDir string) string {
	windDirs := map[string]string{
		"nw": "северо-западное",
		"n":  "северное",
		"ne": "северо-восточное",
		"e":  "восточное",
		"se": "юго-восточное",
		"s":  "южное",
		"sw": "юго-западное",
		"w":  "западное",
		"c":  "штиль"}

	return windDirs[windDir]
}

func hourWithMin() string {

	timeStamp := time.Unix(time.Now().Unix(), 0)

	hr, min, _ := timeStamp.Clock()

	finalTime := "%d:%d"

	result := fmt.Sprintf(finalTime, hr, min)

	return result
}

func weekDay() rtime.Weekday {
	t := rtime.Now()
	standardTime := time.Now()
	t = rtime.Time(standardTime)

	return t.Weekday()
}

func dayForecastShow() {
	dataJson := apiData()
	var data map[string]interface{}
	json.Unmarshal([]byte(dataJson), &data)

	forecast := data["forecast"].(map[string]interface{})
	fact := data["fact"].(map[string]interface{})

	todayParts := forecast["parts"].([]interface{})

	var todayDay map[string]interface{}
	var todayEvening map[string]interface{}

	for _, value := range todayParts {
		switch value.(map[string]interface{})["part_name"] {
		case "day":
			todayDay = value.(map[string]interface{})
		case "evening":
			todayEvening = value.(map[string]interface{})
		}
	}

	type Forecast struct {
		Now     map[string]interface{}
		Day     map[string]interface{}
		Evening map[string]interface{}
	}

	templateData := Forecast{
		Now:     fact,
		Day:     todayDay,
		Evening: todayEvening,
	}

	funcmap := template.FuncMap{
		"conditionTranslate": conditionTranslate,
		"windDirTranslate":   windDirTranslate,
		"weekDay":            weekDay,
		"hourWithMin":        hourWithMin,
	}

	text, err := getTemplate("day_forecast_show.gohtml", funcmap, templateData)
	if err != nil {
		log.Fatal(err)
	}
	sendToHorn(text)
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dayForecastShow()

	tasks()
}
