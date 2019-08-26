package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/justcompile/gload/internal"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gload url",
	Short: "gload - load testing with go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		duration, _ := cmd.Flags().GetDuration("duration")
		workers, _ := cmd.Flags().GetInt("workers")

		d := internal.NewDispatcher(
			&internal.Options{
				Duration: duration,
				Workers:  workers,
			},
		)

		d.Run(args[0])
	},
}

// Execute does the thing
func Execute() {
	rootCmd.Flags().DurationP("duration", "d", time.Second*30, "Duration for the test")
	rootCmd.Flags().IntP("workers", "w", 10, "Number of workers to execute")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
