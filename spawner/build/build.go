package build

import (
  "flag"
  "fmt"
  "os"
)

var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion string
)

var Info = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildUser": BuildUser,
	"buildDate": BuildDate,
	"goVersion": GoVersion,
}


var versionFlag = flag.Bool("version", false, "print spawner version and exit")

func LoadInfo() {

  flag.Parse()

  if *versionFlag {
		fmt.Printf("Version %s (%s) - %s - %s - %s\n", Info["version"], Info["revision"], Info["branch"], Info["buildUser"], Info["buildDate"])
		os.Exit(0)
	}
}

