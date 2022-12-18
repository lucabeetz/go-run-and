package main

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
	"time"

	"github.com/theckman/yacspin"
	"github.com/tidwall/gjson"
)

func makeRequest(prompt string, apiKey string) string {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          "",
		SuffixAutoColon: true,
		Message:         "",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}
	spinner, err := yacspin.New(cfg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	spinner.Start()
	spinner.Message("Making API request")

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

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	// Read response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	spinner.Stop()
	return string(resBody)
}

func runCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %s", err)
		os.Exit(1)
	}
}

func main() {
	// Check for correct usage
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

	// Construct prompt and make completion request
	textRequest := strings.TrimSuffix(os.Args[1], "\n")
	prompt := fmt.Sprintf("%s (on macos)\n```bash\n#!/bin/bash\n", textRequest)
	resBody := makeRequest(prompt, apiKey)

	// Get completion and trim newlines
	completion := gjson.Get(string(resBody), "choices.0.text").String()
	completion = strings.TrimPrefix(completion, "\n")
	completion = strings.TrimSuffix(completion, "\n")

	// Ask user for confirmation to run suggested command
	fmt.Printf("Suggested:\n%s\n", completion)
	fmt.Printf("Run? [y/N] | [e] for explanation\n")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	// Get input
	input = strings.ToLower(strings.TrimSuffix(input, "\n"))

	// Run suggested command if user confirms
	if input == "y" {
		runCommand(completion)
		return
	} else if input == "e" {
		explanationPrompt := fmt.Sprintf("%s (on macos)\n```bash\n#!/bin/bash\n%s\nExplanation\n", textRequest, completion)
		resBody = makeRequest(explanationPrompt, apiKey)

		explanation := gjson.Get(resBody, "choices.0.text").String()
		explanation = strings.TrimPrefix(explanation, "\n")
		explanation = strings.TrimSuffix(explanation, "\n")

		fmt.Printf("Explanation:\n%s\n", explanation)
	} else {
		fmt.Println("Aborting")
		os.Exit(1)
	}

	fmt.Printf("Run? [y/N]\n")

	reader = bufio.NewReader(os.Stdin)
	input, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	// Get input
	input = strings.ToLower(strings.TrimSuffix(input, "\n"))

	if input == "y" {
		runCommand(completion)
	}
}
