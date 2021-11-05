package shipa

import "context"

// PoolConfig - represents Shipa pool
type PoolConfig struct {
	Name      string         `json:"shipaFramework"`
	Resources *PoolResources `json:"resources,omitempty"`
}

// PoolResources - describes main pool configuration part
type PoolResources struct {
	General *PoolGeneral `json:"general,omitempty"`
	Node    *PoolNode    `json:"shipaNode,omitempty"`
}

// PoolNode - part of PoolResources configurations
type PoolNode struct {
	Drivers   []string       `json:"drivers,omitempty"`
	AutoScale *PoolAutoScale `json:"autoScale,omitempty"`
}

// PoolAutoScale - part of PoolNode
type PoolAutoScale struct {
	MaxContainer int     `json:"maxContainer"`
	MaxMemory    int     `json:"maxMemory"`
	ScaleDown    float64 `json:"scaleDown"`
	Rebalance    bool    `json:"rebalance"`
}

// PoolGeneral - pool general configuration
type PoolGeneral struct {
	Setup           *PoolSetup           `json:"setup,omitempty"`
	Plan            *PoolPlan            `json:"plan,omitempty"`
	Security        *PoolSecurity        `json:"security,omitempty"`
	Access          *PoolServiceAccess   `json:"access,omitempty"`
	Services        *PoolServiceAccess   `json:"services,omitempty"`
	Router          string               `json:"router,omitempty"`
	Volumes         []string             `json:"volumes,omitempty"`
	AppQuota        *PoolAppQuota        `json:"appQuota,omitempty"`
	ContainerPolicy *PoolContainerPolicy `json:"containerPolicy,omitempty"`
	NetworkPolicy   *PoolNetworkPolicy   `json:"networkPolicy,omitempty"`
}

// PoolContainerPolicy - part of PoolGeneral object
type PoolContainerPolicy struct {
	AllowedHosts []string `json:"allowedHosts,omitempty"`
}

// PoolAppQuota - part of PoolGeneral object
type PoolAppQuota struct {
	Limit string `json:"limit,omitempty"`
}

// PoolServiceAccess - part of PoolGeneral object
type PoolServiceAccess struct {
	Append    []string `json:"append,omitempty"`
	Blacklist []string `json:"blacklist,omitempty"`
}

// PoolSecurity - part of PoolGeneral object
type PoolSecurity struct {
	DisableScan        bool     `json:"disableScan"`
	ScanPlatformLayers bool     `json:"scanPlatformLayers"`
	IgnoreComponents   []string `json:"ignoreComponents,omitempty"`
	IgnoreCVES         []string `json:"ignoreCves,omitempty"`
}

// PoolPlan - part of PoolGeneral object
type PoolPlan struct {
	Name string `json:"name,omitempty"`
}

// PoolSetup - part of PoolGeneral object
type PoolSetup struct {
	Default             bool   `json:"default"`
	Public              bool   `json:"public"`
	Provisioner         string `json:"provisioner,omitempty"`
	KubernetesNamespace string `json:"kubernetesNamespace,omitempty"`
}

// PoolNetworkPolicy - part of PoolGeneral object
type PoolNetworkPolicy struct {
	Ingress            *NetworkPolicyConfig `json:"ingress,omitempty"`
	Egress             *NetworkPolicyConfig `json:"egress,omitempty"`
	DisableAppPolicies bool                 `json:"disableAppPolicies"`
}

// NetworkPolicyConfig - part of PoolNetworkPolicy
type NetworkPolicyConfig struct {
	PolicyMode        string               `json:"policy_mode,omitempty"`
	CustomRules       []*NetworkPolicyRule `json:"custom_rules,omitempty"`
	ShipaRules        []*NetworkPolicyRule `json:"shipa_rules,omitempty"`
	ShipaRulesEnabled []string             `json:"shipa_rules_enabled,omitempty"`
}

// NetworkPolicyRule - part of NetworkPolicy
type NetworkPolicyRule struct {
	ID                string         `json:"id,omitempty"`
	Enabled           bool           `json:"enabled"`
	Description       string         `json:"description,omitempty"`
	Ports             []*NetworkPort `json:"ports,omitempty"`
	Peers             []*NetworkPeer `json:"peers,omitempty"`
	AllowedApps       []string       `json:"allowed_apps,omitempty"`
	AllowedPools      []string       `json:"allowed_pools,omitempty"`
	AllowedFrameworks []string       `json:"allowed_frameworks,omitempty"` // uses for transferring data from crossplane
}

// NetworkPort - part of NetworkPolicyRule
type NetworkPort struct {
	Protocol string `json:"protocol,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// NetworkPeer - part of NetworkPolicyRule
type NetworkPeer struct {
	PodSelector       *NetworkPeerSelector `json:"podSelector,omitempty"`
	NamespaceSelector *NetworkPeerSelector `json:"namespaceSelector,omitempty"`
	IPBlock           []string             `json:"ipBlock,omitempty"`
}

// NetworkPeerSelector - part of NetworkPeer
type NetworkPeerSelector struct {
	MatchLabels      map[string]string     `json:"matchLabels,omitempty"`
	MatchExpressions []*SelectorExpression `json:"matchExpressions,omitempty"`
}

// SelectorExpression - part of NetworkPeerSelector
type SelectorExpression struct {
	Key      string   `json:"key,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Values   []string `json:"values,omitempty"`
}

func (n *NetworkPolicyRule) updateAllowedPools() {
	if len(n.AllowedFrameworks) > 0 {
		n.AllowedPools = n.AllowedFrameworks
		n.AllowedFrameworks = nil
	}
}

func (n *NetworkPolicyConfig) updateAllowedPools() {
	for i := range n.CustomRules {
		n.CustomRules[i].updateAllowedPools()
	}
	for i := range n.ShipaRules {
		n.ShipaRules[i].updateAllowedPools()
	}
}

func (p *PoolConfig) updateAllowedPools() {
	if p.Resources != nil && p.Resources.General != nil && p.Resources.General.NetworkPolicy != nil {
		p.Resources.General.NetworkPolicy.updateAllowedPools()
	}
}

func (p *PoolNetworkPolicy) updateAllowedPools() {
	if p.Ingress != nil {
		p.Ingress.updateAllowedPools()
	}
	if p.Egress != nil {
		p.Egress.updateAllowedPools()
	}
}

// GetPoolConfig - retrieves pool
func (c *Client) GetPoolConfig(ctx context.Context, name string) (*PoolConfig, error) {
	poolConfig := &PoolConfig{}
	err := c.get(ctx, poolConfig, apiPoolsConfig, name)
	if err != nil {
		return nil, err
	}

	return poolConfig, nil
}

// CreatePoolConfig - creates pool
func (c *Client) CreatePoolConfig(ctx context.Context, pool *PoolConfig) error {
	pool.updateAllowedPools()
	return c.post(ctx, pool, apiPoolsConfig)
}

// UpdatePoolConfig - updates pool
func (c *Client) UpdatePoolConfig(ctx context.Context, req *PoolConfig) error {
	req.updateAllowedPools()
	return c.put(ctx, req, apiPoolsConfig)
}
