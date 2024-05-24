package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	OWAPIKey string `yaml:"OWapiKey"`
	DDAPIKey string `yaml:"DDapiKey"`
}

type WeatherData struct {
	List []City `json:"list"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

type City struct {
	ID    int          `json:"id"`
	Name  string       `json:"name"`
	Main  *Main        `json:"main"`
	Coord *Coordinates `json:"coord"`
}

type Coordinates struct {
	Latitude float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type Dadata struct {
	Suggestion *Suggestions `json:"suggestions"` 
}

type Suggestions struct {
	Data *Data `json:"data"`
}

type Data struct {
	Region string `json:"region_with_type"`
}

var config Config

func main() {
	if err := cleanenv.ReadConfig("config.yaml", &config); err != nil {
		fmt.Println("Ошибка при загрузке конфигурации:", err)
		return
	}

	var town string
	fmt.Scan(&town)

	urlOW := fmt.Sprintf("https://api.openweathermap.org/data/2.5/find?q=%v&type=like&APPID=%s&units=metric&cnt=15", town, config.OWAPIKey)

	respOW, err := http.Get(urlOW)
	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}
	defer respOW.Body.Close()

	bodyOW, err := io.ReadAll(respOW.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	var data WeatherData
	err = json.Unmarshal(bodyOW, &data)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	var lat, lon float64

	for _, city := range data.List {
		fmt.Printf("Current temperature °C in %s: %.0f°C\n", town, city.Main.Temp)
		fmt.Printf("Feels like: %.0f°C\n", city.Main.FeelsLike)
		fmt.Printf("Current temperature °F in %v: %.0f°F\n", town, float64(((city.Main.Temp)*(9/5)))+32)
		fmt.Printf("Feels like: %.0f°F\n", float64(((city.Main.FeelsLike)*(9/5)))+32)
		fmt.Printf("Latitude: %.4f\n", city.Coord.Latitude)
		fmt.Printf("Longitude: %.4f\n", city.Coord.Longitude)
		lat = city.Coord.Latitude
		lon = city.Coord.Longitude
	}
	postRequest(lat, lon)
}

func postRequest(latitude, longitude float64) {

	str := fmt.Sprintf(`{ "lat": %f, "lon": %f, "count": 1 }`, latitude, longitude)
	// dataJs := map[string]float64{
	// 	"lat": city.Coord.Lattitude,
	// 	"lon": city.Coord.Longitude,
	// 	"count": 1
	// }

	jsonData, err := json.Marshal(str)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	urlDD := fmt.Sprintf("http://suggestions.dadata.ru/suggestions/%v/4_1/rs/geolocate/address",config.DDAPIKey )
	respDD, err := http.Post(urlDD, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при запросе Post:", err)
		return
	}

	bodyDD, err := io.ReadAll(respDD.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	if respDD.StatusCode == http.StatusOK {
		fmt.Println("Успешный ответ от сервера:")
	} else {
		fmt.Printf("Ошибка: сервер вернул статус %d\n", respDD.StatusCode)
	}

	var cou Dadata
	err = json.Unmarshal(bodyDD, &cou)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}
	
	fmt.Println(cou)
}