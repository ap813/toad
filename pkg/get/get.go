package get

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// HTTPGet is the exported function for load testing
// Get calls for web services
func HTTPGet(c *cli.Context) error {

	url, headers, err := validateParameters(c)

	if err != nil {
		return err
	}

	_ = loadTest(url, headers)

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
func loadTest(url string, headers map[string]string) int {
	return 0
}
