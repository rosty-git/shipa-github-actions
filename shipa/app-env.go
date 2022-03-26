package shipa

import "context"

// AppEnv represents application env variable
type AppEnv struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// CreateAppEnv - request to create AppEnv
type CreateAppEnv struct {
	App       string    `json:"-" yaml:"app"`
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
