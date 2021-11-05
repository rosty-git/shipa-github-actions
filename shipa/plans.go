package shipa

import (
	"context"
	"errors"
)

// GetPlan - retrieves plan by name
func (c *Client) GetPlan(ctx context.Context, name string) (*Plan, error) {
	plans, err := c.ListPlans(ctx)
	if err != nil {
		return nil, err
	}
	for _, plan := range plans {
		if plan.Name == name {
			return plan, nil
		}
	}

	return nil, errors.New("plan not found")
}

// ListPlans - list all plans
func (c *Client) ListPlans(ctx context.Context) ([]*Plan, error) {
	plans := make([]*Plan, 0)
	err := c.get(ctx, &plans, apiPlans)
	if err != nil {
		return nil, err
	}

	return plans, nil
}

// CreatePlan - creates plan
func (c *Client) CreatePlan(ctx context.Context, req *CreatePlanRequest) error {
	return c.post(ctx, req, apiPlans)
}

// DeletePlan - deletes plan
func (c *Client) DeletePlan(ctx context.Context, name string) error {
	return c.delete(ctx, apiPlans, name)
}
