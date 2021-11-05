package shipa

import "context"

// Team - represents Shipa team
type Team struct {
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
}

// UpdateTeamRequest - request for team update
type UpdateTeamRequest struct {
	Name string   `json:"newname,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

// GetTeam - retreives team
func (c *Client) GetTeam(ctx context.Context, name string) (*Team, error) {
	team := &Team{}
	err := c.get(ctx, team, apiTeams, name)
	if err != nil {
		return nil, err
	}

	return team, nil
}

// CreateTeam - creates team
func (c *Client) CreateTeam(ctx context.Context, req *Team) error {
	return c.post(ctx, req, apiTeams)
}

// UpdateTeam - updates team
func (c *Client) UpdateTeam(ctx context.Context, name string, req *UpdateTeamRequest) error {
	return c.put(ctx, req, apiTeams, name)
}

// DeleteTeam - deletes team
func (c *Client) DeleteTeam(ctx context.Context, name string) error {
	return c.delete(ctx, apiTeams, name)
}
