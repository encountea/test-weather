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
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
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

	fmt.Printf("Current temperature °C in %s: %.0f°C\n", city, data.Main.Temp-271)
	fmt.Printf("Feels like: %.0f°C\n", data.Main.FeelsLike-271)
	fmt.Printf("Current temperature °F in %v: %.0f°F\n", city, float64(((data.Main.Temp-271)*(9/5)))+32)
	fmt.Printf("Feels like: %.0f°F", float64(((data.Main.FeelsLike-271)*(9/5)))+32)
}
