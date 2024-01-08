package graphql

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/graphql-utils/client"
	"github.com/hgiasac/hasura-router/go/types"
)

const (
	HasuraClientName = "hasura-client-name"
)

var (
	errPromoteAdminDenied = errors.New("cannot promote to admin")
)

type options struct {
	timeout     time.Duration
	clientName  string
	adminSecret string
	debug       bool
}

var defaultOptions = options{
	timeout: 30 * time.Second,
}

type Option func(*options)

// WithTimeout set timeout option to hasura client
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WithClientName set timeout option to hasura client
func WithClientName(clientName string) Option {
	return func(opts *options) {
		opts.clientName = clientName
	}
}

// WithAdminSecret set admin secret option to hasura client
func WithAdminSecret(adminSecret string) Option {
	return func(opts *options) {
		opts.adminSecret = adminSecret
	}
}

// WithDebug enables debug option to hasura client
func WithDebug(value bool) Option {
	return func(opts *options) {
		opts.debug = value
	}
}

// HasuraClientConfig input config for Client
type HasuraClientConfig struct {
	BaseURL     string            `envconfig:"BASE_URL"`
	URL         string            `envconfig:"URL"`
	AdminSecret string            `envconfig:"ADMIN_SECRET"`
	Headers     map[string]string `envconfig:"HEADERS"`
	Timeout     time.Duration     `envconfig:"TIMEOUT" default:"60s"`
	Debug       bool              `envconfig:"DEBUG" default:"false"`
}

// HasuraClient represents a graphql client with Hasura credential
type HasuraClient struct {
	client.Client
	adminSecret      string
	clientName       string
	sessionVariables types.SessionVariables
}

// NewHasuraClient creates a new GraphQL client for Hasura with the HTTP transport
// that uses a RoundTripper to set the headers on every request.
// Headers come in two forms:
// * From the actor information and admin secret provided in ActorAwareClient initialization
// * Previously set in the context object
func NewHasuraClient(endpoint string, options ...Option) *HasuraClient {
	opts := defaultOptions
	for _, apply := range options {
		apply(&opts)
	}

	sessionVariables := types.SessionVariables{}
	if opts.adminSecret != "" {
		sessionVariables.Set(types.XHasuraAdminSecret, opts.adminSecret)
	}
	if opts.clientName != "" {
		sessionVariables.Set(HasuraClientName, opts.clientName)
	}

	return &HasuraClient{
		Client:           client.NewClient(endpoint, buildHttpClient(opts.timeout)).WithDebug(opts.debug),
		adminSecret:      opts.adminSecret,
		clientName:       opts.clientName,
		sessionVariables: sessionVariables,
	}
}

// NewAdminClient creates a new Hasura GraphQL client with admin role
func NewAdminClient(endpoint string, adminSecret string, options ...Option) *HasuraClient {
	return NewHasuraClient(endpoint, append(options, WithAdminSecret(adminSecret))...)
}

// NewAdminClient creates a new Hasura GraphQL client with admin role
func NewHasuraClientFromConfig(config HasuraClientConfig) *HasuraClient {
	endpoint := config.URL
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s/v1/graphql", config.BaseURL)
	}

	sessionVariables := types.SessionVariables{}
	if config.AdminSecret != "" {
		sessionVariables.Set(types.XHasuraAdminSecret, config.AdminSecret)
	}

	for k, v := range config.Headers {
		sessionVariables[k] = v
	}

	return &HasuraClient{
		Client:           client.NewClient(endpoint, buildHttpClient(config.Timeout)).WithDebug(config.Debug),
		adminSecret:      config.AdminSecret,
		clientName:       sessionVariables.Get(HasuraClientName),
		sessionVariables: sessionVariables,
	}
}

// ToSessionVariables create session variables from options
func (c HasuraClient) getDefaultSessionVariables() types.SessionVariables {
	sessionVariables := types.SessionVariables{}
	if c.adminSecret != "" {
		sessionVariables.Set(types.XHasuraAdminSecret, c.adminSecret)
	}
	if c.clientName != "" {
		sessionVariables.Set(HasuraClientName, c.clientName)
	}

	return sessionVariables
}

func (c *HasuraClient) Query(ctx context.Context, q any, variables map[string]any, options ...graphql.Option) error {
	ctx = setHeaders(ctx, c.sessionVariables.ToStringMap())
	return c.Client.Query(ctx, q, variables, options...)
}

func (c *HasuraClient) QueryRaw(ctx context.Context, q any, variables map[string]any, options ...graphql.Option) ([]byte, error) {
	ctx = setHeaders(ctx, c.sessionVariables.ToStringMap())
	return c.Client.QueryRaw(ctx, q, variables, options...)
}

func (c *HasuraClient) Mutate(ctx context.Context, m any, variables map[string]any, options ...graphql.Option) error {
	ctx = setHeaders(ctx, c.sessionVariables.ToStringMap())
	return c.Client.Mutate(ctx, m, variables, options...)
}

func (c *HasuraClient) MutateRaw(ctx context.Context, m any, variables map[string]any, options ...graphql.Option) ([]byte, error) {
	ctx = setHeaders(ctx, c.sessionVariables.ToStringMap())
	return c.Client.MutateRaw(ctx, m, variables, options...)
}

func (c *HasuraClient) Exec(ctx context.Context, query string, m any, variables map[string]any, options ...graphql.Option) error {
	ctx = setHeaders(ctx, c.sessionVariables.ToStringMap())
	return c.Client.Exec(ctx, query, m, variables, options...)
}

func (c *HasuraClient) ExecRaw(ctx context.Context, query string, variables map[string]any, options ...graphql.Option) ([]byte, error) {
	ctx = setHeaders(ctx, c.sessionVariables.ToStringMap())
	return c.Client.ExecRaw(ctx, query, variables, options...)
}

// AsRole allows the client to act on behalf of a new role
func (c *HasuraClient) As(variables map[string]string) (*HasuraClient, error) {
	sessionVariables := c.getDefaultSessionVariables()

	for k, v := range variables {
		sessionVariables[k] = v
	}

	return &HasuraClient{
		Client:           c.Client,
		adminSecret:      c.adminSecret,
		clientName:       c.clientName,
		sessionVariables: sessionVariables,
	}, nil
}

// AsRole allows the client to act on behalf of a new role
func (c *HasuraClient) AsRole(role string, userId string) (*HasuraClient, error) {
	if c.adminSecret == "" {
		return nil, fmt.Errorf("cannot promote to role <%s>", role)
	}

	sessionVariables := c.sessionVariables.FilterKey(types.XHasuraRole, XHasuraUserID)
	sessionVariables[types.XHasuraRole] = role
	if userId != "" {
		sessionVariables[XHasuraUserID] = userId
	}

	return &HasuraClient{
		Client:           c.Client,
		adminSecret:      c.adminSecret,
		clientName:       c.clientName,
		sessionVariables: sessionVariables,
	}, nil
}

// AsAdmin allows the client to act on behalf of an admin
func (c *HasuraClient) AsAdmin() (*HasuraClient, error) {
	if c.adminSecret == "" {
		return nil, errPromoteAdminDenied
	}
	sessionVariables := c.sessionVariables.FilterKey(types.XHasuraRole, XHasuraUserID)

	return &HasuraClient{
		Client:           c.Client,
		adminSecret:      c.adminSecret,
		clientName:       c.clientName,
		sessionVariables: sessionVariables,
	}, nil
}

// ForceAdmin allows the client to act on behalf of an admin, this function panics if the client cannot
// be promoted to an Admin client. Prefer AsAdmin instead.
func (c *HasuraClient) ForceAdmin() *HasuraClient {
	admin, err := c.AsAdmin()
	if err != nil {
		panic(err)
	}
	return admin
}

// AsAnonymous allows the client to act on behalf of an anonymous user
func (c *HasuraClient) AsAnonymous() (*HasuraClient, error) {
	newSession := types.SessionVariables{}
	if c.clientName != "" {
		newSession[HasuraClientName] = c.clientName
	}

	return &HasuraClient{
		Client:           c.Client,
		adminSecret:      c.adminSecret,
		clientName:       c.clientName,
		sessionVariables: newSession,
	}, nil
}

type headerRoundTripper struct {
	setHeaders func(req *http.Request)
	rt         http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.setHeaders(req)
	return h.rt.RoundTrip(req)
}

func buildHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: headerRoundTripper{
			setHeaders: func(req *http.Request) {
				// set headers in the context
				for hn, hv := range getHeadersFromContext(req.Context()) {
					req.Header.Set(hn, hv)
				}
			},
			rt: http.DefaultTransport,
		},
		Timeout: timeout,
	}
}
