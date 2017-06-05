package weather

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	"time"
)

type Forecast struct {
	Dt int64 `json:"dt"`
	Main Metrics `json:"main"`
	Weather []Weather `json:"weather"`
	Clouds Clouds `json:"clouds"`
	Wind Wind `json:"wind"`
	Rain Precipitations `json:"rain"`
	Snow Precipitations `json:"snow"`
	Sys Sys `json:"sys"`
	DtTxt string `json:"dt_txt"`
	Date time.Time `json:"-"`
}

type Metrics struct {
	Temp float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
	Pressure float64 `json:"pressure"`
	SeaLevel float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity float64 `json:"humidity"`
	TempKf float64 `json:"temp_kf"`
}

type Weather struct {
	ID int `json:"id"`
	Main string `json:"main"`
	Description string `json:"description"`
	Icon string `json:"icon"`
}

type Clouds struct {
	All float64 `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg float64 `json:"deg"`
}

type Precipitations struct {
	ThreeHours float64 `json:"3h"`
}

type Sys struct {
	Pod string `json:"pod"`
}

type City struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Coord Coordinates `json:"coord"`
	Country string `json:"country"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lon"`
}

type ForecastResponse struct {
	Cod string `json:"cod"`
	Message float64 `json:"message"`
	Cnt int `json:"cnt"`
	List []Forecast `json:"list"`
	City City `json:"city"`
}


func (c *Client) ForecastByID(cityID int64) ([]Forecast, error) {
	var forecastResp ForecastResponse

	req, err := http.NewRequest("GET", c.genForecastByIDURL(cityID), nil)
	if err != nil {
		return nil, err
	}
	<-c.throttle
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forcast returned errorneous response code: %d", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&forecastResp)
	if err != nil {
		return nil, err
	}

	err = parseTimes(forecastResp.List)

	return forecastResp.List, err
}

func (c *Client) genForecastByIDURL(cityID int64) string {
	return "http://api.openweathermap.org/data/2.5/forecast?APPID="+
		c.apiKey+"&units=metric&id="+strconv.FormatInt(cityID, 10)
}

func parseTimes(forecasts []Forecast) error {
	var err error

	for i, forecast := range forecasts {
		forecast, err = parseTime(forecast)
		if err != nil {
			return err
		}

		forecasts[i] = forecast
	}

	return nil
}

func parseTime(forecast Forecast) (Forecast, error) {
	t, err := time.Parse("2006-01-02 15:04:05", forecast.DtTxt)
	if err != nil {
		return forecast, err
	}

	forecast.Date = t

	return forecast, nil
}