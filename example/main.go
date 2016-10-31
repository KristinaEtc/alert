package main

import (
	"github.com/KristinaEtc/alert"
	_ "github.com/KristinaEtc/slflog"

	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("alert.go")

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	Version    string
)

// ConfFile is a file with all program options
type ConfFile struct {
	FileName string
}

var globalOpt = ConfFile{
	FileName: "7.aiff",
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "music options")

	log.Infof("%s", globalOpt.FileName)
	alert.PlayMusic(globalOpt.FileName)

	log.Error("----------------------------------------------")

	log.Infof("BuildDate=%s\n", BuildDate)
	log.Infof("GitCommit=%s\n", GitCommit)
	log.Infof("GitBranch=%s\n", GitBranch)
	log.Infof("GitState=%s\n", GitState)
	log.Infof("GitSummary=%s\n", GitSummary)
	log.Infof("VERSION=%s\n", Version)

	log.Info("Starting working...")

}
