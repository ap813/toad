package get

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"toad/pkg/vus"

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

	debug := c.Bool("debug")

	err = loadTest(url, headers, timeout, vus, delay, duration, debug)

	return err
}

// Valid params from user call to cli
func validateParameters(c *cli.Context) (string, map[string]string, error) {

	url := c.String("url")

	if len(url) <= 2 {
		return "", nil, errors.New("'url' too short to be valid")
	}

	fmt.Println("Calling URL: ", url)

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
func loadTest(url string, headers map[string]string, timeout, users, delay, duration int, debug bool) error {

	// Create http client used for calls
	client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	var wg sync.WaitGroup

	// 1 virtual user = 1 go routine
	wg.Add(users)

	// Spawn all the virtual users, close worker once finished
	for i := 0; i < users; i++ {
		go func() {
			vus.VirtualUser(client, "GET", url, headers, delay, duration, debug)
			wg.Done()
		}()
	}

	// Wait for all virtual users to finish execution
	wg.Wait()

	fmt.Println("Load test complete")

	return nil
}
