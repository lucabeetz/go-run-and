package gra

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

	// Construct request body
	textRequest := strings.TrimSuffix(os.Args[1], "\n")
	prompt := fmt.Sprintf("%s (on macos)\n```bash\n#!/bin/bash\n", textRequest)
	requestBody := map[string]interface{}{
		"model":       "code-davinci-002",
		"prompt":      prompt,
		"temperature": 0,
		"max_tokens":  256,
		"stop":        "```",
	}

	// Construct completion request body
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
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

	// Get completion and trim newlines
	completion := gjson.Get(string(resBody), "choices.0.text").String()
	completion = strings.TrimPrefix(completion, "\n")
	completion = strings.TrimSuffix(completion, "\n")

	// Ask user for confirmation to run suggested command
	fmt.Printf("Suggested:\n%s\n", completion)
	fmt.Printf("Run? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	// Run suggested command if user confirms
	if strings.ToLower(strings.TrimSuffix(input, "\n")) == "y" {
		cmd := exec.Command("bash", "-c", completion)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error running command: %s", err)
			os.Exit(1)
		}
	}
}
