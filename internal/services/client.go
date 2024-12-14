package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devmegablaster/bashform/internal/models"
)

type Client struct {
	BaseUrl    string
	PubKey     string
	httpClient *http.Client
}

func NewClient(baseUrl string, pubKey string) *Client {
	return &Client{
		BaseUrl:    baseUrl,
		PubKey:     pubKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetForm(code string) models.Form {
	req, err := http.NewRequest("GET", c.BaseUrl+"/forms/"+code, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("PubKey", c.PubKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return models.Form{}
	}

	defer resp.Body.Close()

	var formResponse models.FormResponse
	if err := json.NewDecoder(resp.Body).Decode(&formResponse); err != nil {
		fmt.Println(err)
		return models.Form{}
	}

	return formResponse.Data
}

func (c *Client) SubmitForm(code string, response models.Response) {
	reqBody, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", c.BaseUrl+"/forms/"+code, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("PubKey", c.PubKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
}
