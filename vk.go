package vk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const clientID = "2274003"                        //VK for Android app client_id
const clientSecret = "hHbZxrka2uZ6jB1inYsH"       //VK for Android app client_secret
const authURL = "https://oauth.vk.com/token?"     //Direct Authorization URL
const apiMethodURL = "https://api.vk.com/method/" //Method request URL

// AuthResponse structure contains all parameters of response of authorization request
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`

	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Auth function makes authorization request
func Auth(login string, password string) (*AuthResponse, error) {
	var jsonResponse *AuthResponse
	var requestURL = url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {login},
		"password":      {password},
	}
	response, err := http.Get(authURL + requestURL.Encode())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &jsonResponse); err != nil {
		return nil, err
	}
	return jsonResponse, nil
}

// Request function makes api method request
func Request(methodName string, parameters map[string]string, user *AuthResponse) ([]map[string]interface{}, error) {
	requestURL, err := url.Parse(apiMethodURL + methodName)
	if err != nil {
		return make([]map[string]interface{}, 0), err
	}
	requestQuery := requestURL.Query()
	for key, value := range parameters {
		requestQuery.Set(key, value)
	}
	requestQuery.Set("access_token", user.AccessToken)
	requestURL.RawQuery = requestQuery.Encode()
	response, err := http.Get(requestURL.String())
	if err != nil {
		return make([]map[string]interface{}, 0), err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return make([]map[string]interface{}, 0), err
	}
	return parseJSONResponse(content)
}

func parseJSONResponse(rawJSON []byte) ([]map[string]interface{}, error) {
	responseMap := make(map[string][]map[string]interface{})
	responseMapAlt := make(map[string]map[string]interface{})
	if json.Unmarshal(rawJSON, &responseMap) != nil {
		if err := json.Unmarshal(rawJSON, &responseMapAlt); err != nil {
			return make([]map[string]interface{}, 0), err
		}
		var key string
		for k := range responseMapAlt {
			key = k
			break
		}
		responseMap[key] = []map[string]interface{}{responseMapAlt[key]}
	}

	if _, value := responseMap["response"]; value {
		return responseMap["response"], nil
	} else if _, value := responseMap["error"]; value {
		return make([]map[string]interface{}, 0), errors.New(responseMap["error"][0]["error_msg"].(string))
	}
	return make([]map[string]interface{}, 0), errors.New("Response not clear")
}
