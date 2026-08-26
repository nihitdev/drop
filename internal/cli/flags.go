package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"time"
)

var ErrUsage = errors.New("usage requested")

type Config struct {
	Port   int
	Keep   bool
	Expiry time.Duration
	Path   string
}

func Parse(args []string, output io.Writer) (Config, error) {
	var cfg Config
	flags := flag.NewFlagSet("drop", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.IntVar(&cfg.Port, "port", 0, "port to listen on (0 = automatic)")
	flags.BoolVar(&cfg.Keep, "keep", false, "keep server running after download")
	flags.DurationVar(&cfg.Expiry, "expires", 0, "expire after duration, e.g. 10m")
	flags.Usage = func() { printUsage(output) }

	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	if flags.NArg() == 0 {
		printUsage(output)
		return Config{}, ErrUsage
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("port must be between 0 and 65535")
	}
	if cfg.Expiry < 0 {
		return Config{}, fmt.Errorf("expiry must not be negative")
	}

	cfg.Path = flags.Arg(0)
	return cfg, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: drop [options] <file-or-directory>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "examples:")
	fmt.Fprintln(w, "  drop file.zip")
	fmt.Fprintln(w, "  drop ./folder")
	fmt.Fprintln(w, "  drop --keep file.zip")
	fmt.Fprintln(w, "  drop --expires 10m file.zip")
	fmt.Fprintln(w, "  drop --port 9000 file.zip")
}
