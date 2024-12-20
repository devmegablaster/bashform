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

func (c *Client) GetForm(code string) (models.Form, error) {
	req, err := http.NewRequest("GET", c.BaseUrl+"/forms/"+code, nil)
	if err != nil {
		return models.Form{}, err
	}

	req.Header.Add("PubKey", c.PubKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return models.Form{}, err
	}

	defer resp.Body.Close()

	var formResponse models.FormResponse
	if err := json.NewDecoder(resp.Body).Decode(&formResponse); err != nil {
		fmt.Println(err)
		return models.Form{}, err
	}

	return formResponse.Data, nil
}

func (c *Client) GetForms() ([]models.Form, error) {
	req, err := http.NewRequest("GET", c.BaseUrl+"/forms", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("PubKey", c.PubKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	var formResponse models.FormsResponse
	if err := json.NewDecoder(resp.Body).Decode(&formResponse); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return formResponse.Data, nil
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

func (c *Client) CreateForm(form models.FormRequest) (*models.FormResponse, error) {
	reqBody, err := json.Marshal(form)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseUrl+"/forms", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("PubKey", c.PubKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var formResponse models.FormResponse
	if err := json.NewDecoder(resp.Body).Decode(&formResponse); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &formResponse, nil
}

func (c *Client) GetResponses(formID string) (*models.Form, error) {
	req, err := http.NewRequest("GET", c.BaseUrl+"/forms/"+formID+"/responses", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("PubKey", c.PubKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	var responseResponse models.FormResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseResponse); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &responseResponse.Data, nil
}
