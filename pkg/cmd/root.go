package cmd

import (
	"io"
	"os"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/formatting"
	"github.com/cyberark/conjur-inspect/pkg/log"
	"github.com/cyberark/conjur-inspect/pkg/version"
	"github.com/spf13/cobra"
)

var defaultReportConstructor = NewDefaultReport

func newRootCommand() *cobra.Command {
	var debug bool
	var jsonOutput bool

	var containerID string
	var rawDataDir string
	var reportID string

	rootCmd := &cobra.Command{
		Use:   "conjur-inspect",
		Short: "Qualification CLI for common Conjur Enterprise self-hosted issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				log.EnableDebugMode()
			}

			report, err := defaultReportConstructor(reportID, rawDataDir)
			if err != nil {
				return err
			}

			log.Debug("Running report...")
			result := report.Run(containerID)

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
			err = writer.Write(cmd.OutOrStdout(), &result)
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
	rootCmd.PersistentFlags().StringVarP(
		&containerID,
		"container-id",
		"", // No shorthand
		"", // No default
		"Conjur Enterprise container ID or name to inspect",
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

	rootCmd.PersistentFlags().StringVarP(
		&rawDataDir,
		"data-output-dir",
		"",  // No shorthand
		".", // Default is the current working directory
		"Where to save the raw data archive",
	)

	rootCmd.PersistentFlags().StringVarP(
		&reportID,
		"report-id",
		"", // No shorthand

		// This time stamp defines a custom format in golang, see here for more info:
		// https://pkg.go.dev/time#pkg-constants
		time.Now().Format("2006-01-02-15-04-05"), // Default is the current timestamp
		"Correlation identifier used for the raw data archive and report output",
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
