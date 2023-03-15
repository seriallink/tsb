package cmd

import (
	"errors"
	"os"

	"github.com/seriallink/timescale/src"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flags struct {
	input   string
	workers int
}

func InitCmd() *cobra.Command {

	cmd := cobra.Command{
		Use: "tsb",
		RunE: func(cmd *cobra.Command, args []string) error {

			logrus.SetOutput(os.Stderr)
			logrus.SetLevel(logrus.InfoLevel)

			if flags.workers < 1 {
				return errors.New("workers must be a positive number")
			}

			reader := cmd.InOrStdin()
			if flags.input != "" {
				file, err := os.Open(flags.input)
				if err != nil {
					return err
				}
				defer file.Close()
				reader = file
			}

			err := src.InitConnectionPool()
			if err != nil {
				return err
			}
			defer src.CloseConnectionPool()

			src.NewTimescaleBenchmark(reader, flags.workers).Do()
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&flags.input, "input", "i", "", "csv file name (leave empty to use stdin)")
	cmd.PersistentFlags().IntVarP(&flags.workers, "workers", "w", 10, "number of concurrent workers")

	return &cmd

}
