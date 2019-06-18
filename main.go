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

	templates, err = templates.ParseFiles(fileName)
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

	forecast := data["forecasts"].([]interface{})
	//fact := data["fact"].(map[string]interface{})

	todayForecast := forecast[0].(map[string]interface{})
	todayParts := todayForecast["parts"].(map[string]interface{})
	todayMorning := todayParts["morning"].(map[string]interface{})
	todayDay := todayParts["day"]
	todayEvening := todayParts["evening"]
	fmt.Println(todayDay, todayEvening)

	//for _, value := range fact {
	//	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(value)
	//}

	//var text string = "Сегодня утром будет " + conditionTranslate(todayMorning["condition"].(string))

	fmt.Println(getTemplate("day_forecast_show.gohtml", todayMorning))
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
