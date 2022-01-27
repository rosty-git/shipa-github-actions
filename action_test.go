package main

import (
	"testing"

	"github.com/brunoa19/shipa-github-actions/shipa"
	"github.com/brunoa19/shipa-github-actions/types"
	"github.com/stretchr/testify/assert"
)

// cluster:
//  name: gke-actions
//  endpoint:
//    addresses: ["https://35.237.203.24"]
//    caCert: ./infra/cert.crt
//    token: ./infra/token
//  resources:
//    frameworks:
//      name: ["dev-policy", "gha-prod"]

func Test_createClusterIfNotExist(t *testing.T) {
	client, _ := shipa.New()
	client.SetDebugMode(true)

	cluster := &types.Cluster{
		Name: "gke-actions",
		Endpoint: &types.ClusterEndpoint{
			Addresses:   []string{"https://35.237.203.24"},
			Certificate: "./infra/cert.crt",
			Token:       "./infra/token",
		},
		Resources: &types.ClusterResources{
			Frameworks: &types.Framework{
				Name: []string{"dev-policy", "gha-prod"},
			},
		},
	}

	err := createClusterIfNotExist(client, cluster)
	expErr := "failed to create shipa cluster: Framework does not exist."
	assert.EqualError(t, err, expErr)
}

func Test_createFrameworkIfNotExist(t *testing.T) {
	client, _ := shipa.New()
	client.SetDebugMode(true)

	framework := &shipa.PoolConfig{
		Name: "test-fr-1",
		Resources: &shipa.PoolResources{
			General: &shipa.PoolGeneral{
				Setup: &shipa.PoolSetup{
					Default:     true,
					Public:      true,
					Provisioner: "kubernetes",
				},
				Plan: &shipa.PoolPlan{
					Name: "shipa-plan",
				},
				Security: &shipa.PoolSecurity{
					DisableScan: true,
					IgnoreComponents: []string{
						"busybox", "bash", "curl", "dpkg",
					},
				},
				Router: "nginx",
				ContainerPolicy: &shipa.PoolContainerPolicy{
					AllowedHosts: []string{
						"docker.io",
					},
				},
				NodeSelectors: &shipa.NodeSelectors{
					Terms: &shipa.NodeSelectorsTerms{
						Environment: "team1",
						OS:          "linux",
					},
					Strict: true,
				},
				PodAutoScaler: &shipa.PodAutoScaler{
					MinReplicas:                    1,
					MaxReplicas:                    10,
					TargetCPUUtilizationPercentage: 50,
					DisableAppOverride:             true,
				},
				DomainPolicy: &shipa.DomainPolicy{
					AllowedCnames: []string{
						"*.example.com", "*.acme.bar",
					},
				},
				AppAutoDiscovery: &shipa.AppAutoDiscovery{
					AppSelector: []*shipa.AppSelectorLabels{
						{Label: "app"},
					},
				},
				NetworkPolicy: &shipa.PoolNetworkPolicy{
					Ingress: &shipa.NetworkPolicyConfig{
						PolicyMode: "allow-all",
					},
					Egress: &shipa.NetworkPolicyConfig{
						PolicyMode: "allow-all",
					},
					DisableAppPolicies: false,
				},
			},
		},
	}

	err := createFrameworkIfNotExist(client, framework)
	assert.NoError(t, err)
}
