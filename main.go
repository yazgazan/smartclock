package main

import (
	"github.com/yazgazan/smartclock/weather"
	"log"
	"fmt"
	"time"
)

const (
	AmsterdamCityID = 2759794
	APIKey = ""
)

func main() {
	c := weather.NewClient(APIKey)
	defer c.Close()

	forecasts, err := c.ForecastByID(AmsterdamCityID)
	if err != nil {
		log.Fatalln(err)
	}

	forecasts = todaysForecasts(forecasts)
	forecasts = dayTimeForecasts(forecasts)
	fmt.Printf("AverageTemp: %f\n", averageTemp(forecasts))
	fmt.Printf("AverageHumidity: %f\n", averageHumidity(forecasts))
	fmt.Printf("AverageCloudiness: %f\n", averageCloudiness(forecasts))
}

func todaysForecasts(forecasts []weather.Forecast) []weather.Forecast {
	today := make([]weather.Forecast, 0, len(forecasts))
	date := time.Now()

	for _, forecast := range forecasts {
		if forecast.Date.Year() != date.Year() {
			continue
		}
		if forecast.Date.Month() != date.Month() {
			continue
		}
		if forecast.Date.Day() != date.Day() {
			continue
		}
		today = append(today, forecast)
	}

	return today
}

func dayTimeForecasts(forecasts []weather.Forecast) []weather.Forecast {
	daytimeForecasts := make([]weather.Forecast, 0, len(forecasts))

	for _, forecast := range forecasts {
		if forecast.Date.Hour() < 8 {
			continue
		}
		if forecast.Date.Hour() > 22 {
			continue
		}
		daytimeForecasts = append(daytimeForecasts, forecast)
	}

	return daytimeForecasts
}

func averageTemp(forecasts []weather.Forecast) float64 {
	temp := 0.0

	for _, forecast := range forecasts {
		temp += forecast.Main.Temp / float64(len(forecasts))
	}

	return temp
}

func averageHumidity(forecasts []weather.Forecast) float64 {
	humidity := 0.0

	for _, forecast := range forecasts {
		humidity += forecast.Main.Humidity / float64(len(forecasts))
	}

	return humidity
}

func averageCloudiness(forecasts []weather.Forecast) float64 {
	cloudiness := 0.0

	for _, forecast := range forecasts {
		cloudiness += forecast.Clouds.All / float64(len(forecasts))
	}

	return cloudiness
}