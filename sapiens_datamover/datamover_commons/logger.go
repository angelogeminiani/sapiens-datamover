package datamover_commons

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"fmt"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Logger struct {
	mode   string
	root   string
	logger *gg_log.Logger
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewLogger(mode string, logger interface{}) *Logger {
	instance := new(Logger)
	instance.mode = mode
	instance.root = gg.Paths.WorkspacePath("logging")

	// reset file
	_ = gg.IO.RemoveAll(instance.root)
	_ = gg.Paths.Mkdir(instance.root + gg_utils.OS_PATH_SEPARATOR)

	if instance.logger, _ = logger.(*gg_log.Logger); nil == instance.logger {
		instance.logger = gg_log.NewLogger()
		instance.logger.SetFileName(gg.Paths.Concat(instance.root, "logging.log"))
	}

	if mode == ModeDebug {
		instance.logger.SetLevel(gg_log.LEVEL_DEBUG)
	} else {
		instance.logger.SetLevel(gg_log.LEVEL_INFO)
	}
	instance.logger.SetOutput(gg_log.OUTPUT_FILE)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Logger) Logger() *gg_log.Logger {
	return instance.logger
}

func (instance *Logger) Close() {
	instance.logger.Close()
}

func (instance *Logger) SetLevel(level string) {
	instance.logger.SetLevelName(level)
}

func (instance *Logger) GetLevel() int {
	return instance.logger.GetLevel()
}

func (instance *Logger) Debug(args ...interface{}) {
	// file logging
	instance.logger.Debug(args...)

	if instance.mode == ModeDebug {
		// console logging
		fmt.Println(args...)
	}
}

func (instance *Logger) Info(args ...interface{}) {
	// file logging
	instance.logger.Info(args...)

	if instance.mode == ModeDebug {
		// console logging
		fmt.Println(args...)
	}
}

func (instance *Logger) Error(args ...interface{}) {
	// file logging
	instance.logger.Error(args...)

	if instance.mode == ModeDebug {
		// console logging
		fmt.Println(args...)
	}
}

func (instance *Logger) Warn(args ...interface{}) {
	// file logging
	instance.logger.Warn(args...)

	if instance.mode == ModeDebug {
		// console logging
		fmt.Println(args...)
	}
}
