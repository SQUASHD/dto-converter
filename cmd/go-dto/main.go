package main

import (
	"dto-converter/internal/config"
	"dto-converter/internal/projects"
	"fmt"
	"os"
)

var (
	configDir  = os.Getenv("HOME") + "/.config/go-dto"
	configPath = configDir + "/godto.json"
)

func handleInit() {
	err := config.Init(configDir, configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}
}

func main() {
	cfg, err := config.Read(configPath)
	if err != nil {
		handleInit()
		return
	}
	ctx := &projects.CommandContext{Config: &cfg, ConfigPath: configPath}

	if len(os.Args) < 2 {
		ctx.HandleHelp()
		return
	}

	switch os.Args[1] {
	case "i", "init":
		handleInit()

	case "h", "help":
		ctx.HandleHelp()

	case "a", "add":
		ctx.HandleAdd(os.Args)

	case "p", "projects":
		ctx.HandleProjects(os.Args)

	case "s", "set":
		ctx.HandleSetDirectories(os.Args[2], os.Args[3], os.Args[4])

	case "r", "remove":
		ctx.HandleRemove(os.Args[2])

	case "go", "run":
		ctx.HandleRun(os.Args)

	default:
		ctx.HandleConvertDto(os.Args[1])
	}
}
