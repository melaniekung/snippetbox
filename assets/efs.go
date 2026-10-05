package assets

import "embed"

/* comment directive to store files in an embedded filesystem */
//go:embed "html" "static"
var Files embed.FS
