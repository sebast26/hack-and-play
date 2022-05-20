package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	baseAuthorizationURL   = "https://accounts.google.com/o/oauth2/auth"
	accessTokenExchangeURL = "https://oauth2.googleapis.com/token"
	redirectURI            = "/auth2callback"
	taskAPIURL             = "https://www.googleapis.com/tasks/v1/lists/@default/tasks"
)

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func step1ObtainAuthorizationLinkLink(writer http.ResponseWriter, request *http.Request) {
	authorizationLinkParams := buildAuthorizationLinkParams()
	authorizationLink := fmt.Sprintf("%s?%s", baseAuthorizationURL, authorizationLinkParams)
	log.Printf("Auth URL: %s", authorizationLink)

	fmt.Fprintf(writer, "<a href='%s'>%s</a>", authorizationLink, authorizationLink)
}

func buildAuthorizationLinkParams() string {
	clientID := os.Getenv("CLIENT_ID")
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", "http://localhost:8090"+redirectURI)
	params.Add("scope", "https://www.googleapis.com/auth/tasks")
	params.Add("response_type", "code")
	params.Add("approval_prompt", "force")
	params.Add("access_type", "offline")
	return params.Encode()
}

func step2ExchangeCodeToToken(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query().Get("code")
	log.Printf("Auth code: %s", code)

	params := buildExchangeCodeParams(code)
	tokenRequest, err := http.NewRequest(http.MethodPost, accessTokenExchangeURL, strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatalf("%v: unable to build token request", err)
	}

	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		log.Fatalf("%v: exchange code for token request failed", err)
	}
	defer func() {
		_ = tokenResponse.Body.Close()
	}()

	var tokenResp tokenResp
	err = json.NewDecoder(tokenResponse.Body).Decode(&tokenResp)
	if err != nil {
		log.Fatalf("%v: error unmarshaling token response", err)
	}

	log.Printf("Access token: %s", tokenResp.AccessToken)

	tasks := callTaskAPI(tokenResp.AccessToken)

	fmt.Fprintf(writer, tasks)
}

func callTaskAPI(token string) string {
	request, err := http.NewRequest(http.MethodGet, taskAPIURL, nil)
	if err != nil {
		log.Fatalf("%v: cannot create Task API request")
	}

	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("%v: error when calling Task API")
	}
	defer func() {
		_ = response.Body.Close()
	}()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("%v: error reading response from Task API")
	}

	return string(b)
}

func buildExchangeCodeParams(code string) url.Values {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	params := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://localhost:8090" + redirectURI},
	}
	return params
}

func main() {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("app requires CLIENT_ID and CLIENT_SECRET")
	}

	http.HandleFunc(redirectURI, step2ExchangeCodeToToken)
	http.HandleFunc("/", step1ObtainAuthorizationLinkLink)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("%w: could not start http server")
	}
}
