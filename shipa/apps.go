package shipa

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// UpdateAppRequest - request for App update
type UpdateAppRequest struct {
	Pool        string   `json:"pool,omitempty"`
	TeamOwner   string   `json:"teamowner,omitempty"`
	Description string   `json:"description,omitempty"`
	Plan        string   `json:"plan,omitempty"`
	Platform    string   `json:"platform,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// NewUpdateAppRequest - converts App object to UpdateAppRequest
func NewUpdateAppRequest(a *App) *UpdateAppRequest {
	var plan string
	if a.Plan != nil {
		plan = a.Plan.Name
	}

	return &UpdateAppRequest{
		Pool:        a.Pool,
		TeamOwner:   a.TeamOwner,
		Description: a.Description,
		Plan:        plan,
		Platform:    a.Platform,
		Tags:        a.Tags,
	}
}

// App - represents shipa app
type App struct {
	Name        string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Pool        string        `json:"pool,omitempty" yaml:"framework,omitempty"`
	TeamOwner   string        `json:"teamowner,omitempty" yaml:"teamowner,omitempty"`
	Plan        *Plan         `json:"plan,omitempty" yaml:"plan,omitempty"`
	Units       []*Unit       `json:"units,omitempty" yaml:"units,omitempty"`
	Cname       []string      `json:"cname,omitempty" yaml:"cname,omitempty"`
	IP          string        `json:"ip,omitempty" yaml:"ip,omitempty"`
	Org         string        `json:"org,omitempty" yaml:"org,omitempty"`
	Entrypoints []*Entrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Routers     []*Router     `json:"routers,omitempty" yaml:"routers,omitempty"`
	Lock        *Lock         `json:"lock,omitempty" yaml:"lock,omitempty"`
	Tags        []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Platform    string        `json:"platform,omitempty" yaml:"platform,omitempty"`
	Status      string        `json:"status,omitempty" yaml:"status,omitempty"`
	Error       string        `json:"error,omitempty" yaml:"error,omitempty"` // not shows in API response
}

// Plan - part of App object
type Plan struct {
	Name     string   `json:"name,omitempty" yaml:"name,omitempty"`
	Memory   int64    `json:"memory" yaml:"memory"`
	Swap     int64    `json:"swap" yaml:"swap"`
	CPUShare int64    `json:"cpushare" yaml:"cpushare"`
	Default  bool     `json:"default" yaml:"default"`
	Public   bool     `json:"public" yaml:"public"`
	Org      string   `json:"org,omitempty" yaml:"org,omitempty"`
	Teams    []string `json:"teams,omitempty" yaml:"teams,omitempty"`
}

// CreatePlanRequest - create request for Plan
type CreatePlanRequest struct {
	Name     string   `json:"name,omitempty"`
	Memory   string   `json:"memory"`
	Swap     string   `json:"swap"`
	CPUShare int64    `json:"cpushare"`
	Default  bool     `json:"default"`
	Public   bool     `json:"public"`
	Org      string   `json:"org,omitempty"`
	Teams    []string `json:"teams,omitempty"`
}

// BytesToHuman - converts number in bytes to shorten form
func BytesToHuman(input int64) string {
	nBytes := int64(1024)
	items := []string{"K", "M", "G"}

	if input < nBytes {
		return strconv.FormatInt(input, 10)
	}

	for _, k := range items {
		input /= nBytes
		if input < nBytes {
			return fmt.Sprintf("%d%s", input, k)
		}
	}

	return fmt.Sprintf("%d%s", input, items[len(items)-1])
}

// Unit - part of App object
type Unit struct {
	ID          string   `json:"ID,omitempty" yaml:"ID,omitempty"`
	Name        string   `json:"Name,omitempty" yaml:"Name,omitempty"`
	AppName     string   `json:"AppName,omitempty" yaml:"AppName,omitempty"`
	ProcessName string   `json:"ProcessName,omitempty" yaml:"ProcessName,omitempty"`
	Type        string   `json:"Type,omitempty" yaml:"Type,omitempty"`
	IP          string   `json:"IP,omitempty" yaml:"IP,omitempty"`
	Status      string   `json:"Status,omitempty" yaml:"Status,omitempty"`
	Version     string   `json:"Version,omitempty" yaml:"Version,omitempty"`
	Org         string   `json:"Org,omitempty" yaml:"Org,omitempty"`
	HostAddr    string   `json:"HostAddr,omitempty" yaml:"HostAddr,omitempty"`
	HostPort    string   `json:"HostPort,omitempty" yaml:"HostPort,omitempty"`
	Address     *Address `json:"Address,omitempty" yaml:"Address,omitempty"`
}

// Address - part of Unit object
type Address struct {
	Scheme      string `json:"Scheme,omitempty" yaml:"Scheme,omitempty"`
	Host        string `json:"Host,omitempty" yaml:"Host,omitempty"`
	Opaque      string `json:"Opaque,omitempty" yaml:"Opaque,omitempty"`
	User        string `json:"User,omitempty" yaml:"User,omitempty"`
	Path        string `json:"Path,omitempty" yaml:"Path,omitempty"`
	RawPath     string `json:"RawPath,omitempty" yaml:"RawPath,omitempty"`
	ForceQuery  bool   `json:"ForceQuery" yaml:"ForceQuery"`
	RawQuery    string `json:"RawQuery,omitempty" yaml:"RawQuery,omitempty"`
	Fragment    string `json:"Fragment,omitempty" yaml:"Fragment,omitempty"`
	RawFragment string `json:"RawFragment,omitempty" yaml:"RawFragment,omitempty"`
}

// Entrypoint - part of App object
type Entrypoint struct {
	Cname  string `json:"cname,omitempty" yaml:"cname,omitempty"`
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
}

// Router - part of App object
type Router struct {
	Name    string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Opts    map[string]interface{} `json:"opts,omitempty" yaml:"opts,omitempty"`
	Type    string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Address string                 `json:"address,omitempty" yaml:"address,omitempty"`
	Default bool                   `json:"default" yaml:"default"` // not show in API response
}

// Lock - part of App object
type Lock struct {
	Locked      bool   `json:"Locked" yaml:"Locked"`
	Reason      string `json:"Reason,omitempty" yaml:"Reason,omitempty"`
	Owner       string `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	AcquireDate string `json:"AcquireDate,omitempty" yaml:"AcquireDate,omitempty"`
}

// ListApps - retrieves all apps
func (c *Client) ListApps(ctx context.Context) ([]*App, error) {
	apps := make([]*App, 0)
	err := c.get(ctx, &apps, apiApps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// GetApp - retrieves app
func (c *Client) GetApp(ctx context.Context, name string) (*App, error) {
	app := &App{}
	err := c.get(ctx, app, apiApps, name)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// CreateApp - creates app
func (c *Client) CreateApp(ctx context.Context, app *App) error {
	return c.post(ctx, app, apiApps)
}

// UpdateApp - updates app
func (c *Client) UpdateApp(ctx context.Context, name string, app *UpdateAppRequest) error {
	return c.put(ctx, app, apiApps, name)
}

// DeleteApp - delets app
func (c *Client) DeleteApp(ctx context.Context, name string) error {
	return c.delete(ctx, apiApps, name)
}

// AppEnv represents application env variable
type AppEnv struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// CreateAppEnv - request to create AppEnv
type CreateAppEnv struct {
	App       string    `yaml:"app"`
	Envs      []*AppEnv `json:"envs" yaml:"envs"`
	NoRestart bool      `json:"norestart" yaml:"norestart"`
	Private   bool      `json:"private" yaml:"private"`
}

// CreateAppEnvs - create app envs
func (c *Client) CreateAppEnvs(ctx context.Context, req *CreateAppEnv) error {
	return c.post(ctx, req, apiAppEnvs(req.App))
}

// GetAppEnvs - retrieves app envs
func (c *Client) GetAppEnvs(ctx context.Context, appName string) ([]*AppEnv, error) {
	envs := make([]*AppEnv, 0)
	err := c.get(ctx, &envs, apiAppEnvs(appName))
	if err != nil {
		return nil, err
	}

	return envs, nil
}

// DeleteAppEnvs - deletes app env
func (c *Client) DeleteAppEnvs(ctx context.Context, req *CreateAppEnv) error {
	params := []*QueryParam{
		{Key: "norestart", Val: req.NoRestart},
	}
	for _, p := range req.Envs {
		params = append(params, &QueryParam{Key: "env", Val: p.Name})
	}

	if len(params) > 1 {
		return c.deleteWithParams(ctx, params, apiAppEnvs(req.App))
	}

	return nil
}

// AppCname - represents app cname
type AppCname struct {
	Cname   string `json:"cname"`
	Encrypt bool   `json:"encrypt"`
}

// CreateAppCname - allows to create app cname
func (c *Client) CreateAppCname(ctx context.Context, appName string, req *AppCname) error {
	return c.post(ctx, req, apiAppCname(appName))
}

// UpdateAppCname - allows to update app cname
func (c *Client) UpdateAppCname(ctx context.Context, appName string, req *AppCname) error {
	return c.put(ctx, req, apiAppCname(appName))
}

// DeleteAppCname - deletes app cname
func (c *Client) DeleteAppCname(ctx context.Context, appName string, req *AppCname) error {
	return c.deleteWithPayload(ctx, req, nil, apiAppCname(appName))
}

// AppDeploy - represents app deploy object
type AppDeploy struct {
	App            string `yaml:"app"`
	Image          string `json:"image" yaml:"image"`
	PrivateImage   bool   `json:"private-image,omitempty" yaml:"private-image,omitempty"`
	RegistryUser   string `json:"registry-user,omitempty" yaml:"registry-user,omitempty"`
	RegistrySecret string `json:"registry-secret,omitempty" yaml:"registry-secret,omitempty"`
	Steps          int64  `json:"steps,omitempty" yaml:"steps,omitempty"`
	StepWeight     int64  `json:"step-weight,omitempty" yaml:"step-weight,omitempty"`
	StepInterval   string `json:"step-interval,omitempty" yaml:"step-interval,omitempty"`
	Port           int64  `json:"port,omitempty" yaml:"port,omitempty"`
	Detach         bool   `json:"detach" yaml:"detach"`
	Message        string `json:"message,omitempty" yaml:"message,omitempty"`
	ShipaYaml      string `json:"shipayaml,omitempty" yaml:"shipayaml,omitempty"`
}

// DeployApp - sends request to deploy app with giving parameters
func (c *Client) DeployApp(ctx context.Context, req *AppDeploy) error {
	params := map[string]string{
		"image": req.Image,
	}
	if req.PrivateImage {
		params["private-image"] = "true"
		params["registry-user"] = req.RegistryUser
		params["registry-secret"] = req.RegistrySecret
	}
	if req.Steps > 0 {
		params["steps"] = strconv.FormatInt(req.Steps, 10)
	}
	if req.StepWeight > 0 {
		params["step-weight"] = strconv.FormatInt(req.StepWeight, 10)
	}

	if req.StepInterval != "" {
		interval, err := parseStepInterval(req.StepInterval)
		if err != nil {
			return err
		}
		params["step-interval"] = interval
	}

	if req.Port > 0 {
		params["port-number"] = strconv.FormatInt(req.Port, 10)
		params["port-protocol"] = "TCP"
	}

	if req.Detach {
		params["detach"] = "true"
	}
	if req.Message != "" {
		params["message"] = req.Message
	}
	if req.ShipaYaml != "" {
		yamlContent, err := getShipaYamlBase64Enc(req.ShipaYaml)
		if err != nil {
			return err
		}
		params["shipayaml"] = yamlContent
	}

	return c.postURLEncoded(ctx, params, apiAppDeploy(req.App))
}

func getShipaYamlBase64Enc(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func parseStepInterval(input string) (string, error) {
	if input == "" {
		return "0", nil
	}

	interval, err := time.ParseDuration(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse step-interval="+input)
	}

	return strconv.FormatInt(int64(interval.Seconds()), 10), nil
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
