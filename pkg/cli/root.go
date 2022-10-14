package cli

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/telemetry"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func RunRootSchemaCommand() {
	command := SetupRootSchemaCommand()
	utils.DoOrDie(errors.Wrapf(command.Execute(), "run root schema command"))
}

type RootSchemaFlags struct {
	Verbosity string
}

func SetupRootSchemaCommand() *cobra.Command {
	flags := &RootSchemaFlags{}
	command := &cobra.Command{
		Use:   "scaling",
		Short: "scaling hacking",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return telemetry.SetUpLogger(flags.Verbosity)
		},
	}

	command.PersistentFlags().StringVarP(&flags.Verbosity, "verbosity", "v", "info", "log level; one of [info, debug, trace, warn, error, fatal, panic]")

	command.AddCommand(SetupVersionCommand())
	//command.AddCommand(SetupUploadCommand())
	//command.AddCommand(setupWebServerCommand())
	//command.AddCommand(SetupParserCommand())
	//command.AddCommand(SetupAnalyzeCommand())

	return command
}

var (
	version   = "development"
	gitSHA    = "development"
	buildTime = "development"
)

func SetupVersionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "print out version information",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunVersionCommand()
		},
	}

	return command
}

func RunVersionCommand() {
	jsonString, err := json.MarshalToString(map[string]string{
		"Version":   version,
		"GitSHA":    gitSHA,
		"BuildTime": buildTime,
	})
	utils.DoOrDie(err)
	fmt.Printf("scaling version: \n%s\n", jsonString)
}
