package shipa

import "context"

// PoolConfig - represents Shipa pool
type PoolConfig struct {
	Name      string         `json:"shipaFramework" yaml:"name"`
	Resources *PoolResources `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// PoolResources - describes main pool configuration part
type PoolResources struct {
	General *PoolGeneral `json:"general,omitempty" yaml:"general,omitempty"`
	Node    *PoolNode    `json:"shipaNode,omitempty" yaml:"shipaNode,omitempty"`
}

// PoolNode - part of PoolResources configurations
type PoolNode struct {
	Drivers   []string       `json:"drivers,omitempty" yaml:"drivers,omitempty"`
	AutoScale *PoolAutoScale `json:"autoScale,omitempty" yaml:"autoScale,omitempty"`
}

// PoolAutoScale - part of PoolNode
type PoolAutoScale struct {
	MaxContainer int     `json:"maxContainer" yaml:"maxContainer"`
	MaxMemory    int     `json:"maxMemory" yaml:"maxMemory"`
	ScaleDown    float64 `json:"scaleDown" yaml:"scaleDown"`
	Rebalance    bool    `json:"rebalance" yaml:"rebalance"`
}

// PoolGeneral - pool general configuration
type PoolGeneral struct {
	Setup            *PoolSetup           `json:"setup,omitempty" yaml:"setup,omitempty"`
	Plan             *PoolPlan            `json:"plan,omitempty" yaml:"plan,omitempty"`
	Security         *PoolSecurity        `json:"security,omitempty" yaml:"security,omitempty"`
	Access           *PoolServiceAccess   `json:"access,omitempty" yaml:"access,omitempty"`
	Services         *PoolServiceAccess   `json:"services,omitempty" yaml:"services,omitempty"`
	Router           string               `json:"router,omitempty" yaml:"router,omitempty"`
	Volumes          []string             `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	ContainerPolicy  *PoolContainerPolicy `json:"containerPolicy,omitempty" yaml:"containerPolicy,omitempty"`
	NodeSelectors    *NodeSelectors       `json:"nodeSelectors,omitempty" yaml:"nodeSelectors,omitempty"`
	PodAutoScaler    *PodAutoScaler       `json:"podAutoScaler,omitempty" yaml:"podAutoScaler,omitempty"`
	DomainPolicy     *DomainPolicy        `json:"domainPolicy,omitempty" yaml:"domainPolicy,omitempty"`
	AppAutoDiscovery *AppAutoDiscovery    `json:"appAutoDiscovery,omitempty" yaml:"appAutoDiscovery,omitempty"`
	NetworkPolicy    *PoolNetworkPolicy   `json:"networkPolicy,omitempty" yaml:"networkPolicy,omitempty"`
}

// AppAutoDiscovery - part of PoolGeneral object
type AppAutoDiscovery struct {
	AppSelector []*AppSelectorLabels `json:"appSelector,omitempty" yaml:"appSelector,omitempty"`
	Suffix      string               `json:"suffix,omitempty" yaml:"suffix,omitempty"`
}

// AppSelectorLabels - part of AppAutoDiscovery object
type AppSelectorLabels struct {
	Label string `json:"label,omitempty" yaml:"label,omitempty"`
}

// DomainPolicy - part of PoolGeneral object
type DomainPolicy struct {
	AllowedCnames []string `json:"allowedCnames,omitempty" yaml:"allowedCnames,omitempty"`
}

// PodAutoScaler - part of PoolGeneral object
type PodAutoScaler struct {
	MinReplicas                    int  `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas                    int  `json:"maxReplicas" yaml:"maxReplicas"`
	TargetCPUUtilizationPercentage int  `json:"targetCPUUtilizationPercentage" yaml:"targetCPUUtilizationPercentage"`
	DisableAppOverride             bool `json:"disableAppOverride" yaml:"disableAppOverride"`
}

// NodeSelectors - part of PoolGeneral object
type NodeSelectors struct {
	Terms  *NodeSelectorsTerms `json:"terms,omitempty" yaml:"terms,omitempty"`
	Strict bool                `json:"strict" yaml:"strict"`
}

// NodeSelectorsTerms - part of NodeSelectors object
type NodeSelectorsTerms struct {
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
	OS          string `json:"os,omitempty" yaml:"os,omitempty"`
}

// PoolContainerPolicy - part of PoolGeneral object
type PoolContainerPolicy struct {
	AllowedHosts []string `json:"allowedHosts,omitempty" yaml:"allowedHosts,omitempty"`
}

// PoolServiceAccess - part of PoolGeneral object
type PoolServiceAccess struct {
	Append    []string `json:"append,omitempty" yaml:"append,omitempty"`
	Blacklist []string `json:"blacklist,omitempty" yaml:"blacklist,omitempty"`
}

// PoolSecurity - part of PoolGeneral object
type PoolSecurity struct {
	DisableScan        bool     `json:"disableScan" yaml:"disableScan"`
	ScanPlatformLayers bool     `json:"scanPlatformLayers" yaml:"scanPlatformLayers"`
	IgnoreComponents   []string `json:"ignoreComponents,omitempty" yaml:"ignoreComponents,omitempty"`
	IgnoreCVES         []string `json:"ignoreCves,omitempty" yaml:"ignoreCves,omitempty"`
}

// PoolPlan - part of PoolGeneral object
type PoolPlan struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

// PoolSetup - part of PoolGeneral object
type PoolSetup struct {
	Default             bool   `json:"default" yaml:"default"`
	Public              bool   `json:"public" yaml:"public"`
	Provisioner         string `json:"provisioner,omitempty" yaml:"provisioner,omitempty"`
	KubernetesNamespace string `json:"kubernetesNamespace,omitempty" yaml:"kubernetesNamespace,omitempty"`
}

// PoolNetworkPolicy - part of PoolGeneral object
type PoolNetworkPolicy struct {
	Ingress            *NetworkPolicyConfig `json:"ingress,omitempty" yaml:"ingress,omitempty"`
	Egress             *NetworkPolicyConfig `json:"egress,omitempty" yaml:"egress,omitempty"`
	DisableAppPolicies bool                 `json:"disableAppPolicies" yaml:"disableAppPolicies"`
}

// NetworkPolicyConfig - part of PoolNetworkPolicy
type NetworkPolicyConfig struct {
	PolicyMode        string               `json:"policy_mode,omitempty" yaml:"policy_mode,omitempty"`
	CustomRules       []*NetworkPolicyRule `json:"custom_rules,omitempty" yaml:"custom_rules,omitempty"`
	ShipaRules        []*NetworkPolicyRule `json:"shipa_rules,omitempty" yaml:"shipa_rules,omitempty"`
	ShipaRulesEnabled []string             `json:"shipa_rules_enabled,omitempty" yaml:"shipa_rules_enabled,omitempty"`
}

// NetworkPolicyRule - part of NetworkPolicy
type NetworkPolicyRule struct {
	ID           string         `json:"id,omitempty" yaml:"id,omitempty"`
	Enabled      bool           `json:"enabled" yaml:"enabled"`
	Description  string         `json:"description,omitempty" yaml:"description,omitempty"`
	Ports        []*NetworkPort `json:"ports,omitempty" yaml:"ports,omitempty"`
	Peers        []*NetworkPeer `json:"peers,omitempty" yaml:"peers,omitempty"`
	AllowedApps  []string       `json:"allowed_apps,omitempty" yaml:"allowed_apps,omitempty"`
	AllowedPools []string       `json:"allowed_pools,omitempty" yaml:"allowed_frameworks,omitempty"`
}

// NetworkPort - part of NetworkPolicyRule
type NetworkPort struct {
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
}

// NetworkPeer - part of NetworkPolicyRule
type NetworkPeer struct {
	PodSelector       *NetworkPeerSelector `json:"podSelector,omitempty" yaml:"podSelector,omitempty"`
	NamespaceSelector *NetworkPeerSelector `json:"namespaceSelector,omitempty" yaml:"namespaceSelector,omitempty"`
	IPBlock           []string             `json:"ipBlock,omitempty" yaml:"ipBlock,omitempty"`
}

// NetworkPeerSelector - part of NetworkPeer
type NetworkPeerSelector struct {
	MatchLabels      map[string]string     `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
	MatchExpressions []*SelectorExpression `json:"matchExpressions,omitempty" yaml:"matchExpressions,omitempty"`
}

// SelectorExpression - part of NetworkPeerSelector
type SelectorExpression struct {
	Key      string   `json:"key,omitempty" yaml:"key,omitempty"`
	Operator string   `json:"operator,omitempty" yaml:"operator,omitempty"`
	Values   []string `json:"values,omitempty" yaml:"values,omitempty"`
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
	return c.post(ctx, pool, apiPoolsConfig)
}

// UpdatePoolConfig - updates pool
func (c *Client) UpdatePoolConfig(ctx context.Context, req *PoolConfig) error {
	return c.put(ctx, req, apiPoolsConfig)
}
