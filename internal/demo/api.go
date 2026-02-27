package demo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func FetchISS(ctx context.Context, client *http.Client) (float64, float64, bool) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.wheretheiss.at/v1/satellites/25544", nil)
	if err != nil {
		return 0, 0, false
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, 0, false
	}
	var payload struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, 0, false
	}
	return payload.Latitude, payload.Longitude, true
}

func FetchLocalGeo(ctx context.Context, client *http.Client) (*IPGeoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipapi.co/json/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipapi status: %d", resp.StatusCode)
	}
	var geo IPGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return nil, err
	}
	if geo.City == "" {
		geo.City = "Unknown"
	}
	return &geo, nil
}

func NextSunrise(ctx context.Context, client *http.Client, geo *IPGeoResponse, now time.Time) (*time.Time, error) {
	if geo == nil {
		return nil, nil
	}
	firstDate := now.UTC().Format("2006-01-02")
	first, err := fetchSunriseForDate(ctx, client, geo.Lat, geo.Lon, firstDate)
	if err != nil {
		return nil, err
	}
	if first.After(now) {
		return &first, nil
	}
	secondDate := now.UTC().Add(24 * time.Hour).Format("2006-01-02")
	second, err := fetchSunriseForDate(ctx, client, geo.Lat, geo.Lon, secondDate)
	if err != nil {
		return nil, err
	}
	return &second, nil
}

func fetchSunriseForDate(ctx context.Context, client *http.Client, lat, lon float64, date string) (time.Time, error) {
	url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%.4f&lng=%.4f&date=%s&formatted=0", lat, lon, date)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return time.Time{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("sunrise-sunset status: %d", resp.StatusCode)
	}
	var out SunriseResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return time.Time{}, err
	}
	if out.Status != "OK" {
		return time.Time{}, fmt.Errorf("sunrise-sunset returned status %q", out.Status)
	}
	tm, err := time.Parse(time.RFC3339, out.Results.Sunrise)
	if err != nil {
		return time.Time{}, err
	}
	return tm, nil
}
