package gh

import (
	"context"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

// NewRESTClient returns GitHub REST API client for the given token (that may be empty)
// and debug logging function (that may be nil).
func NewRESTClient(token string, debugf Printf) (*github.Client, error) {
	var src oauth2.TokenSource
	if token != "" {
		src = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
	}

	httpClient := oauth2.NewClient(context.Background(), src)
	if debugf != nil {
		httpClient.Transport = NewTransport(httpClient.Transport, debugf)
	}

	c := github.NewClient(httpClient)
	c.UserAgent = "FerretDB-gh/1.0 (+https://github.com/FerretDB/gh)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query rate limit to check that the client is able to make queries.
	// See https://docs.github.com/en/rest/rate-limit.
	// We can't use https://docs.github.com/en/rest/users/users#get-the-authenticated-user API,
	// because short-lived automatic GITHUB_TOKEN is provided by GitHub Actions App that can't access this API
	// (and doesn't have authenticated user).
	_, _, err := c.RateLimits(ctx)
	return c, err
}
