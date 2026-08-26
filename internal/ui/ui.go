package ui

import (
	"fmt"
	"os"
	"time"

	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/nihitdev/drop/internal/server"
	"github.com/nihitdev/drop/internal/share"
)

func ShowShare(target *share.Target, address, ip string, port int, keep bool, expiry time.Duration) {
	fmt.Println()
	fmt.Println("╭────────────────────────────────────╮")
	fmt.Println("│              drop                  │")
	fmt.Println("│        local file sharing          │")
	fmt.Println("╰────────────────────────────────────╯")
	fmt.Println()
	fmt.Printf("  File      %s\n", target.DisplayName)
	fmt.Printf("  Size      %s\n", HumanSize(target.Size))
	fmt.Printf("  Address   %s:%d\n", ip, port)
	if keep {
		fmt.Println("  Mode      unlimited downloads")
	} else {
		fmt.Println("  Mode      one download")
	}
	if expiry > 0 {
		fmt.Printf("  Expires   %s\n", expiry)
	}
	fmt.Printf("\n  %s\n\n", address)
	qrterminal.GenerateWithConfig(address, qrterminal.Config{
		Level: qrterminal.M, Writer: os.Stdout, HalfBlocks: true, QuietZone: 1,
	})
	fmt.Println("\n  waiting for download...")
	fmt.Println("  Ctrl+C to stop")
	fmt.Println()
}

func HumanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit && exp < 3; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), []string{"KB", "MB", "GB", "TB"}[exp])
}

func ClientConnected(ip string)    { fmt.Printf("\n→ %s connected\n", ip) }
func TransferComplete(name string) { fmt.Printf("✓ %s transferred\n", name) }

func ServerStopped(reason server.StopReason) {
	if reason == server.StoppedAfterExpiry {
		fmt.Println("\n⌛ share expired")
	}
	fmt.Println("✓ server stopped")
}
