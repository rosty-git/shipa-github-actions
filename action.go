package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/brunoa19/shipa-github-actions/shipa"
	"gopkg.in/yaml.v2"
)

func main() {
	if _, ok := os.LookupEnv("SHIPA_HOST"); !ok {
		log.Fatal("SHIPA_HOST env not set")
	}

	if _, ok := os.LookupEnv("SHIPA_TOKEN"); !ok {
		log.Fatal("SHIPA_TOKEN env not set")
	}

	path, err := getFilePath()
	if err != nil {
		log.Fatal("invalid file path:", err)
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	var app shipa.App
	err = yaml.Unmarshal(yamlFile, &app)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	client, err := shipa.New()
	if err != nil {
		log.Fatal("failed to create shipa client:", err)
	}

	err = client.CreateApp(context.TODO(), &app)
	log.Fatal("failed to create shipa app:", err)
}

func getFilePath() (string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "", errors.New("no input arg")
	}

	path := args[0]

	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}
