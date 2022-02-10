package shipa

import (
	"context"
	"encoding/json"
)

// Job defines job data
type Job struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Framework    string          `json:"framework"`
	Containers   []*JobContainer `json:"containers"`
	Policy       *JobPolicy      `json:"policy"`
	Cluster      string          `json:"cluster,omitempty"`
	Org          string          `json:"org,omitempty"`
	Owner        string          `json:"owner,omitempty"`
	BackoffLimit int64           `json:"backoffLimit"`
	Completions  int64           `json:"completions"`
	Parallelism  int64           `json:"parallelism"`
	Suspend      bool            `json:"suspend"`
	Description  string          `json:"description,omitempty"`
	Team         string          `json:"team,omitempty"`
	Teams        []string        `json:"teams,omitempty"`
	Type         string          `json:"type,omitempty"`
	Version      string          `json:"version,omitempty"`
	CreatedAt    string          `json:"createdAt,omitempty"`
	DeletedAt    string          `json:"deletedAt,omitempty"`
	UpdatedAt    string          `json:"updatedAt,omitempty"`
}

// JobCreateRequest defines fields to create Job
type JobCreateRequest struct {
	Name       string          `json:"name" yaml:"name"`
	Framework  string          `json:"framework" yaml:"framework"`
	Containers []*JobContainer `json:"containers" yaml:"containers"`
	Policy     *JobPolicy      `json:"policy" yaml:"policy"`

	// optional
	BackoffLimit int64  `json:"backoffLimit" yaml:"backoffLimit"`
	Completions  int64  `json:"completions" yaml:"completions"`
	Parallelism  int64  `json:"parallelism" yaml:"parallelism"`
	Suspend      bool   `json:"suspend" yaml:"suspend"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty"`
	Team         string `json:"team,omitempty" yaml:"team,omitempty"`
	Type         string `json:"type,omitempty" yaml:"type,omitempty"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
}

// JobPolicy defines restart policy
type JobPolicy struct {
	RestartPolicy string `json:"restartPolicy" yaml:"restartPolicy"`
}

// JobContainer defines container
type JobContainer struct {
	Command []string `json:"command" yaml:"command"`
	Image   string   `json:"image" yaml:"image"`
	Name    string   `json:"name" yaml:"name"`
}

// GetJob - retrieves job by id
func (c *Client) GetJob(ctx context.Context, id string) (*Job, error) {
	job := &Job{}
	err := c.get(ctx, job, apiJobs, id)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// ListJobs - list all jobs
func (c *Client) ListJobs(ctx context.Context) ([]*Job, error) {
	jobs := make([]*Job, 0)
	err := c.get(ctx, &jobs, apiJobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// CreateJob - creates job
func (c *Client) CreateJob(ctx context.Context, req *JobCreateRequest) (*Job, error) {
	body, err := c.postWithResult(ctx, req, apiJobs)
	if err != nil {
		return nil, err
	}

	job := &Job{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// DeleteJob - deletes job
func (c *Client) DeleteJob(ctx context.Context, id string) error {
	return c.delete(ctx, apiJobs, id)
}
