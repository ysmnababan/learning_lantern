package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func GetExchangeRate(from_rate, to_rate string) (float64, error) {
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/3d954472bc85adbafad0a4f9/latest/%s", from_rate)

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	// req.Header.Set("Content-Type", "application/json")
	// Encode API key to Base64
	// auth := fmt.Sprintf("%s:", apiKey)
	// b64 := base64.StdEncoding.EncodeToString([]byte(auth))
	// req.Header.Set("Authorization", "Basic "+b64)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// For debugging, log the request and response
	reqDump, _ := httputil.DumpRequest(req, true)
	log.Printf("REQUEST:\n%s", string(reqDump))
	log.Printf("RESPONSE:\n%s", body)

	// unmarshall body
	type ExchangeRateResponse struct {
		Result             string             `json:"result"`
		Documentation      string             `json:"documentation"`
		TermsOfUse         string             `json:"terms_of_use"`
		TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
		TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
		TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
		TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
		BaseCode           string             `json:"base_code"`
		ConversionRates    map[string]float64 `json:"conversion_rates"`
	}
	pr := ExchangeRateResponse{}
	if err := json.Unmarshal(body, &pr); err != nil {
		return 0, err
	}

	return pr.ConversionRates[to_rate], nil
}
