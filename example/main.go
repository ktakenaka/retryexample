package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	"github.com/ktakenaka/retryexample/retry"
)

const (
	demoEndpoint = "http://google.com/api/demo"
)

var (
	retryableErr = errors.New("timeout")
)

func isRetryable(err error) bool {
	if errors.Is(err, retryableErr) {
		return true
	}
	return false
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET", demoEndpoint,
		func(req *http.Request) (*http.Response, error) { return nil, retryableErr },
	)

	var (
		res *http.Response
		err error

		t = time.Now()
		i = 1
	)
	err = retry.Do(
		ctx,
		func() error {
			fmt.Println(i, time.Now().Sub(t))
			res, err = request()
			i++
			return err
		},
		retry.CheckRetryable(isRetryable),
	)

	fmt.Printf("res: %v\n", res)
	fmt.Printf("err: %v\n", err)
}

func request() (*http.Response, error) {
	return http.Get(demoEndpoint)
}
