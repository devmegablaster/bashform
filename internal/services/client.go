package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/devmegablaster/bashform/internal/constants"
	"github.com/devmegablaster/bashform/internal/models"
)

type Client struct {
	baseURL    string
	pubKey     string
	httpClient *http.Client
}

func NewClient(baseURL string, pubKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		pubKey:     pubKey,
		httpClient: &http.Client{},
	}
}

// Get form by code
func (c *Client) GetForm(code string) (models.Form, error) {
	var formResponse models.FormResponse
	if err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s", constants.API_FORMS_PATH, code), nil, &formResponse, true); err != nil {
		log.Error(err)
		return models.Form{}, fmt.Errorf(err.Data.Error)
	}

	return formResponse.Data, nil
}

// CodeResponse is the response from the checkCode endpoint
type CodeResponse struct {
	CodeData struct {
		Avaliable bool `json:"available"`
	} `json:"data"`
}

// Check if code is available
func (c *Client) CheckCode(code string) (bool, error) {
	var codeResponse CodeResponse
	if err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s/available", constants.API_FORMS_PATH, code), nil, &codeResponse, true); err != nil {
		log.Error(err)
		return false, fmt.Errorf(err.Data.Error)
	}

	return codeResponse.CodeData.Avaliable, nil
}

// Get all the forms created by the user
func (c *Client) GetForms() ([]models.Form, error) {
	var formResponse models.FormsResponse
	if err := c.doRequest(http.MethodGet, constants.API_FORMS_PATH, nil, &formResponse, true); err != nil {
		log.Error(err)
		return nil, fmt.Errorf(err.Data.Error)
	}

	return formResponse.Data, nil
}

// Submit a response to a form using formID and response
func (c *Client) SubmitForm(id string, response models.Response) error {
	if err := c.doRequest(http.MethodPost, fmt.Sprintf("%s/%s", constants.API_FORMS_PATH, id), response, nil, true); err != nil {
		log.Error(err)
		return fmt.Errorf(err.Data.Error)
	}

	return nil
}

// Create a form
func (c *Client) CreateForm(form models.FormRequest) (*models.FormResponse, error) {
	var formResponse models.FormResponse
	if err := c.doRequest(http.MethodPost, constants.API_FORMS_PATH, form, &formResponse, true); err != nil {
		log.Error(err)
		return nil, fmt.Errorf(err.Data.Error)
	}

	return &formResponse, nil
}

// Get form with responses for a form using formID
func (c *Client) GetResponses(formID string) (*models.Form, error) {
	var responseResponse models.FormResponse
	err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s/responses", constants.API_FORMS_PATH, formID), nil, &responseResponse, true)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf(err.Data.Error)
	}

	return &responseResponse.Data, nil
}

// doRequest performs the HTTP request
func (s *Client) doRequest(method, path string, data interface{}, response interface{}, auth bool) *models.ApiError {
	url := fmt.Sprintf("%s%s", s.baseURL, path)

	var buf bytes.Buffer
	if data != nil {
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return models.ErrorToApiError(err)
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return models.ErrorToApiError(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if auth {
		req.Header.Set(constants.API_AUTH_HEADER, s.pubKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return models.ErrorToApiError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr models.ApiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return models.ErrorToApiError(err)
		}

		return &apiErr
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return models.ErrorToApiError(err)
		}
	}

	return nil
}
