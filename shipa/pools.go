package shipa

import (
	"context"
	"errors"
)

// Pool - represents Shipa pool
type Pool struct {
	Name        string   `json:"name"`
	Default     bool     `json:"default"`
	Provisioner string   `json:"provisioner"`
	Public      bool     `json:"public"`
	Teams       []string `json:"teams,omitempty"`
	Allowed     *Allowed `json:"allowed,omitempty"`
}

// Allowed - part of Pool
type Allowed struct {
	Driver []string `json:"driver,omitempty"`
	Plan   []string `json:"plan,omitempty"`
	Team   []string `json:"team,omitempty"`
}

// CreatePoolRequest - request to create Pool
type CreatePoolRequest struct {
	Name        string `json:"name"`
	Default     bool   `json:"default"`
	Provisioner string `json:"provisioner"`
	Public      bool   `json:"public"`
	Force       bool   `json:"force"`
}

// UpdatePoolRequest - update request for Pool
type UpdatePoolRequest struct {
	Name    string `json:"-"`
	Default bool   `json:"default"`
	Public  bool   `json:"public"`
	Force   bool   `json:"force"`
}

// GetPool - retrieves pool by name
func (c *Client) GetPool(ctx context.Context, name string) (*Pool, error) {
	pools, err := c.ListPools(ctx)
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		if pool.Name == name {
			return pool, nil
		}
	}

	return nil, errors.New("framework not found")
}

// ListPools - lists all pools
func (c *Client) ListPools(ctx context.Context) ([]*Pool, error) {
	pools := make([]*Pool, 0)
	err := c.get(ctx, &pools, apiPools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

// CreatePool - creates pool
func (c *Client) CreatePool(ctx context.Context, req *CreatePoolRequest) error {
	return c.post(ctx, req, apiPools)
}

// UpdatePool - updates pool
func (c *Client) UpdatePool(ctx context.Context, req *UpdatePoolRequest) error {
	return c.put(ctx, req, apiPools, req.Name)
}

// DeletePool - deletes pool
func (c *Client) DeletePool(ctx context.Context, name string) error {
	return c.delete(ctx, apiPools, name)
}
