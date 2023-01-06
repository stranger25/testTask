package httphandler

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testTask/internal/model"
)

func (c *HTTPClient) GetData(date string) (*model.Rates, error) {
	req, err := http.NewRequest("GET", c.Url+date, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var rates model.Rates

	err = xml.Unmarshal(body, &rates)
	if err != nil {
		return nil, err
	}
	return &rates, nil
}
