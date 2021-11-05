package shipa

import "context"

// Role - represents Shipa role
type Role struct {
	Name        string `json:"name"`
	Context     string `json:"context"`
	Description string `json:"description,omitempty"`
}

// GetRole - retreives role
func (c *Client) GetRole(ctx context.Context, name string) (*Role, error) {
	role := &Role{}
	err := c.get(ctx, role, apiRoles, name)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// CreateRole - creates role
func (c *Client) CreateRole(ctx context.Context, req *Role) error {
	return c.post(ctx, req, apiRoles)
}

// DeleteRole - deletes role
func (c *Client) DeleteRole(ctx context.Context, name string) error {
	return c.delete(ctx, apiRoles, name)
}

// Email - email struct
type Email struct {
	Email string `json:"email"`
}

// AssociateRoleToUser - adds role to user
func (c *Client) AssociateRoleToUser(ctx context.Context, role, email string) error {
	return c.post(ctx, &Email{email}, apiRoleUser(role))
}

// DisassociateRoleFromUser - removes role from user
func (c *Client) DisassociateRoleFromUser(ctx context.Context, role, email string) error {
	return c.deleteWithPayload(ctx, &Email{email}, nil, apiRoleUser(role), email)
}
