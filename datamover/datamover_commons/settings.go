package datamover_commons

import (
	"bitbucket.org/digi-sense/gg-core"
	"fmt"
)

type DataMoverSettings struct {
	PathJobs string           `json:"path_jobs"`
	Postman  *PostmanSettings `json:"postman"`
}

// ---------------------------------------------------------------------------------------------------------------------
//	s e t t i n g s
// ---------------------------------------------------------------------------------------------------------------------

func NewSettings(mode string) (*DataMoverSettings, error) {
	instance := new(DataMoverSettings)
	err := instance.init(mode)

	return instance, err
}

func (instance *DataMoverSettings) String() string {
	return gg.JSON.Stringify(instance)
}

func (instance *DataMoverSettings) init(mode string) error {
	if len(mode) == 0 {
		mode = "production"
	}
	dir := gg.Paths.WorkspacePath("./")
	filename := gg.Paths.Concat(dir, fmt.Sprintf("settings.%s.json", mode))

	// load settings
	return gg.JSON.ReadFromFile(filename, &instance)
}
