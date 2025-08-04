package web

import (
	"embed"
)

//go:embed css
var Css embed.FS

//go:embed js
var Js embed.FS

//go:embed images
var Images embed.FS

//go:embed public
var Public embed.FS

//go:embed doc
var Doc embed.FS
