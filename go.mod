module github.com/FerretDB/gh

go 1.24

toolchain go1.24.2

// We use @master for go-github for now to include issue types (https://github.com/google/go-github/pull/3525).
// We will switch to v71 once it is released.

require (
	github.com/google/go-github/v70 v70.0.1-0.20250402125210-3a3f51bc7c5d
	golang.org/x/oauth2 v0.28.0
)

require github.com/google/go-querystring v1.1.0 // indirect
