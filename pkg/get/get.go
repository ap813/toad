package get

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli"
)

var client *http.Client

// HTTPGet is the exported function for load testing
// Get calls for web services
func HTTPGet(c *cli.Context) error {

	url, headers, err := validateParameters(c)

	if err != nil {
		return err
	}

	timeout := c.Int("timeout")

	vus := c.Int("vus")

	delay := c.Int("delay")

	duration := c.Int("duration")

	_, _ = loadTest(url, headers, timeout, vus, delay, duration)

	return nil
}

// Valid params from user call to cli
func validateParameters(c *cli.Context) (string, map[string]string, error) {

	url := c.String("url")

	if len(url) <= 2 {
		return "", nil, errors.New("'url' too short to be valid")
	}

	// Headers stored in map for calls later
	headersMap := make(map[string]string)

	headers := c.String("headers")

	if len(headers) != 0 {

		fmt.Println(headers)

		headers = strings.Trim(headers, " ")

		headers = strings.Trim(headers, ",")

		headersArray := strings.Split(headers, ",")

		for _, val := range headersArray {

			val = strings.Trim(val, " ")

			valArray := strings.Split(val, ":")

			if len(valArray) != 2 {
				fmt.Println(valArray, len(valArray))
				return "", nil, errors.New("'headers' had key that didn't correlate to a value")
			}

			headersMap[valArray[0]] = valArray[1]
		}

		fmt.Println("Headers:", headersMap)
	}

	return url, headersMap, nil
}

// Perform the load testing on the service
func loadTest(url string, headers map[string]string, timeout, vus, delay, duration int) (float64, float64, error) {

	// Create http client used for calls
	client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Get Request
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return 0, 0, errors.New("Invalid URL passed")
	}

	// Apply headers
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	var wg sync.WaitGroup

	wg.Add(vus)

	for i := 0; i < vus; i++ {
		go func() {
			virtualUser(request, delay, duration)
			wg.Done()
		}()
	}
}

func virtualUser(request *http.Request, delay int, duration int) {

	// Create channel to stop virtual user
	stopChannel := make(chan bool)

	// Run
	go func() {
		for {
			select {
			case <-stopChannel:
				break
			default:
				getCall(*request)
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}()

	// Stop virtual user after duration
	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		stopChannel <- true
	})
}

// Function that is used to call service with request
func getCall(request http.Request) {
	req := request
	resp, err := client.Do(&req)
	if err != nil {
		fmt.Println("Request failed")
	} else {
		fmt.Println("Status Code: ", resp.StatusCode)
	}
}
