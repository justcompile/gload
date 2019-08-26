/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"time"

	"github.com/justcompile/gload/internal"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run URL",
	Short: "Runs a load test against a given URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		duration, _ := cmd.Flags().GetDuration("duration")
		output, _ := cmd.Flags().GetString("output")
		workers, _ := cmd.Flags().GetInt("workers")

		d := internal.NewDispatcher(
			&internal.Options{
				Duration: duration,
				Output:   output,
				Workers:  workers,
			},
		)

		defer d.Close()

		d.Run(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().DurationP("duration", "d", time.Second*30, "Duration for the test")
	runCmd.Flags().StringP("output", "o", "-", "File containing test results [default: stdout]")
	runCmd.Flags().IntP("workers", "w", 10, "Number of workers to execute")

}
