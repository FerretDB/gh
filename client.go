package gh

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

// NewRESTClient returns GitHub REST API client for the given token (that may be empty)
// and debug logging function (that may be nil).
func NewRESTClient(token string, debugf Printf) (*github.Client, error) {
	// don't use http.DefaultClient and oauth2.NewClient to avoid data races

	httpTransport := http.DefaultTransport
	if debugf != nil {
		httpTransport = NewTransport(http.DefaultTransport, debugf)
	}

	var httpClient *http.Client
	if token == "" {
		httpClient = &http.Client{
			Transport: httpTransport,
		}
	} else {
		httpClient = &http.Client{
			Transport: &oauth2.Transport{
				Base: httpTransport,
				Source: oauth2.ReuseTokenSource(nil, oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: token},
				)),
			},
		}
	}

	c := github.NewClient(httpClient)
	c.UserAgent = "FerretDB-gh/1.0 (+https://github.com/FerretDB/gh)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query rate limit to check that the client is able to make queries.
	// See https://docs.github.com/en/rest/rate-limit.
	// We can't use https://docs.github.com/en/rest/users/users#get-the-authenticated-user API,
	// because short-lived automatic GITHUB_TOKEN is provided by GitHub Actions App that can't access this API.
	rl, _, err := c.RateLimit.Get(ctx)

	if rl != nil && debugf != nil {
		debugf(
			"Rate limit: %d/%d, resets at: %s.",
			rl.Core.Remaining, rl.Core.Limit, rl.Core.Reset.Format(time.RFC3339),
		)
	}

	return c, err
}
