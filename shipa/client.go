package shipa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// API endpoints
const (
	apiClusters    = "provisioner/clusters"
	apiPoolsConfig = "frameworks-config"
	apiPools       = "pools"
	apiApps        = "apps"
	apiUsers       = "users"
	apiPlans       = "plans"
	apiTeams       = "teams"
	apiRoles       = "roles"
	apiVolumes     = "volumes"
	apiVolumePlans = "volume-plans"
	apiJobs        = "jobs"
)

func apiAppNetworkPolicy(appName string) string {
	return fmt.Sprintf("%s/%s/network-policy", apiApps, appName)
}

func apiAppDeployments(appName string) string {
	return fmt.Sprintf("%s/%s/deployments", apiApps, appName)
}

func apiAppEnvs(appName string) string {
	return fmt.Sprintf("%s/%s/env", apiApps, appName)
}

func apiAppCname(appName string) string {
	return fmt.Sprintf("%s/%s/cname", apiApps, appName)
}

func apiAppDeploy(appName string) string {
	return fmt.Sprintf("%s/%s/deploy", apiApps, appName)
}

func apiRolePermissions(role string) string {
	return fmt.Sprintf("%s/%s/permissions", apiRoles, role)
}

func apiRoleUser(role string) string {
	return fmt.Sprintf("%s/%s/user", apiRoles, role)
}

func apiVolumeBind(volumeName string) string {
	return fmt.Sprintf("%s/%s/bind", apiVolumes, volumeName)
}

// Client - represents shipa client
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	debug      bool
}

// New returns a Shipa client, trying to get host and token from ENVs
func New() (*Client, error) {
	return NewClient(os.Getenv("SHIPA_HOST"), os.Getenv("SHIPA_TOKEN"))
}

// NewClient returns a new Shipa client.
func NewClient(host, token string) (*Client, error) {
	if host == "" {
		return nil, errors.New("shipa client init failed: host can not be empty")
	}

	if token == "" {
		return nil, errors.New("shipa client init failed: token can not be empty")
	}

	c := &Client{
		HostURL:    host,
		HTTPClient: &http.Client{Timeout: 1500 * time.Second},
		Token:      token,
	}

	err := c.testAuthentication()
	if err != nil {
		return nil, fmt.Errorf("shipa client auth failed: %w", err)
	}

	return c, nil
}

func (c *Client) SetDebugMode(debug bool) {
	c.debug = debug
}

func (c *Client) doRequest(req *http.Request) ([]byte, int, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	res, err := c.HTTPClient.Do(req)
	// ignore EOF error, when host drops connection
	if err != nil && strings.HasSuffix(err.Error(), ": EOF") {
		return []byte{}, http.StatusOK, nil
	}

	if err != nil {
		return nil, 0, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, res.StatusCode, err
}

func closeBody(res *http.Response) {
	if !res.Close {
		if err := res.Body.Close(); err != nil {
			fmt.Println("ERR: failed to close response body:", err.Error())
		}
	}
}

func (c *Client) get(ctx context.Context, out interface{}, urlPath ...string) error {
	req, err := c.newRequest(ctx, "GET", nil, urlPath...)
	if err != nil {
		return err
	}

	body, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}
	return json.Unmarshal(body, out)
}

func (c *Client) newURLEncodedRequest(ctx context.Context, method string, params map[string]string, urlPath ...string) (*http.Request, error) {
	URL := strings.Join(append([]string{c.HostURL}, urlPath...), "/")

	if c.debug {
		log.Printf("\n> %s: %s\n", method, URL)
		log.Printf(">>> Payload: %+v\n", params)
	}

	data := url.Values{}
	for key, val := range params {
		data.Set(key, val)
	}

	return http.NewRequestWithContext(ctx, method, URL, strings.NewReader(data.Encode())) // URL-encoded payload
}

func (c *Client) newRequest(ctx context.Context, method string, payload interface{}, urlPath ...string) (*http.Request, error) {
	var body io.Reader
	URL := strings.Join(append([]string{c.HostURL}, urlPath...), "/")

	if c.debug {
		log.Printf("\n> %s: %s\n", method, URL)
	}

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		if c.debug {
			log.Printf(">>> Payload: %s\n", string(data))
		}

		body = bytes.NewBuffer(data)
	}

	return http.NewRequestWithContext(ctx, method, URL, body)
}

func (c *Client) newRequestWithParams(ctx context.Context, method string, payload interface{}, urlPath []string, params map[string]string) (*http.Request, error) {
	var body io.Reader
	URL := strings.Join(append([]string{c.HostURL}, urlPath...), "/")

	paramValues := make([]string, 0)
	for key, val := range params {
		paramValues = append(paramValues, fmt.Sprintf("%s=%s", key, val))
	}
	paramsStr := strings.Join(paramValues, "&")

	if paramsStr != "" {
		URL = fmt.Sprintf("%s?%s", URL, paramsStr)
	}

	if c.debug {
		log.Printf("\n> %s: %s\n", method, URL)
	}

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		if c.debug {
			log.Printf(">>> Payload: %s\n", string(data))
		}

		body = bytes.NewBuffer(data)
	}

	return http.NewRequestWithContext(ctx, method, URL, body)
}

func (c *Client) newRequestWithParamsList(ctx context.Context, method string, payload interface{}, urlPath []string, params []*QueryParam) (*http.Request, error) {
	var body io.Reader
	URL := strings.Join(append([]string{c.HostURL}, urlPath...), "/")

	paramValues := make([]string, 0)
	for _, p := range params {
		paramValues = append(paramValues, fmt.Sprintf("%s=%v", p.Key, p.Val))
	}
	paramsStr := strings.Join(paramValues, "&")

	if paramsStr != "" {
		URL = fmt.Sprintf("%s?%s", URL, paramsStr)
	}

	if c.debug {
		log.Printf("\n> %s: %s\n", method, URL)
	}

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		if c.debug {
			log.Printf(">>> Payload: %s\n", string(data))
		}

		body = bytes.NewBuffer(data)
	}

	return http.NewRequestWithContext(ctx, method, URL, body)
}

func (c *Client) updateRequest(ctx context.Context, method string, payload interface{}, urlPath ...string) ([]byte, int, error) {
	req, err := c.newRequest(ctx, method, payload, urlPath...)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

func (c *Client) updateURLEncodedRequest(ctx context.Context, method string, params map[string]string, urlPath ...string) ([]byte, int, error) {
	req, err := c.newURLEncodedRequest(ctx, method, params, urlPath...)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

func (c *Client) post(ctx context.Context, payload interface{}, urlPath ...string) error {
	body, statusCode, err := c.updateRequest(ctx, "POST", payload, urlPath...)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated && statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}

	return parseError(body)
}

func (c *Client) postWithResult(ctx context.Context, payload interface{}, urlPath ...string) ([]byte, error) {
	body, statusCode, err := c.updateRequest(ctx, "POST", payload, urlPath...)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusCreated && statusCode != http.StatusOK {
		return nil, ErrStatus(statusCode, body)
	}

	return body, parseError(body)
}

type replyMsg struct {
	Message string
	Error   string
}

func parseError(body []byte) error {
	msgs := bytes.Split(body, []byte("\n"))
	if len(msgs) == 0 {
		return nil
	}

	msg := getLastMessage(msgs)
	if msg == nil || len(msg) == 0 {
		return nil
	}

	var m replyMsg
	err := json.Unmarshal(msg, &m)
	if err != nil {
		// failed to unmarshal, probably there is different payload
		return nil
	}

	if m.Error == "" {
		return nil
	}

	return errors.New(m.Error)
}

func getLastMessage(msgs [][]byte) []byte {
	for i := len(msgs) - 1; i >= 0; i-- {
		if len(msgs[i]) > 0 {
			return msgs[i]
		}
	}

	return nil
}

func (c *Client) postURLEncoded(ctx context.Context, params map[string]string, urlPath ...string) error {
	body, statusCode, err := c.updateURLEncodedRequest(ctx, "POST", params, urlPath...)
	if err != nil {
		return err
	}

	if c.debug {
		log.Println("### Deploy app RESP:", string(body))
	}

	if statusCode != http.StatusAccepted && statusCode != http.StatusCreated && statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}

	if bytes.Contains(body, []byte("There are vulnerabilities!")) {
		return errors.New("found vulnerabilities")
	}

	if bytes.Contains(bytes.ToLower(body), []byte(`"error"`)) {
		return fmt.Errorf("app deploy failed, body: %s", body)
	}

	return nil
}

func (c *Client) put(ctx context.Context, payload interface{}, urlPath ...string) error {
	body, statusCode, err := c.updateRequest(ctx, "PUT", payload, urlPath...)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}
	return nil
}

func (c *Client) delete(ctx context.Context, urlPath ...string) error {
	req, err := c.newRequest(ctx, "DELETE", nil, urlPath...)
	if err != nil {
		return err
	}

	body, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}
	return nil
}

// QueryParam - query parameter
type QueryParam struct {
	Key string
	Val interface{}
}

func (c *Client) deleteWithParams(ctx context.Context, params []*QueryParam, urlPath ...string) error {
	req, err := c.newRequestWithParamsList(ctx, "DELETE", nil, urlPath, params)
	if err != nil {
		return err
	}

	body, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}
	return nil
}

func (c *Client) deleteWithPayload(ctx context.Context, payload interface{}, params map[string]string, urlPath ...string) error {
	req, err := c.newRequestWithParams(ctx, "DELETE", payload, urlPath, params)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	body, statusCode, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return ErrStatus(statusCode, body)
	}
	return nil
}

// ErrStatus - returns error with status and message
func ErrStatus(statusCode int, body []byte) error {
	return fmt.Errorf("status: %d, body: %s", statusCode, body)
}

func (c *Client) testAuthentication() error {
	_, err := c.ListPlans(context.TODO())
	return err
}
