package vus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client *http.Client
var debug bool

// VirtualUserBody is called by all http types to represent a virtual user for
// http requests with a json body
func VirtualUserBody(c *http.Client, method, url string, body []byte, headers map[string]string, delay int, duration int, debugging bool) {

	// Make sure that request can be made properly
	_, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Request parameters are not valid. Error at:", err.Error())
		return
	}

	// Set client and debug preference
	client = c
	debug = debugging

	// Create channel to stop virtual user
	stopChannel := make(chan bool)
	stopUser := make(chan bool)

	// Run
	go func() {
		for {
			select {
			case <-stopChannel:
				break
			default:
				// Make http call then delay specified by user
				call(method, url, body, headers)
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}()

	// Stop virtual user after duration
	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		stopChannel <- true // Stops http calls
		stopUser <- true    // Allows continuation of function execution
	})

	<-stopUser // Wait for timer
}

// Function that is used to call service with request
func call(method, url string, body []byte, headers map[string]string) {

	// Create a local request so that the original is reusable
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	// Don't sending bad requests
	if err != nil {
		fmt.Println("Problem with creating request inside call method, please report issue")
		return
	}

	// Apply headers
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// Option to print response data from http calls
	if debug {
		resp, err := client.Do(req)

		// Present relevant information to the user
		if err == nil {
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Service responded with bad json and status code ", resp.StatusCode)
				return
			}
			fmt.Println("Status Code: ", resp.StatusCode, " Response Body: ", string(bodyBytes))
		} else {
			fmt.Println("Request failed to send to service with:", err.Error())
		}
	} else {
		client.Do(req)
	}
}

// VirtualUser is called by all http types to represent a virtual user
func VirtualUser(c *http.Client, method, url string, headers map[string]string, delay int, duration int, debugging bool) {

	// Make sure that request can be made properly
	_, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Request parameters are not valid. Error at:", err.Error())
		return
	}

	// Set client and debug preference
	client = c
	debug = debugging

	// Create channel to stop virtual user
	stopChannel := make(chan bool)
	stopUser := make(chan bool)

	// Run
	go func() {
		for {
			select {
			case <-stopChannel:
				break
			default:
				// Make http call then delay specified by user
				callNoBody(method, url, headers)
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}()

	// Stop virtual user after duration
	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		stopChannel <- true // Stops http calls
		stopUser <- true    // Allows continuation of function execution
	})

	<-stopUser // Wait for timer
}

// Function that is used to call service with request
func callNoBody(method, url string, headers map[string]string) {

	// Create a local request so that the original is reusable
	req, err := http.NewRequest(method, url, nil)

	// Don't sending bad requests
	if err != nil {
		fmt.Println("Problem with creating request inside call method, please report issue")
		return
	}

	// Apply headers
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// Option to print response data from http calls
	if debug {
		resp, err := client.Do(req)

		// Present relevant information to the user
		if err == nil {
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Service responded with bad json and status code ", resp.StatusCode)
				return
			}
			fmt.Println("Status Code: ", resp.StatusCode, " Response Body: ", string(bodyBytes))
		} else {
			fmt.Println("Request failed to send to service with:", err.Error())
		}
	} else {
		client.Do(req)
	}
}
