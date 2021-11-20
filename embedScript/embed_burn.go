package embedScript

import (
	"embed"
	_ "embed"
)
//go:embed burntest.sh

var Fs embed.FS


