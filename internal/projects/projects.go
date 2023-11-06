package projects

import (
	"encoding/json"
	"fmt"
	"github.com/SQUASHD/dto-converter/internal/factory"
	"github.com/SQUASHD/dto-converter/internal/structs"
	"github.com/SQUASHD/dto-converter/internal/utils"
	"os"
)

type CommandContext struct {
	Config     *structs.Config
	ConfigPath string
}

func (ctx *CommandContext) Write() error {
	jsonData, err := json.MarshalIndent(ctx.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}
	if err = os.WriteFile(ctx.ConfigPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}
	return nil
}

func (ctx *CommandContext) List() (fmtProjects []string) {
	var width int
	for _, project := range ctx.Config.Projects {
		if len(project.Name) > width {
			width = len(project.Name)
		}
	}
	for _, project := range ctx.Config.Projects {
		fmtProject := fmt.Sprintf("%*s: %s, in=%s, out=%s", width, project.Name, project.Language, project.InPath, project.OutPath)
		fmtProjects = append(fmtProjects, fmtProject)
	}
	return
}

func (ctx *CommandContext) Add(name, language, inPath, outPath string) error {

	reserved := []string{
		"csharp", "go", "java", "python", "typescript", "javascript", "cpp",
		"swift", "kotlin", "rust", "scala", "php", "objective-c", "fsharp",
		"dart", "ada", "haskell",
		"run", "list", "add", "remove", "help", "init", "projects", "set",
		"r", "l", "a", "rm", "h", "i", "p", "s",
	}
	for _, r := range reserved {
		if r == name {
			return fmt.Errorf("error: `%v` is reserved and can't be used as a project's name", name)
		}
	}

	for i, project := range ctx.Config.Projects {
		if project.Name == name {
			ctx.Config.Projects[i] = structs.Project{Name: name, Language: language, InPath: inPath, OutPath: outPath}
			return nil
		}
	}

	ctx.Config.Projects = append(ctx.Config.Projects, structs.Project{Name: name, Language: language, InPath: inPath, OutPath: outPath})
	return nil
}

func (ctx *CommandContext) HandleProjects(args []string) {
	if len(args) == 2 {
		ctx.ListProjects()
	} else if len(args) == 3 {
		ctx.ShowProjectDetail(args[2])
	} else {
		fmt.Println("Incorrect number of arguments for 'projects'.")
	}
}

func (ctx *CommandContext) ListProjects() {
	for _, project := range ctx.Config.Projects {
		fmt.Println(project.Name)
	}
}

func (ctx *CommandContext) ShowProjectDetail(projectName string) {
	for _, project := range ctx.Config.Projects {
		if projectName == project.Name {
			fmt.Printf("Language: %s, In=%s, Out=%s\n", project.Language, project.InPath, project.OutPath)
			return
		}
	}
	fmt.Println("Project doesn't exist")
}

func (ctx *CommandContext) AddProject(name, language, inputPath, outputPath string) {
	absInput, err := utils.NormalizePath(inputPath)
	utils.ErrFatal(err)

	absOutput, err := utils.NormalizePath(outputPath)
	utils.ErrFatal(err)

	err = ctx.Add(name, language, absInput, absOutput)
	utils.ErrFatal(err)

	err = ctx.Write()
	utils.ErrFatal(err)
}

func (ctx *CommandContext) HandleConvertDto(projectName string) {

	projectData, found := ctx.FindProjectByName(projectName)
	if !found {
		fmt.Println("Project not found.")
		return
	}

	converter, err := factory.ConverterFactory(projectData)
	utils.ErrFatal(err)
	if err = converter.ConvertDirectory(projectData.InPath); err != nil {
		fmt.Println(err)
	}
}

func (ctx *CommandContext) FindProjectByName(name string) (*structs.Project, bool) {
	for i, project := range ctx.Config.Projects {
		if name == project.Name {
			return &ctx.Config.Projects[i], true
		}
	}
	return nil, false
}

func (ctx *CommandContext) HandleAdd(args []string) {

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments for 'add'.")
		return
	}
	name := args[2]
	language := args[3]
	inputPath := args[4]
	outputPath := args[5]

	supported := false
	for _, supportedLanguage := range factory.SupportedLanguages {
		if language == supportedLanguage {
			supported = true
			break
		}
	}
	if !supported {
		fmt.Printf("Language `%s` is not supported.\n", language)
		return
	}

	absInput, err := utils.NormalizePath(inputPath)
	utils.ErrFatal(err)

	absOutput, err := utils.NormalizePath(outputPath)
	utils.ErrFatal(err)

	err = ctx.Add(name, language, absInput, absOutput)
	utils.ErrFatal(err)

	err = ctx.Write()
	utils.ErrFatal(err)
}

func (ctx *CommandContext) HandleHelp() {
	helpText := `
Usage: go-dto COMMAND [OPTIONS]

Commands:
  init, i             Initialize a new configuration file if it doesn't exist.
  add, a              Add a new project to the configuration.
  help, h             Show detailed information about commands and usage.
  projects, p         Display a list of all configured projects.
  set, s              Assign input or output directories to a specified project.
  remove, r           Remove an existing project from the configuration.
  run, go             Execute the converter process for all or specific projects.

Arguments:
  projectName         Specify the target project for the command.
  in/out              Indicate whether to set the 'input' or 'output' directory.
  language            Define the programming language for a new project.
  inputPath           Path to the source files for conversion.
  outputPath          Destination path for the converted files.

Examples:
  Initialize configuration:
    go-dto init

  List all projects:
    go-dto projects

  Add a new project:
    go-dto add projectName language inputPath outputPath

  Set the input directory to current directory:
    go-dto set projectName in .

  Remove a project:
    go-dto remove projectName

  Run the converter:
    go-dto go projectName   # Optional: Specify a project name to run individually.
`
	fmt.Print(helpText)
}

func (ctx *CommandContext) HandleSetDirectories(projectName, inOrOut, path string) {
	project, exists := ctx.FindProjectByName(projectName)
	if !exists {
		fmt.Println("Project not found.")
		return
	}

	if inOrOut != "in" && inOrOut != "out" {
		fmt.Println("Incorrect arguments for 'set': expected 'in' or 'out'")
		return
	}

	// If the path is ".", then use the current working directory
	if path == "." {
		dir, err := os.Getwd()
		utils.ErrFatal(err)
		path = dir
	} else {
		dir, err := utils.NormalizePath(path)
		utils.ErrFatal(err)
		path = dir
	}
	switch inOrOut {
	case "in":
		fmt.Println(path)
		project.InPath = path
	case "out":
		fmt.Println(path)
		project.OutPath = path
	}

	err := ctx.Write()
	utils.ErrFatal(err)
}

func (ctx *CommandContext) HandleRemove(projectName string) {
	for i, project := range ctx.Config.Projects {
		if project.Name == projectName {
			ctx.Config.Projects = append(ctx.Config.Projects[:i], ctx.Config.Projects[i+1:]...)
			return
		}
	}
	fmt.Println("Project doesn't exist")
}

func (ctx *CommandContext) HandleRun(args []string) {
	if len(args) == 2 {
		for _, project := range ctx.Config.Projects {
			ctx.HandleConvertDto(project.Name)
		}
	}
	if len(args) == 3 {
		ctx.HandleConvertDto(args[2])
	}
}
