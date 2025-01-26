package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rifaideen/talkative"
)

type AI struct {
	url string
}

type PullResponse struct {
	Status    string `json:"status"`
	Digest    string `json:"digest"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
}

type PullCallBack func(*PullResponse, error)

func New(url string) (*AI, error) {
	return &AI{url: url}, nil
}

func (ai *AI) Pull(model string, cb PullCallBack) (<-chan bool, error) {
	if cb == nil {
		return nil, talkative.ErrCallback
	}

	client := http.Client{}

	payload := map[string]string{
		"model": model,
	}

	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, fmt.Errorf("%w:%v", talkative.ErrEncoding, err)
	}

	res, err := client.Post(
		fmt.Sprintf("%s/api/pull", ai.url),
		"application/json",
		body,
	)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			return nil, fmt.Errorf("%w\n%v", errors.New("the requested resource not found"), body)
		case http.StatusBadRequest:
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			return nil, fmt.Errorf("%w\n%v", talkative.ErrBadRequest, body)
		default:
			return nil, fmt.Errorf("%w: please make sure ollama server is running and url is correct", talkative.ErrInvoke)
		}
	}

	chDone := make(chan bool)

	go func() {
		talkative.StreamResponse(res.Body, cb)

		chDone <- true
	}()

	return chDone, nil
}
