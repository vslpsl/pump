package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vslpsl/pump/internal/datasource/file"
	"github.com/vslpsl/pump/internal/datasource/stream"
	"github.com/vslpsl/pump/internal/limiter"
	"github.com/vslpsl/pump/internal/pipe"
	"io"
	"log"
	"os"
)

const (
	DefaultBufferSize = 1024 * 32
	DefaultRateLimit  = -1
)

var (
	sourcePath string
	targetPath string
	bufferSize int64
	rateLimit  int64
)

var pumpCmd = &cobra.Command{
	Use:           "pump",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		source, err := getSource()
		if err != nil {
			return err
		}
		defer func() { _ = source.Close() }()

		target, err := getTarget()
		if err != nil {
			return err
		}
		defer func() { _ = target.Close() }()
		p := pipe.New(source, target, bufferSize)

		return p.Pump()
	},
}

func init() {
	pumpCmd.PersistentFlags().StringVar(&sourcePath, "source", "", "path to source, if not set os.Stdin is used")
	pumpCmd.PersistentFlags().StringVar(&targetPath, "target", "", "path to target, if not set os.Stdout is used")
	pumpCmd.PersistentFlags().Int64Var(&bufferSize, "buffer-size", DefaultBufferSize, "size of buffer")
	pumpCmd.PersistentFlags().Int64Var(&rateLimit, "rate-limit", DefaultRateLimit, "limit number of bytes per second")
}

func Execute() {
	if err := pumpCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func getSource() (io.ReadCloser, error) {
	if sourcePath == "" {
		return stream.NewReader(os.Stdin), nil
	}

	return file.NewReader(sourcePath)
}

func getTarget() (io.WriteCloser, error) {
	if targetPath == "" {
		return stream.NewWriter(os.Stdout), nil
	}

	return file.NewWriter(targetPath)
}

func getLimiter() pipe.Limiter {
	if rateLimit <= 0 {
		return nil
	}

	return limiter.New(rateLimit)
}
