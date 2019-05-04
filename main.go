package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("YW_API_KEY")

	client := &http.Client{
	}

	req, err := http.NewRequest("GET", "https://api.weather.yandex.ru/v1/informers?lat=55.75396&lon=37.620393", nil)

	req.Header.Add("X-Yandex-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(apiKey)
	fmt.Println(string(body))
}