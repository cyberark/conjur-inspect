package cmd

import (
	"io"
	"os"

	"github.com/cyberark/conjur-inspect/pkg/formatting"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/report"
	"github.com/cyberark/conjur-inspect/pkg/version"
	"github.com/spf13/cobra"
)

func init() {

}

func newRootCommand() *cobra.Command {
	var debug bool
	var jsonOutput bool

	rootCmd := &cobra.Command{
		Use:   "conjur-inspect",
		Short: "Qualification CLI for common Conjur Enterprise self-hosted issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				log.EnableDebugMode()
			}

			report := report.NewDefaultReport(debug)

			log.Debug("Running report...")
			result := report.Run()

			// Determine which output format we'll use
			var writer formatting.Writer
			switch {
			case jsonOutput:
				log.Debug("Using JSON report formatting")
				writer = &formatting.JSON{}
			case isTerminal(cmd.OutOrStdout()):
				log.Debug("Using rich text report formatting")
				writer = &formatting.Text{
					FormatStrategy: &formatting.RichANSIFormatStrategy{},
				}
			default:
				log.Debug("Using plain text report formatting")
				writer = &formatting.Text{
					FormatStrategy: &formatting.PlainFormatStrategy{},
				}
			}

			// Write the report result
			err := writer.Write(cmd.OutOrStdout(), &result)
			if err != nil {
				return err
			}

			log.Debug("Inspection finished!")
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

	// Create json flag for the conjur-inspect command to output a report.
	// Usage: conjur-inspect --json or -j
	rootCmd.PersistentFlags().BoolVarP(
		&jsonOutput,
		"json",
		"j",
		false,
		"Output report in JSON",
	)

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

func isTerminal(writer io.Writer) bool {
	// Test if the writer is for a file. If not, we know it isn't a terminal
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}

	o, err := file.Stat()

	// If there's an error stat-ing the file, then assume it's not a terminal
	if err != nil {
		return false
	}

	// Check to see whether this is a device or a regular file
	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}

var rootCmd = newRootCommand()
