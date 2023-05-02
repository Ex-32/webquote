package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var help = flag.Bool("help", false, "Show help")
var api_key string

const api_url string = "https://api.api-ninjas.com/v1/quotes"

func exit_on_err(err *error) {
	if err == nil {
		return
	}
	if *err != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", (*err).Error()))
		os.Exit(1)
	}
}

func main() {
	flag.StringVar(
		&api_key,
		"api-key",
		"",
		"Custom API key for api-ninjas.com",
	)

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// manually inserting default API key value prevents it from showing
	// default value
	if len(api_key) == 0 {
		api_key = "ffbKimU4CDIf9HXmAPQvNQ==5N3ORyXLoJ4QOHzF"
	}

	// setup API GET request and add header value with API key
	req, err := http.NewRequest("GET", api_url, nil)
	exit_on_err(&err)
	req.Header.Set("X-Api-Key", api_key)

	// create http client and perform http request to get content
	client := &http.Client{}
	resp, err := client.Do(req)
	exit_on_err(&err)
	defer resp.Body.Close()

	// read content of http request's body into []byte
	body, err := io.ReadAll(resp.Body)
	exit_on_err(&err)

	if resp.StatusCode == 200 {
		// parse json from request body
		var quotes []map[string]interface{}
		err = json.Unmarshal(body, &quotes)
		exit_on_err(&err)

		// print content of request body
		quote := quotes[0]
		fmt.Printf(
			"%s\n~ %s\n",
			quote["quote"],
			quote["author"],
		)

		// return success
		return
	} else if resp.StatusCode == 400 {
		var error_json map[string]interface{}
		err = json.Unmarshal(body, &error_json)
		exit_on_err(&err)
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", error_json["error"]))
	} else if resp.StatusCode/100 == 4 {
		os.Stderr.WriteString(fmt.Sprintf("error: client: %s", resp.Status))
	} else if resp.StatusCode/100 == 5 {
		os.Stderr.WriteString(fmt.Sprintf("error: server: %s\n", resp.Status))
	} else {
		os.Stderr.WriteString(
			fmt.Sprintf("error: unexpected response: %s\n", resp.Status),
		)
	}

	// unless we successfully returned, exit with failure
	os.Exit(1)
}
