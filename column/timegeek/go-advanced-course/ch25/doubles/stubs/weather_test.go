package doubles

import (
	"fmt"
	"testing"
)

// --- Stub Implementation ---
type WeatherProviderStub struct {
	FixedTemperature int
	FixedError       error
	CityCalled       string // Optional: to verify if the correct city was passed
}

func (s *WeatherProviderStub) GetTemperature(city string) (int, error) {
	s.CityCalled = city // Record the city that was called
	return s.FixedTemperature, s.FixedError
}

func TestGetWeatherAdvice_WithStub(t *testing.T) {
	t.Run("HotWeather", func(t *testing.T) {
		stub := &WeatherProviderStub{FixedTemperature: 30, FixedError: nil}
		advice := GetWeatherAdvice(stub, "Dubai")
		expected := "It's hot in Dubai (30°C)! Stay hydrated."
		if advice != expected {
			t.Errorf("Expected '%s', got '%s'", expected, advice)
		}
		if stub.CityCalled != "Dubai" {
			t.Errorf("Expected GetTemperature to be called with 'Dubai', got '%s'", stub.CityCalled)
		}
	})

	t.Run("ProviderError", func(t *testing.T) {
		stub := &WeatherProviderStub{FixedError: fmt.Errorf("API quota exceeded")}
		advice := GetWeatherAdvice(stub, "London")
		expected := "Sorry, could not get weather for London: API quota exceeded"
		if advice != expected {
			t.Errorf("Expected '%s', got '%s'", expected, advice)
		}
	})
}
