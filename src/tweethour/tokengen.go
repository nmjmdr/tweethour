package tweethour

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const consumerKey = "CONSUMER_KEY"
const consumerSecret = "CONSUMER_SECRET"


const tokenUrl = "https://api.twitter.com/oauth2/token"
const authHeader = "Authorization"
const authHeaderValue = "Basic"
const contentTypeHeader = "Content-Type"
const contentTypeValue = "application/x-www-form-urlencoded;charset=UTF-8"
const grantTypeBody = "grant_type=client_credentials"

type TokenGen interface {
	Generate() (string, Error)
}

type tokenGen struct {
	client *http.Client
}

func NewTokenGen(client *http.Client) TokenGen {
	t := new(tokenGen)
	t.client = client
	return t
}

func (t *tokenGen) Generate() (string, Error) {

	response, err := t.makeRequest()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", NewTokenError(errors.New("Failed to fetch token"))
	}

	var v map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&v)

	if err != nil {
		return "", NewTokenError(err)
	}

	tokenType, err := getValue("token_type", v)
	if err != nil {
		return "", NewTokenError(err)
	}

	if tokenType != "bearer" {
		return "", NewTokenError(errors.New("Unable to get bearer token, an unrecognized token was returned: " + tokenType))
	}

	var token string
	token, err = getValue("access_token", v)

	if err != nil {
		return "", NewTokenError(err)
	}

	return token, nil

}

func (t *tokenGen) makeRequest() (*http.Response, error) {
	buf := bytes.NewBuffer([]byte(grantTypeBody))

	req, err := http.NewRequest("POST", tokenUrl, buf)

	if err != nil {
		return nil, err
	}

	encodedKey := getEncodedKey(consumerKey,consumerSecret)

	req.Header.Add(authHeader, fmt.Sprintf("%s %s", authHeaderValue, encodedKey))
	req.Header.Add(contentTypeHeader, contentTypeValue)

	var response *http.Response

	response, err = t.client.Do(req)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Error response : " + http.StatusText(response.StatusCode))
	}

	return response, nil

}
