package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	APIKey string `yaml:"apiKey"`
}

type WeatherData struct {
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
	} `json:"main"`
}

func main() {
	var config Config
	if err := cleanenv.ReadConfig("config.yaml", &config); err != nil {
		fmt.Println("Ошибка при загрузке конфигурации:", err)
		return
	}

	var city string
	fmt.Scan(&city)
	 
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, config.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	var data WeatherData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	fmt.Printf("Current temperature °C in %s: %.0f°C\n" + "Feels like: %.0f°C\n" + "Current temperature °F in %v: %.0f°F\n" + 
	"Feels like: %.0f°F", city, data.Main.Temp-271, data.Main.Feels_like-271, city, float64(((data.Main.Temp-271)*(9/5))) +
	32, float64(((data.Main.Feels_like-271)*(9/5))) + 32)
}
