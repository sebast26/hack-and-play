package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

const (
	baseAuthorizationURL   = "https://accounts.google.com/o/oauth2/auth"
	accessTokenExchangeURL = "https://accounts.google.com/o/oauth2/token"
	redirectURI            = "/auth2callback"
)

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
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
		log.Fatalf("%w: unable to build token request")
	}

	client := &http.Client{}
	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		log.Fatalf("%w: exchange code for token request failed")
	}
	defer func() {
		_ = tokenResponse.Body.Close()
	}()

	//var tokenResp tokenResp
	//err = json.NewDecoder(tokenResponse.Body).Decode(&tokenResp)
	//if err != nil {
	//	log.Fatalf("%w: error unmarshaling token response")
	//}

	//log.Printf("Access token: %s", tokenResp.AccessToken)

	bytes, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		return
	}

	//  00000000  7b 0a 20 20 22 65 72 72  6f 72 22 3a 20 7b 0a 20  |{.  "error": {. |
	// 00000010  20 20 20 22 63 6f 64 65  22 3a 20 34 30 30 2c 0a  |   "code": 400,.|
	// 00000020  20 20 20 20 22 6d 65 73  73 61 67 65 22 3a 20 22  |    "message": "|
	// 00000030  49 6e 76 61 6c 69 64 20  4a 53 4f 4e 20 70 61 79  |Invalid JSON pay|
	// 00000040  6c 6f 61 64 20 72 65 63  65 69 76 65 64 2e 20 55  |load received. U|
	// 00000050  6e 65 78 70 65 63 74 65  64 20 74 6f 6b 65 6e 2e  |nexpected token.|
	// 00000060  5c 6e 63 6c 69 65 6e 74  5f 69 64 3d 32 38 36 36  |\nclient_id=2866|
	// 00000070  39 31 31 36 30 39 5c 6e  5e 22 2c 0a 20 20 20 20  |911609\n^",.    |
	// 00000080  22 73 74 61 74 75 73 22  3a 20 22 49 4e 56 41 4c  |"status": "INVAL|
	// 00000090  49 44 5f 41 52 47 55 4d  45 4e 54 22 0a 20 20 7d  |ID_ARGUMENT".  }|
	// 000000a0  0a 7d 0a                                          |.}.|
	spew.Dump(bytes)
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
