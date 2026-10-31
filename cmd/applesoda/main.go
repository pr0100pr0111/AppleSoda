package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
	"github.com/pr0100pr0111/AppleSoda/internal/ui"
)

var version = "0.1.0"

func main() {
	demoFlag := flag.Bool("demo", false, "run in demo mode with synthetic data")
	versionFlag := flag.Bool("version", false, "show version")
	configFlag := flag.String("config", "", "path to config file")
	flag.StringVar(configFlag, "c", "", "path to config file (shorthand)")
	flag.BoolVar(versionFlag, "v", false, "show version (shorthand)")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("applesoda %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	var p player.Player
	if *demoFlag {
		p = player.NewFakePlayer()
	} else {
		p = player.NewMusicAppPlayer()
	}

	model := ui.NewModel(p, cfg)
	prog := tea.NewProgram(model)
	if _, err := prog.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "applesoda error: %v\n", err)
		os.Exit(1)
	}
}
