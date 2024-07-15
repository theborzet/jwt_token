package externalapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExternalApiClient struct {
	ApiUrl string
}

func NewExternalApiClient(apiUrl string) *ExternalApiClient {
	return &ExternalApiClient{ApiUrl: apiUrl}
}

type response struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func (e *ExternalApiClient) FetchPeopleInfo(passportSerie, passportNumber string) (*response, error) {
	url := fmt.Sprintf("%spassportSerie=%s&passportNumber=%s", e.ApiUrl, passportSerie, passportNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch people info: status code %d", resp.StatusCode)
	}

	var response response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
