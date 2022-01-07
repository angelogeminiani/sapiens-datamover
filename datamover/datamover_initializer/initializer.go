package datamover_initializer

import (
	"bitbucket.org/digi-sense/gg-core"
	_ "embed"
	"fmt"
)

//go:embed default_settings.json
var DefaultSettings string

func Initialize(mode string) (err error) {
	wpRoot := gg.Paths.WorkspacePath("./")
	tmpRoot := gg.Paths.Concat(wpRoot, "temp")

	gg.Paths.SetTempRoot(tmpRoot)
	gg.IO.RemoveAllSilent(tmpRoot)

	// settings
	filename := gg.Paths.Concat(wpRoot, fmt.Sprintf("settings.%s.json", mode))
	// ensure settings exists
	_ = gg.Paths.Mkdir(filename)
	if b, _ := gg.Paths.Exists(filename); !b {
		_, _ = gg.IO.WriteTextToFile(DefaultSettings, filename)
	}

	return
}
