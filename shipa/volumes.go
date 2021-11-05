package shipa

import "context"

// Volume - represents Shipa Volume
type Volume struct {
	Name        string         `json:"Name"`
	Capacity    string         `json:"Capacity"`
	TeamOwner   string         `json:"TeamOwner"`
	Pool        string         `json:"Pool"`
	AccessModes string         `json:"AccessModes"`
	Plan        VolumePlanName `json:"Plan"`
}

// VolumePlanName - internal struct for Shipa Volume
type VolumePlanName struct {
	Name string `json:"Name"`
}

// VolumeBinding - represents Shipa volume binding
type VolumeBinding struct {
	Volume string `json:"-"`

	App        string `json:"App"`
	MountPoint string `json:"MountPoint"`
	// optional
	NoRestart bool `json:"NoRestart"`
}

// CreateVolume - creates volume
func (c *Client) CreateVolume(ctx context.Context, req *Volume) error {
	return c.post(ctx, req, apiVolumes)
}

// GetVolume - retrieve volume by name
func (c *Client) GetVolume(ctx context.Context, name string) (*Volume, error) {
	resp := &Volume{}
	err := c.get(ctx, resp, apiVolumes, name)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateVolume - updates volume
func (c *Client) UpdateVolume(ctx context.Context, req *Volume) error {
	return c.post(ctx, req, apiVolumes, req.Name)
}

// DeleteVolume - deletes volume by name
func (c *Client) DeleteVolume(ctx context.Context, name string) error {
	return c.delete(ctx, apiVolumes, name)
}

// BindVolume - binds volume
func (c *Client) BindVolume(ctx context.Context, req *VolumeBinding) error {
	return c.post(ctx, req, apiVolumeBind(req.Volume))
}

// UnbindVolume - unbinds volume
func (c *Client) UnbindVolume(ctx context.Context, req *VolumeBinding) error {
	return c.deleteWithPayload(ctx, req, nil, apiVolumeBind(req.Volume))
}
