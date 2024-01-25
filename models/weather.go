package models

import (
	"encoding/json"
	"fmt"
	"github.com/daniel-z-johnson/peronalWeatherSite/config"
	"log/slog"
	"net/http"
	"net/url"
)

type GeoLocation struct {
	Version    int32  `json:"Version"`
	Key        string `json:"Key"`
	Type       string `json:"Type"`
	ParentCity struct {
		Key           string `json:"Key"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"ParentCity"`
}

type WeatherAPI struct {
	config *config.Config
	logger *slog.Logger
}

func WeatherService(logger *slog.Logger, config *config.Config) *WeatherAPI {
	return &WeatherAPI{logger: logger, config: config}
}

func (wa *WeatherAPI) GetGeoPoints() ([]*GeoLocation, error) {
	geoPoints := make([]*GeoLocation, 0)
	for _, zip := range wa.config.Zipcodes {
		u, err := url.Parse("https://dataservice.accuweather.com/locations/v1/postalcodes/search")
		if zip.CountryCode != "" {
			u, err = url.Parse(
				fmt.Sprintf("https://dataservice.accuweather.com/locations/v1/postalcodes/%s/search", zip.CountryCode))
		}
		if err != nil {
			return nil, err
		}
		values := u.Query()
		values.Set("q", zip.PostalCode)
		values.Set("apikey", wa.config.WeatherAPI.Key)
		u.RawQuery = values.Encode()
		wa.logger.Info(u.String())
		geoLoc := make([]*GeoLocation, 0)
		r, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		err = json.NewDecoder(r.Body).Decode(&geoLoc)
		if err != nil {
			return nil, err
		}
		geoPoints = append(geoPoints, geoLoc[0])
	}
	return geoPoints, nil
}