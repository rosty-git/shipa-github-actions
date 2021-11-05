package shipa

import (
	"context"
	"errors"
)

var (
	// ErrUserNotFound - uses when user not found
	ErrUserNotFound = errors.New("user not found")
)

// User - represents Shipa user
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetUser - retrieves user
func (c *Client) GetUser(ctx context.Context, email string) (*User, error) {
	users, err := c.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// ListUsers - lists all user
func (c *Client) ListUsers(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)
	err := c.get(ctx, &users, apiUsers)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser - create new user
func (c *Client) CreateUser(ctx context.Context, req *User) error {
	return c.post(ctx, req, apiUsers)
}

// DeleteUser - deletes user
func (c *Client) DeleteUser(ctx context.Context, email string) error {
	return nil

	// TODO: uncomment after delete user will be fixed
	// params := map[string]string{
	// 	"email": email,
	// }
	// return c.deleteWithParams(ctx, params, apiUsers)
}
