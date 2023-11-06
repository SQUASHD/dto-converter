# go-dto - A Simple CLI for Record-Based DTO Conversion

## Overview

`go-dto` is a straightforward CLI tool designed to convert C# record-based DTOs to TypeScript types. It was
created out of personal need to streamline conversions in my projects.

Inspired by the functionality of Waseem's [Gopen project](https://github.com/wipdev-tech/gopen), `go-dto` offers a
similar ease of managing project configurations for type conversions. It works best with DTOs that are structured as
records, such as:

```csharp
public record ExampleDto(
   int PropertyOne,
   string PropertyTwo
   );
```

`go-dto` is intended to do one thing well but comes with no guarantees. It's perfect for my use case, and if it fits
yours, that's a win. If not, you might need to extend it to meet your specific requirements.

## Installation

Just run:

```bash
go install ./... "github.com/SQUASHD/dto-converter"
```
Test on the example project:

```bash
go-dto add example csharp example/input example/output
go-dto go
```

## Usage

Here's how to use the tool:

```text
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
```

Just clone, configure, and run.

### Rooms for Improvement

* Autocomplete for commands and arguments.
* Support for more languages.
* Add flags or config for getting feedback on the conversion process.
