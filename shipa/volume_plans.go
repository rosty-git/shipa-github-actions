package shipa

import "context"

// VolumePlan - represents Shipa VolumePlan
type VolumePlan struct {
	Name         string `json:"name"`
	Team         string `json:"team"`
	StorageClass string `json:"storage_class"`
}

func (v *VolumePlan) toCreateRequest() *createVolumePlanRequest {
	return &createVolumePlanRequest{
		Name:         v.Name,
		Teams:        []string{v.Team},
		StorageClass: v.StorageClass,
	}
}

type createVolumePlanRequest struct {
	Name         string   `json:"name"`
	Teams        []string `json:"teams"`
	StorageClass string   `json:"storage_class"`
}

type getVolumePlanResponse struct {
	Name         string   `json:"Name"`
	Teams        []string `json:"Teams"`
	StorageClass string   `json:"StorageClass"`
}

func (r *getVolumePlanResponse) toVolumePlan() *VolumePlan {
	var team string
	if len(r.Teams) > 0 {
		team = r.Teams[0]
	}
	return &VolumePlan{
		Name:         r.Name,
		Team:         team,
		StorageClass: r.StorageClass,
	}
}

// CreateVolumePlan - creates volume plan
func (c *Client) CreateVolumePlan(ctx context.Context, req *VolumePlan) error {
	return c.post(ctx, req.toCreateRequest(), apiVolumePlans)
}

// DeleteVolumePlan - deletes volume plan by name
func (c *Client) DeleteVolumePlan(ctx context.Context, name string) error {
	return c.delete(ctx, apiVolumePlans, name)
}

// GetVolumePlan - retrieves volume plan by name
func (c *Client) GetVolumePlan(ctx context.Context, name string) (*VolumePlan, error) {
	resp := &getVolumePlanResponse{}
	err := c.get(ctx, resp, apiVolumePlans, name)
	if err != nil {
		return nil, err
	}

	return resp.toVolumePlan(), nil
}

// UpdateVolumePlan - updates volume plan
func (c *Client) UpdateVolumePlan(ctx context.Context, req *VolumePlan) error {
	return c.put(ctx, req.toCreateRequest(), apiVolumePlans, req.Name)
}
