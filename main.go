// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

// Command tfproviderdocs validates Terraform provider documentation.
//
// Deprecated: tfproviderdocs is no longer maintained. All functionality has
// been superseded by swissshepherd. Please migrate to:
// https://github.com/YakDriver/swissshepherd
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/YakDriver/tfproviderdocs/command"
	"github.com/YakDriver/tfproviderdocs/version"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

const (
	Name = `tfproviderdocs`
)

func main() {
	printDeprecationNotice(os.Stderr)

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		InfoColor:  cli.UiColorGreen,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      colorable.NewColorableStdout(),
			ErrorWriter: colorable.NewColorableStderr(),
		},
	}

	c := &cli.CLI{
		Name:     Name,
		Version:  version.GetVersion().FullVersionNumber(true),
		Args:     os.Args[1:],
		Commands: command.Commands(ui),
	}

	exitStatus, err := c.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(exitStatus)
}

// printDeprecationNotice writes a deprecation banner to the given writer
// (typically os.Stderr) so it does not pollute stdout consumed by CI parsers.
func printDeprecationNotice(w io.Writer) {
	fmt.Fprintln(w, "============================================================")
	fmt.Fprintln(w, "DEPRECATED: tfproviderdocs is no longer maintained.")
	fmt.Fprintln(w, "All functionality has been superseded by swissshepherd:")
	fmt.Fprintln(w, "    https://github.com/YakDriver/swissshepherd")
	fmt.Fprintln(w, "Please migrate. This tool will receive no further releases.")
	fmt.Fprintln(w, "============================================================")
}
