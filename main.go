package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getTemplate(fileName string, data interface{}) (err error, result string) {
	templates, err := template.ParseGlob("templates/*.gohtml")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
	}

	templates, err = templates.ParseFiles("templates/" + fileName)
	if err != nil {
		return
	}

	var tpl bytes.Buffer
	if err := templates.Execute(&tpl, data); err != nil {
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
		now     map[string]interface{}
		day     map[string]interface{}
		evening map[string]interface{}
	}

	templateData := Forecast{
		now:     fact,
		day:     todayDay,
		evening: todayEvening,
	}

	fmt.Println(templateData.now["temp"])
	fmt.Println(getTemplate("day_forecast_show.gohtml", templateData))
	//sendToHorn(text)
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
