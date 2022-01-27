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
