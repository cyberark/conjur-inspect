package cmd

import (
	"fmt"
	"os"

	"github.com/conjurinc/conjur-preflight/pkg/framework"
	"github.com/conjurinc/conjur-preflight/pkg/report"
	"github.com/conjurinc/conjur-preflight/pkg/version"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "conjur-preflight",
		Short: "Qualification CLI for common Conjur Enterprise self-hosted issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			report := report.NewDefaultReport()
			result := report.Run()

			// Determine whether we want to use rich text or plain text based on
			// whether we're outputting directly to a terminal or not
			o, _ := os.Stdout.Stat()
			var formatStrategy framework.FormatStrategy
			if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
				formatStrategy = &framework.RichTextFormatStrategy{}
			} else { //It is not the terminal
				formatStrategy = &framework.PlainTextFormatStrategy{}
			}

			reportText, err := result.ToText(formatStrategy)
			if err != nil {
				return err
			}

			fmt.Println(reportText)

			return nil
		},
		Version: version.FullVersionName,
	}

	// TODO: Add JSON output option
	// TODO: Verbose logging control
	// TODO: Ability to adjust requirement criteria (PASS, WARN, FAIL checks)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = newRootCommand()
