package datamover_initializer

import (
	"bitbucket.org/digi-sense/gg-core"
	_ "embed"
	"fmt"
)

//go:embed default_settings.json
var DefaultSettings string

//go:embed default_job.json
var DefaultSettingsJob string

//go:embed default_services_settings.json
var DefaultSettingsService string

//go:embed default_globals.json
var DefaultGlobals string

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

	// settings net
	filename = gg.Paths.Concat(wpRoot, fmt.Sprintf("services.%s.json", mode))
	// ensure settings exists
	_ = gg.Paths.Mkdir(filename)
	if b, _ := gg.Paths.Exists(filename); !b {
		_, _ = gg.IO.WriteTextToFile(DefaultSettingsService, filename)
	}

	// connections
	filename = gg.Paths.Concat(wpRoot, fmt.Sprintf("globals.%s.json", mode))
	if b, _ := gg.Paths.Exists(filename); !b {
		_, _ = gg.IO.WriteTextToFile(DefaultGlobals, filename)
	}
	return
}

func CreateJobSettings(dir string) (err error) {
	filename := gg.Paths.Concat(dir, "job.json")
	if b, _ := gg.Paths.Exists(filename); !b {
		_, err = gg.IO.WriteTextToFile(DefaultSettingsJob, filename)
	}
	return
}
