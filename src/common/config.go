package common

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	log "third/seelog"
	toml "third/toml"
)

var G gConfig
var C gCommand

type gConfig struct {
	Basic basic `check:"Struct"`
}

var (
	MY_NAME    string = "sample"
	author     string = "Please build by build.sh"
	githash    string = "Please build by build.sh"
	buildstamp string = "Please build by build.sh"
	goversion  string = runtime.Version()
)

type basic struct {
	LogConfigFile    string `check:"StringNotEmpty"`
	DebugBindAddress string `check:"StringNotEmpty"`
}

type gCommand struct {
	ConfigFile   *string
	LogConfFile  *string
	PrintVersion *bool
	Foreground   *bool
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (self *gCommand) parseCommand(args []string) {
	flagset := flag.NewFlagSet(MY_NAME, flag.ExitOnError)
	self.ConfigFile = flagset.String("config", "./etc/"+MY_NAME+".conf", "Path to config file")
	self.LogConfFile = flagset.String("logconf", "", "Path to config file")
	self.PrintVersion = flagset.Bool("version", false, "Print version")
	self.Foreground = flagset.Bool("fg", false, "Start server in foreground")
	flagset.Parse(args[1:])

	if *self.PrintVersion {
		fmt.Fprintf(os.Stdout, "%s\n", MY_NAME)
		fmt.Fprintf(os.Stdout, "Author:           %s\n", author)
		fmt.Fprintf(os.Stdout, "Git Commit Hash:  %s\n", githash)
		fmt.Fprintf(os.Stdout, "Build Time:       %s\n", buildstamp)
		fmt.Fprintf(os.Stdout, "Go Version:       %s\n", goversion)
		os.Exit(0)
	}
}

func (self gConfig) check() error {
	return configCheckStruct(MY_NAME, self)
}

func (self gConfig) String() string {
	return configStringStruct(MY_NAME, self)
}

func ParseCommandAndFile() error {
	C.parseCommand(os.Args)

	_, err := toml.DecodeFile(*C.ConfigFile, &G)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config file parse fail, err=[%s] file=[%s]\n", err.Error(), *C.ConfigFile)
		os.Exit(-1)
	}

	err = G.check()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config check fail, err=[%s]\n", err.Error())
		os.Exit(-1)
	}

	if *C.LogConfFile != "" {
		G.Basic.LogConfigFile = *C.LogConfFile
	}
	logger, err := log.LoggerFromConfigAsFile(G.Basic.LogConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seelog LoggerFromConfigAsFile fail, err=[%s]\n", err.Error())
		os.Exit(-1)
	}
	log.ReplaceLogger(logger)

	log.Infof("parse config succ, %s", G.String())
	return err
}
