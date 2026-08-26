package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nihitdev/drop/internal/cli"
	"github.com/nihitdev/drop/internal/network"
	"github.com/nihitdev/drop/internal/server"
	"github.com/nihitdev/drop/internal/share"
	"github.com/nihitdev/drop/internal/token"
	"github.com/nihitdev/drop/internal/ui"
)

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, cli.ErrUsage) {
			fmt.Fprintf(os.Stderr, "drop: %v\n", err)
		}
		os.Exit(1)
	}
}

func run() error {
	cfg, err := cli.Parse(os.Args[1:], os.Stderr)
	if err != nil {
		return err
	}

	target, err := share.Prepare(cfg.Path)
	if err != nil {
		return fmt.Errorf("failed to prepare share target: %w", err)
	}
	defer target.Cleanup()

	shareToken, err := token.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate download token: %w", err)
	}

	ip, err := network.LocalIP()
	if err != nil {
		return fmt.Errorf("failed to determine LAN IP: %w", err)
	}

	listener, err := network.Listen(cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	address := network.URL(ip, listener, shareToken)
	ui.ShowShare(target, address, ip, network.Port(listener), cfg.Keep, cfg.Expiry)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	reason, err := server.Run(ctx, server.Config{
		Listener: listener,
		Target:   target,
		Token:    shareToken,
		Keep:     cfg.Keep,
		Expiry:   cfg.Expiry,
		OnClient: ui.ClientConnected,
		OnSent:   ui.TransferComplete,
	})
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	ui.ServerStopped(reason)
	return nil
}
