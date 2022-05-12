// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

// Package exe defines QoL functions to simplify and unify creating executables
package exe

import (
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/microsoft/CBL-Mariner/toolkit/tools/pkg/logger"
)

// ToolkitVersion specifies the version of the toolkit and the reported version of all tools in it.
var ToolkitVersion = ""

// InputFlag registers an input flag for k with documentation doc and returns the passed value
func InputFlag(k *kingpin.Application, doc string) *string {
	return k.Flag("input", doc).Required().ExistingFile()
}

// InputStringFlag registers an input flag for k with documentation doc and returns the passed value
func InputStringFlag(k *kingpin.Application, doc string) *string {
	return k.Flag("input", doc).Required().String()
}

// InputDirFlag registers an input flag for k with documentation doc and returns the passed value
func InputDirFlag(k *kingpin.Application, doc string) *string {
	return k.Flag("dir", doc).Required().ExistingDir()
}

// OutputFlag registers an output flag for k with documentation doc and returns the passed value
func OutputFlag(k *kingpin.Application, doc string) *string {
	return k.Flag("output", doc).Required().String()
}

// OutputDirFlag registers an output flag for k with documentation doc and returns the passed value
func OutputDirFlag(k *kingpin.Application, doc string) *string {
	return k.Flag("output-dir", doc).Required().String()
}

// LogFileFlag registers a log file flag for k and returns the passed value
func LogFileFlag(k *kingpin.Application) *string {
	return k.Flag(logger.FileFlag, logger.FileFlagHelp).String()
}

// LogLevelFlag registers a log level flag for k and returns the passed value
func LogLevelFlag(k *kingpin.Application) *string {
	return k.Flag(logger.LevelsFlag, logger.LevelsHelp).PlaceHolder(logger.LevelsPlaceholder).Enum(logger.Levels()...)
}

// PlaceHolderize takes a list of available inputs and returns a corresponding placeholder
func PlaceHolderize(thing []string) string {
	return fmt.Sprintf("(%s)", strings.Join(thing, "|"))
}

// ParseListArgument takes a user provided string list that is space seperated
// and returns a slice of the split and trimmed elements.
func ParseListArgument(input string) (results []string) {
	const delimiter = " "

	trimmedInput := strings.TrimSpace(input)
	if trimmedInput != "" {
		results = strings.Split(trimmedInput, delimiter)
	}
	return
}
