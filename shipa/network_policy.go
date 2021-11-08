package shipa

import "context"

// NetworkPolicy - represents Shipa network-policy
type NetworkPolicy struct {
	App string `yaml:"app"`
	Ingress    *NetworkPolicyConfig `json:"ingress,omitempty" yaml:"ingress,omitempty"`
	Egress     *NetworkPolicyConfig `json:"egress,omitempty" yaml:"egress,omitempty"`
	RestartApp bool                 `json:"restart_app" yaml:"restart_app"`
}

// CreateOrUpdateNetworkPolicy - creates or updates network policy
func (c *Client) CreateOrUpdateNetworkPolicy(ctx context.Context, config *NetworkPolicy) error {
	return c.put(ctx, config, apiAppNetworkPolicy(config.App))
}

// DeleteNetworkPolicy - deletes network policy from the given app
func (c *Client) DeleteNetworkPolicy(ctx context.Context, app string) error {
	return c.delete(ctx, apiAppNetworkPolicy(app))
}

// GetNetworkPolicy - get current policy of an app
func (c *Client) GetNetworkPolicy(ctx context.Context, app string) (*NetworkPolicy, error) {
	config := &NetworkPolicy{}
	err := c.get(ctx, config, apiAppNetworkPolicy(app))
	if err != nil {
		return nil, err
	}

	return config, nil
}
