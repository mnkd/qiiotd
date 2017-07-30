package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
)

var (
	Version  string
	Revision string
)

var app = App{}

func init() {
	var version bool
	var path string
	var yearsAgo int
	var days int
	flag.BoolVar(&version, "v", false, "prints current qiotd version")
	flag.StringVar(&path, "c", "", "/path/to/config.json (default: $HOME/.config/qiotd/config.json)")
	flag.IntVar(&yearsAgo, "ago", 1, "Years ago (default: 1)")
	flag.IntVar(&days, "days", 1, "Days (default: 1)")
	flag.Parse()

	if version {
		fmt.Fprintln(os.Stdout, "Version:", Version)
		fmt.Fprintln(os.Stdout, "Revision:", Revision)
		os.Exit(ExitCodeOK)
	}

	config, err := NewConfig(path)
	if err != nil {
		os.Exit(ExitCodeError)
	}

	app = NewApp(config, yearsAgo, days)
}

func main() {
	os.Exit(app.Run())
}
