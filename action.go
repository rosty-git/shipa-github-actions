package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/brunoa19/shipa-github-actions/shipa"
	"github.com/brunoa19/shipa-github-actions/types"
	"gopkg.in/yaml.v2"
)

func main() {
	shipaActionYml := flag.String("shipa-action", "", "Path to shipa-action.yml")
	debug := flag.Bool("debug", false, "Enables debug mode")
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
	client.SetDebugMode(*debug)

	if *shipaActionYml != "" {
		err := createShipaAction(client, *shipaActionYml)
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
	App           *shipa.CreateAppRequest `yaml:"app,omitempty"`
	AppEnv        *shipa.CreateAppEnv     `yaml:"app-env,omitempty"`
	AppCname      *shipa.AppCname         `yaml:"app-cname,omitempty"`
	NetworkPolicy *shipa.NetworkPolicy    `yaml:"network-policy,omitempty"`
	AppDeploy     *shipa.AppDeploy        `yaml:"app-deploy,omitempty"`
	Framework     *shipa.PoolConfig       `yaml:"framework,omitempty"`
	Cluster       *types.Cluster          `yaml:"cluster,omitempty"`
	Job           *shipa.JobCreateRequest `yaml:"job,omitempty"`
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

	if action.Cluster != nil {
		err = createClusterIfNotExist(client, action.Cluster)
		if err != nil {
			return err
		}
	}

	if action.App != nil {
		err = createAppIfNotExist(client, action.App)
		if err != nil {
			return err
		}
	}

	if action.AppEnv != nil {
		err = client.CreateAppEnvs(context.TODO(), action.AppEnv)
		if err != nil {
			return fmt.Errorf("failed to create shipa app-env: %v", err)
		}
	}

	if action.AppCname != nil {
		err = client.CreateAppCname(context.TODO(), action.AppCname)
		if err != nil {
			return fmt.Errorf("failed to create shipa app-cname: %v", err)
		}
	}

	if action.NetworkPolicy != nil {
		err = client.CreateOrUpdateNetworkPolicy(context.TODO(), action.NetworkPolicy)
		if err != nil {
			return fmt.Errorf("failed to create shipa network-policy: %v", err)
		}
	}

	if action.AppDeploy != nil {
		action.AppDeploy.SetDefaults()
		err = client.DeployApp(context.TODO(), action.AppDeploy)
		if err != nil {
			return fmt.Errorf("failed to deploy shipa app: %v", err)
		}
	}

	if action.Job != nil {
		err = createJobIfNotExist(client, action.Job)
		if err != nil {
			return err
		}
	}

	return nil
}

func createJobIfNotExist(client *shipa.Client, job *shipa.JobCreateRequest) error {
	jobs, err := client.ListJobs(context.TODO())
	if err != nil {
		return err
	}

	for _, j := range jobs {
		if j.Name == job.Name {
			return nil
		}
	}

	_, err = client.CreateJob(context.TODO(), job)
	return err
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

func createAppIfNotExist(client *shipa.Client, app *shipa.CreateAppRequest) error {
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

func createClusterIfNotExist(client *shipa.Client, input *types.Cluster) error {
	cluster, err := input.ToShipaCluster()
	if err != nil {
		return fmt.Errorf("failed to parse shipa cluster: %v", err)
	}

	shipaCluster, err := client.GetCluster(context.TODO(), cluster.Name)
	if err != nil && strings.Contains(err.Error(), "cluster not found") {
		// cluster does not exist
		err = client.CreateCluster(context.TODO(), cluster)
		if err != nil {
			return fmt.Errorf("failed to create shipa cluster: %v", err)
		}
		return nil
	}

	// check if need to add new frameworks to the cluster
	newFrameworks := getNewFrameworks(shipaCluster, cluster)
	if len(newFrameworks) > 0 {
		if shipaCluster.Resources == nil {
			shipaCluster.Resources = &shipa.ClusterResources{}
		}

		for _, name := range newFrameworks {
			shipaCluster.Resources.Frameworks = append(shipaCluster.Resources.Frameworks, &shipa.Framework{
				Name: name,
			})
		}

		err = client.UpdateCluster(context.TODO(), shipaCluster)
		if err != nil {
			return fmt.Errorf("failed to update shipa cluster: %v", err)
		}
	}

	return nil
}

func getNewFrameworks(current *shipa.Cluster, newCluster *shipa.Cluster) []string {
	currentFrameworks := convertFrameworksToMap(current)
	newFrameworks := convertFrameworksToMap(newCluster)

	var result []string
	for name, _ := range newFrameworks {
		if !currentFrameworks[name] {
			result = append(result, name)
		}
	}

	return result
}

func convertFrameworksToMap(cluster *shipa.Cluster) map[string]bool {
	if cluster.Resources == nil || cluster.Resources.Frameworks == nil {
		return nil
	}

	result := make(map[string]bool)
	for _, framework := range cluster.Resources.Frameworks {
		if framework != nil {
			result[framework.Name] = true
		}
	}
	return result
}
