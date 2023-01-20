package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/log"
	"github.com/conjurinc/conjur-preflight/pkg/report"
	"github.com/conjurinc/conjur-preflight/pkg/version"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	var debug bool

	rootCmd := &cobra.Command{
		Use:   "conjur-preflight",
		Short: "Qualification CLI for common Conjur Enterprise self-hosted issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				log.EnableDebugMode()
			}

			report := report.NewDefaultReport(debug)

			log.Debug("Running report...")
			result := report.Run()

			// Determine whether we want to use rich text or plain text based on
			// whether we're outputting directly to a terminal or not
			o, _ := os.Stdout.Stat()
			var formatStrategy framework.FormatStrategy
			if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
				log.Debug("Using rich text report formatting")
				formatStrategy = &framework.RichTextFormatStrategy{}
			} else { //It is not the terminal
				log.Debug("Using plain text report formatting")
				formatStrategy = &framework.PlainTextFormatStrategy{}
			}

			reportText, err := result.ToText(formatStrategy)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), reportText)
			log.Debug("Preflight finished!")
			return nil
		},
		Version: version.FullVersionName,
	}

	rootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"",
		false,
		"debug logging output",
	)

	// TODO: Add JSON output option
	// TODO: Verbose logging control
	// TODO: Ability to adjust requirement criteria (PASS, WARN, FAIL checks)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(stdout, stderr io.Writer) {
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	err := rootCmd.Execute()

	if err != nil {
		log.Error("ERROR: %s\n", err)
		os.Exit(1)
	}
}

var rootCmd = newRootCommand()
