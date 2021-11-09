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
	appYml := flag.String("app", "", "Path to app.yml")
	appEnvYml := flag.String("app-env", "", "Path to app-env.yml")
	networkPolicyYml := flag.String("network-policy", "", "Path to network-policy.yml")
	appDeployYml := flag.String("app-deploy", "", "Path to app-deploy.yml")

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

	} else {

		if *appYml != "" {
			err = createApp(client, *appYml)
			if err != nil {
				log.Fatal(err)
			}
		}

		if *appEnvYml != "" {
			err = createAppEnv(client, *appEnvYml)
			if err != nil {
				log.Fatal(err)
			}
		}

		if *networkPolicyYml != "" {
			err = createNetworkPolicy(client, *networkPolicyYml)
			if err != nil {
				log.Fatal(err)
			}
		}

		if *appDeployYml != "" {
			err = deployApp(client, *appDeployYml)
			if err != nil {
				log.Fatal(err)
			}
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

	if action.App != nil {
		_, err = client.GetApp(context.TODO(), action.App.Name)
		if err == nil {
			// action exists
			return nil
		}

		err = client.CreateApp(context.TODO(), action.App)
		if err != nil {
			return fmt.Errorf("failed to create shipa action: %v", err)
		}
	}

	if action.AppEnv != nil {
		err = client.CreateAppEnvs(context.TODO(), action.AppEnv)
		if err != nil {
			return fmt.Errorf("failed to create shipa appEnv: %v", err)
		}
	}

	if action.NetworkPolicy != nil {
		err = client.CreateOrUpdateNetworkPolicy(context.TODO(), action.NetworkPolicy)
		if err != nil {
			return fmt.Errorf("failed to create shipa networkPolicy: %v", err)
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

func createApp(client *shipa.Client, path string) error {
	yamlFile, err := readFile(path)
	if err != nil {
		return err
	}

	var app shipa.App
	err = yaml.Unmarshal(yamlFile, &app)
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

func createAppEnv(client *shipa.Client, path string) error {
	yamlFile, err := readFile(path)
	if err != nil {
		return err
	}

	var appEnv shipa.CreateAppEnv
	err = yaml.Unmarshal(yamlFile, &appEnv)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	err = client.CreateAppEnvs(context.TODO(), &appEnv)
	if err != nil {
		return fmt.Errorf("failed to create shipa appEnv: %v", err)
	}

	return nil
}

func createNetworkPolicy(client *shipa.Client, path string) error {
	yamlFile, err := readFile(path)
	if err != nil {
		return err
	}

	var networkPolicy shipa.NetworkPolicy
	err = yaml.Unmarshal(yamlFile, &networkPolicy)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	err = client.CreateOrUpdateNetworkPolicy(context.TODO(), &networkPolicy)
	if err != nil {
		return fmt.Errorf("failed to create shipa networkPolicy: %v", err)
	}

	return nil
}

func deployApp(client *shipa.Client, path string) error {
	yamlFile, err := readFile(path)
	if err != nil {
		return err
	}

	var appDeploy shipa.AppDeploy
	err = yaml.Unmarshal(yamlFile, &appDeploy)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	err = client.DeployApp(context.TODO(), &appDeploy)
	if err != nil {
		return fmt.Errorf("failed to deploy shipa app: %v", err)
	}

	return nil
}
