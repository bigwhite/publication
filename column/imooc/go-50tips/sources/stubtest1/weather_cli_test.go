package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var weatherResp = []Weather{
	{
		City:    "nanning",
		TemP:    "26~33",
		Weather: "rain",
		Date:    "05-04",
	},
	{
		City:    "guiyang",
		TemP:    "25~29",
		Weather: "sunny",
		Date:    "05-04",
	},
	{
		City:    "tianjin",
		TemP:    "20~31",
		Weather: "windy",
		Date:    "05-04",
	},
}

func TestGetWeatherInfoOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte

		if r.URL.EscapedPath() != "/weather" {
			w.WriteHeader(http.StatusForbidden)
		}

		r.ParseForm()
		city := r.Form.Get("city")
		if city == "guiyang" {
			data, _ = json.Marshal(&weatherResp[1])
		}
		if city == "tianjin" {
			data, _ = json.Marshal(&weatherResp[2])
		}
		if city == "nanning" {
			data, _ = json.Marshal(&weatherResp[0])
		}

		w.Write(data)
	}))
	defer ts.Close()
	addr := ts.URL

	city := "guiyang"
	w, err := GetWeatherInfo(addr, city)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if w.City != city {
		t.Errorf("want %s, got %s", city, w.City)
	}
	if w.Weather != "sunny" {
		t.Errorf("want %s, got %s", "sunny", w.City)
	}
}
