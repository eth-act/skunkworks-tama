//go:build tamago

package testdata

import (
	"embed"
)

//go:embed *
var FS embed.FS
