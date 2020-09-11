package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Weather struct {
	City    string `json:"city"`
	Date    string `json:"date"`
	TemP    string `json:"temP"`
	Weather string `json:"weather"`
}

func GetWeatherInfo(addr string, city string) (*Weather, error) {
	url := fmt.Sprintf("%s/weather?city=%s", addr, city)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code is %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w Weather
	err = json.Unmarshal(body, &w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
