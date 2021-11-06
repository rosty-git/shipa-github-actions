package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/brunoa19/shipa-github-actions/shipa"
	"gopkg.in/yaml.v2"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("no input arg")
	}

	if len(args) != 2 {
		log.Fatal("expected 2 args")
	}

	if _, ok := os.LookupEnv("SHIPA_HOST"); !ok {
		log.Fatal("SHIPA_HOST env not set")
	}

	if _, ok := os.LookupEnv("SHIPA_TOKEN"); !ok {
		log.Fatal("SHIPA_TOKEN env not set")
	}

	objectType, path := args[0], args[1]
	if _, err := os.Stat(path); err != nil {
		log.Fatal("invalid file path:", err)
	}

	client, err := shipa.New()
	if err != nil {
		log.Fatal("failed to create shipa client:", err)
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	switch objectType {
	case "app":
		err = createApp(client, yamlFile)
	case "appDeploy":
		err = deployApp(client, yamlFile)
	default:
		err = fmt.Errorf("object type not recognized: %v", objectType)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func createApp(client *shipa.Client, yamlFile []byte) error {
	var app shipa.App
	err := yaml.Unmarshal(yamlFile, &app)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	_, err = client.GetApp(context.TODO(), app.Name)
	if err == nil {
		// app exists
		return nil
	}

	err = client.CreateApp(context.TODO(), &app)
	if err != nil {
		return fmt.Errorf("failed to create shipa app: %v", err)
	}

	return nil
}

func deployApp(client *shipa.Client, yamlFile []byte) error {
	var appDeploy shipa.AppDeploy
	err := yaml.Unmarshal(yamlFile, &appDeploy)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	err = client.DeployApp(context.TODO(), &appDeploy)
	if err != nil {
		return fmt.Errorf("failed to deploy shipa app: %v", err)
	}

	return nil
}
