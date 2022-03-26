package shipa

import (
	"context"
	"fmt"
	"strconv"
)

// CreateAppRequest - request to create an App
type CreateAppRequest struct {
	Name      string   `json:"name" yaml:"name,omitempty"`
	Pool      string   `json:"pool,omitempty" yaml:"framework,omitempty"`
	TeamOwner string   `json:"teamOwner,omitempty" yaml:"teamOwner,omitempty"`
	Plan      string   `json:"plan,omitempty" yaml:"plan,omitempty"`
	Tags      []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

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
func (c *Client) CreateApp(ctx context.Context, app *CreateAppRequest) error {
	return c.post(ctx, app, apiApps)
}

// UpdateApp - updates app
func (c *Client) UpdateApp(ctx context.Context, name string, app *UpdateAppRequest) error {
	return c.put(ctx, app, apiApps, name)
}

// DeleteApp - deletes app
func (c *Client) DeleteApp(ctx context.Context, name string) error {
	return c.delete(ctx, apiApps, name)
}
