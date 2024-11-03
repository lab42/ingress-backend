package static

import "embed"

//go:embed *.html *.ttf *.js
var FS embed.FS
