package shipa

import "context"

// NetworkPolicy - represents Shipa network-policy
type NetworkPolicy struct {
	Ingress    *NetworkPolicyConfig `json:"ingress,omitempty"`
	Egress     *NetworkPolicyConfig `json:"egress,omitempty"`
	RestartApp bool                 `json:"restart_app"`
}

func (p *NetworkPolicy) updateAllowedPools() {
	if p.Ingress != nil {
		p.Ingress.updateAllowedPools()
	}
	if p.Egress != nil {
		p.Egress.updateAllowedPools()
	}
}

// CreateOrUpdateNetworkPolicy - creates or updates network policy
func (c *Client) CreateOrUpdateNetworkPolicy(ctx context.Context, app string, config *NetworkPolicy) error {
	config.updateAllowedPools()
	return c.put(ctx, config, apiAppNetworkPolicy(app))
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
