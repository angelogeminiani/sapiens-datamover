package main

import (
	"bitbucket.org/digi-sense/gg-core"
	_ "bitbucket.org/digi-sense/gg-core-x"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"flag"
	"log"
	"os"
)

func main() {
	// PANIC RECOVERY
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			message := gg.Strings.Format("[panic] APPLICATION '%s' ERROR: %s", datamover_commons.AppName, r)
			log.Fatalf(message)
		}
	}()

	//-- command flags --//
	// run
	cmdRun := flag.NewFlagSet("run", flag.ExitOnError)
	dirWork := cmdRun.String("dir_work", gg.Paths.Absolute("./_workspace"), "Set a particular folder as main workspace")
	mode := cmdRun.String("m", datamover_commons.ModeProduction, "Mode allowed: 'debug' or 'production'")
	quit := cmdRun.String("s", "stop", "Quit Command: Write a command (ex: 'stop') to enable stop mode")

	// parse
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "run":
			_ = cmdRun.Parse(os.Args[2:])
		default:
			panic("Command not supported: " + cmd)
		}
	} else {
		panic("Missing command. i.e. 'run'")
	}

	initialize(dirWork, mode)

	app, err := datamover.LaunchApplication(*mode, *quit)
	if nil == err {

		err = app.Start()
		if nil != err {
			log.Panicf("Error starting %s: %s", datamover_commons.AppName, err.Error())
		} else {
			// app is running
			app.Join()
		}

	} else {
		log.Panicf("Error launching %s: %s", datamover_commons.AppName, err.Error())
	}

}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func initialize(dirWork *string, mode *string) {
	gg.Paths.GetWorkspace(datamover_commons.WpDirWork).SetPath(*dirWork)
}
