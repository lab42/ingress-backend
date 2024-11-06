package static

import "embed"

//go:embed *.html *.ttf *.js favicon.ico
var FS embed.FS
