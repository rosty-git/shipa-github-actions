package shipa

import "context"

// AppCname - represents app cname
type AppCname struct {
	App       string `json:"-" yaml:"app"`
	Cname     string `json:"cname" yaml:"cname"`
	Scheme    string `json:"scheme"`
	Encrypted bool   `json:"-" yaml:"encrypted"`
}

func (a *AppCname) setScheme() {
	a.Scheme = "http"
	if a.Encrypted {
		a.Scheme = "https"
	}
}

// CreateAppCname - allows to create app cname
func (c *Client) CreateAppCname(ctx context.Context, req *AppCname) error {
	req.setScheme()
	return c.post(ctx, req, apiAppCname(req.App))
}

// UpdateAppCname - allows to update app cname
func (c *Client) UpdateAppCname(ctx context.Context, req *AppCname) error {
	req.setScheme()
	return c.put(ctx, req, apiAppCname(req.App))
}

// DeleteCnameRequest - request payload to delete cname
type DeleteCnameRequest struct {
	App   string
	Cname []string `json:"cnames"`
}

// DeleteAppCname - deletes app cname
func (c *Client) DeleteAppCname(ctx context.Context, req *DeleteCnameRequest) error {
	return c.deleteWithPayload(ctx, req, nil, apiAppCname(req.App))
}
