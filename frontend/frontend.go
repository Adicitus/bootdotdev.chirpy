/*
frontend.go

This package is just a container for static assets that constitutes the Chirpy
frontend to allow embedding.
*/
package frontend

import "embed"

//go:embed *
var FS embed.FS
