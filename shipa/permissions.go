package shipa

import "context"

// Permission - represents Shipa permission
type Permission struct {
	Role        string   `json:"name"`
	Permissions []string `json:"permission"`
}

type getRolePermissions struct {
	Role        string   `json:"name"`
	Context     string   `json:"context"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"scheme_names,omitempty"`
}

// GetPermission - retrieves permission
func (c *Client) GetPermission(ctx context.Context, role string) (*Permission, error) {
	req := &getRolePermissions{}
	err := c.get(ctx, req, apiRoles, role)
	if err != nil {
		return nil, err
	}

	return &Permission{
		Role:        req.Role,
		Permissions: req.Permissions,
	}, nil
}

// CreatePermission - creates permission
func (c *Client) CreatePermission(ctx context.Context, req *Permission) error {
	return c.post(ctx, req, apiRolePermissions(req.Role))
}

// DeletePermission - deletes permission
func (c *Client) DeletePermission(ctx context.Context, role, permission string) error {
	return c.delete(ctx, apiRolePermissions(role), permission)
}
