//go:build ui

package ui

import "embed"

//go:embed build build/_app build/_app/pages/*.js build/_app/assets/pages/*.css
var Build embed.FS
