package client

import (
	"encoding/json"
	"fmt"
	"goclass/internal/domain"
	"io"
	"log"
	"net/http"
)

const (
	url = "https://jsonplaceholder.typicode.com/posts"
)

type Client struct {
}

func New() Client {
	return Client{}
}

func (c Client) Get() ([]domain.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(body io.Closer) {
		err := body.Close()
		if err != nil {
			log.Println("error closing the response body")
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("the request was unsuccessful: %d (%s)", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	posts := make([]domain.Post, 0)
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
