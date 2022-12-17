package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gra <request>")
		os.Exit(1)
	}

	// Check if OPENAI_API_KEY is set
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		fmt.Println("OPENAI_API_KEY is not set")
		os.Exit(1)
	}

	textRequest := os.Args[1]
	requestBody := fmt.Sprintf(`{"model": "code-davinci-002", "prompt": "%s\n'''bash\n", "temperature": 0, "max_tokens": 100, "stop": "'''"}`, textRequest)

	// Construct completion request body
	jsonBody := []byte(requestBody)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bodyReader)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	completion := gjson.Get(string(resBody), "choices.0.text").String()
	// completion = strings.TrimRight(completion, "\n")

	fmt.Printf("Suggested:\n%s", completion)
}
