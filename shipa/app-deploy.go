package shipa

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// AppDeploy - represents new app deploy object
type AppDeploy struct {
	App            string                   `json:"-" yaml:"app"`
	Image          string                   `json:"image" yaml:"image"`
	AppConfig      *AppDeployConfig         `json:"appConfig" yaml:"appConfig"`
	CanarySettings *AppDeployCanarySettings `json:"canarySettings,omitempty" yaml:"canarySettings,omitempty"`
	PodAutoScaler  *AppDeployPodAutoScaler  `json:"podAutoScaler,omitempty" yaml:"podAutoScaler,omitempty"`
	Port           *AppDeployPort           `json:"port,omitempty" yaml:"port,omitempty"`
	Registry       *AppDeployRegistry       `json:"registry,omitempty" yaml:"registry,omitempty"`
	Volumes        []*AppDeployVolume       `json:"volumesToBind,omitempty" yaml:"volumesToBind,omitempty"`
}

// AppDeployConfig - represents app deploy config
type AppDeployConfig struct {
	Team        string   `json:"team" yaml:"team"`
	Framework   string   `json:"framework" yaml:"framework"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Env         []string `json:"env,omitempty" yaml:"env,omitempty"`
	Plan        string   `json:"plan,omitempty" yaml:"plan,omitempty"`
	Router      string   `json:"router,omitempty" yaml:"router,omitempty"`
	Tags        []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// AppDeployCanarySettings - represents app deploy canary settings
type AppDeployCanarySettings struct {
	StepInterval int64 `json:"stepInterval" yaml:"stepInterval"`
	StepWeight   int64 `json:"stepWeight" yaml:"stepWeight"`
	Steps        int64 `json:"steps" yaml:"steps"`
}

// AppDeployPodAutoScaler - represents app deploy auto scaler
type AppDeployPodAutoScaler struct {
	MaxReplicas                    int64 `json:"maxReplicas" yaml:"maxReplicas"`
	MinReplicas                    int64 `json:"minReplicas" yaml:"minReplicas"`
	TargetCPUUtilizationPercentage int64 `json:"targetCPUUtilizationPercentage" yaml:"targetCPUUtilizationPercentage"`
}

// AppDeployPort - represents app deploy port
type AppDeployPort struct {
	Number   int64  `json:"number" yaml:"number"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

// AppDeployRegistry - represents app deploy registry
type AppDeployRegistry struct {
	User   string `json:"user" yaml:"user"`
	Secret string `json:"secret" yaml:"secret"`
}

// AppDeployVolume - represents app deploy volume
type AppDeployVolume struct {
	Name    string         `json:"volumeName" yaml:"volumeName"`
	Path    string         `json:"volumeMountPath" yaml:"volumeMountPath"`
	Options *VolumeOptions `json:"volumeMountOptions" yaml:"volumeMountOptions"`
}

// VolumeOptions - represents additional volume options
type VolumeOptions struct {
	Prop1 string `json:"additionalProp1" yaml:"additionalProp1"`
	Prop2 string `json:"additionalProp2" yaml:"additionalProp2"`
	Prop3 string `json:"additionalProp3" yaml:"additionalProp3"`
}

// DeployApp - sends request to deploy app with giving parameters
func (c *Client) DeployApp(ctx context.Context, req *AppDeploy) error {
	body, statusCode, err := c.updateRequest(ctx, "POST", req, apiAppDeploy(req.App))
	if err != nil {
		return err
	}

	return validateAppDeployResponse(body, statusCode)
}

func validateAppDeployResponse(body []byte, statusCode int) error {
	fmt.Println("### Deploy app RESP:", string(body))

	if statusCode != http.StatusAccepted && statusCode != http.StatusCreated && statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}

	if bytes.Contains(body, []byte("There are vulnerabilities!")) {
		return errors.New("found vulnerabilities")
	}

	if bytes.Contains(bytes.ToLower(body), []byte(`"error"`)) {
		return fmt.Errorf("app deploy failed, body: %s", body)
	}

	return nil
}

// AppDeployment - represents information about app deployments
type AppDeployment struct {
	ID          string `json:"ID"`
	App         string `json:"App"`
	Active      bool   `json:"Active"`
	Image       string `json:"Image"`
	Version     string `json:"Version"`
	Origin      string `json:"Origin,omitempty"`
	Message     string `json:"Message,omitempty"`
	Commit      string `json:"Commit,omitempty"`
	User        string `json:"User,omitempty"`
	Timestamp   string `json:"Timestamp,omitempty"`
	Error       string `json:"Error,omitempty"`
	CanRollback bool   `json:"CanRollback"`
	Org         string `json:"Org,omitempty"`
}

// ListAppDeployments - lists app deployments
func (c *Client) ListAppDeployments(ctx context.Context, appName string) ([]*AppDeployment, error) {
	deployments := make([]*AppDeployment, 0)
	err := c.get(ctx, &deployments, apiAppDeployments(appName))
	if err != nil {
		return nil, err
	}

	return deployments, nil
}
