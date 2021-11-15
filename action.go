package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/brunoa19/shipa-github-actions/shipa"
	"gopkg.in/yaml.v2"
)

func main() {
	shipaActionYml := flag.String("shipa-action", "", "Path to shipa-action.yml")
	flag.Parse()

	if _, ok := os.LookupEnv("SHIPA_HOST"); !ok {
		log.Fatal("SHIPA_HOST env not set")
	}

	if _, ok := os.LookupEnv("SHIPA_TOKEN"); !ok {
		log.Fatal("SHIPA_TOKEN env not set")
	}

	client, err := shipa.New()
	if err != nil {
		log.Fatal("failed to create shipa client:", err)
	}

	if *shipaActionYml != "" {
		err = createShipaAction(client, *shipaActionYml)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func readFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("invalid file path: %v", err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return bytes, nil
}

type ShipaAction struct {
	App           *shipa.App           `yaml:"app,omitempty"`
	AppEnv        *shipa.CreateAppEnv  `yaml:"app-env,omitempty"`
	NetworkPolicy *shipa.NetworkPolicy `yaml:"network-policy,omitempty"`
	AppDeploy     *shipa.AppDeploy     `yaml:"app-deploy,omitempty"`
	Framework     *shipa.PoolConfig    `yaml:"framework,omitempty"`
}

func createShipaAction(client *shipa.Client, path string) error {
	yamlFile, err := readFile(path)
	if err != nil {
		return err
	}

	var action ShipaAction
	err = yaml.Unmarshal(yamlFile, &action)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	if action.Framework != nil {
		err = createFrameworkIfNotExist(client, action.Framework)
		if err != nil {
			return err
		}
	}

	if action.App != nil {
		err = createFrameworkIfNotExist(client, &shipa.PoolConfig{Name: action.App.Pool})
		if err != nil {
			return err
		}

		err = createAppIfNotExist(client, action.App)
	}

	if action.AppEnv != nil {
		err = client.CreateAppEnvs(context.TODO(), action.AppEnv)
		if err != nil {
			return fmt.Errorf("failed to create shipa app-env: %v", err)
		}
	}

	if action.NetworkPolicy != nil {
		err = client.CreateOrUpdateNetworkPolicy(context.TODO(), action.NetworkPolicy)
		if err != nil {
			return fmt.Errorf("failed to create shipa network-policy: %v", err)
		}
	}

	if action.AppDeploy != nil {
		err = client.DeployApp(context.TODO(), action.AppDeploy)
		if err != nil {
			return fmt.Errorf("failed to deploy shipa app: %v", err)
		}
	}

	return nil
}

func createFrameworkIfNotExist(client *shipa.Client, framework *shipa.PoolConfig) error {
	_, err := client.GetPoolConfig(context.TODO(), framework.Name)
	if err != nil {
		// framework does not exist
		err = client.CreatePoolConfig(context.TODO(), framework)
		if err != nil {
			return fmt.Errorf("failed to create shipa framework: %v", err)
		}
	}
	return nil
}

func createAppIfNotExist(client *shipa.Client, app *shipa.App) error {
	_, err := client.GetApp(context.TODO(), app.Name)
	if err != nil {
		// app does not exist
		err = client.CreateApp(context.TODO(), app)
		if err != nil {
			return fmt.Errorf("failed to create shipa app: %v", err)
		}
	}
	return nil
}
