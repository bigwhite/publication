package doubles

import "fmt"

type WeatherProvider interface {
	GetTemperature(city string) (int, error)
}

func GetWeatherAdvice(provider WeatherProvider, city string) string {
	temp, err := provider.GetTemperature(city)
	if err != nil {
		return fmt.Sprintf("Sorry, could not get weather for %s: %v", city, err)
	}
	if temp > 25 {
		return fmt.Sprintf("It's hot in %s (%d°C)! Stay hydrated.", city, temp)
	} else if temp < 10 {
		return fmt.Sprintf("It's cold in %s (%d°C)! Dress warmly.", city, temp)
	}
	return fmt.Sprintf("The weather in %s is pleasant (%d°C).", city, temp)
}
