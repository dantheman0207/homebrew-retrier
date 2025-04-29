package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
)

var (
	Version = "0.1.16"
)

// Fibonacci backoff
func fibonacciBackoff(attempt int, baseDelay int) time.Duration {
	if attempt <= 2 {
		return time.Duration(baseDelay) * time.Second
	}
	a, b := baseDelay, baseDelay
	for i := 3; i <= attempt; i++ {
		a, b = b, a+b
	}
	return time.Duration(b) * time.Second
}

// Exponential backoff
func exponentialBackoff(attempt int, baseDelay int) time.Duration {
	return time.Duration(baseDelay*int(math.Pow(2, float64(attempt-1)))) * time.Second
}

// Linear backoff
func linearBackoff(attempt int, baseDelay int) time.Duration {
	return time.Duration(baseDelay*attempt) * time.Second
}

// Constant backoff
func constantBackoff(baseDelay int) time.Duration {
	return time.Duration(baseDelay) * time.Second
}

// Parse the backoff strategy based on the flag
func parseBackoffStrategy(strategy string, attempt int, baseDelay int) (time.Duration, string) {
	switch strategy {
	case "f", "fibonacci":
		return fibonacciBackoff(attempt, baseDelay), "fibonacci"
	case "e", "exponential":
		return exponentialBackoff(attempt, baseDelay), "exponential"
	case "l", "linear":
		return linearBackoff(attempt, baseDelay), "linear"
	case "c", "constant":
		return constantBackoff(baseDelay), "constant"
	default:
		fmt.Println("Unknown backoff strategy. Using Fibonacci as default.")
		return fibonacciBackoff(attempt, baseDelay), "fibonacci"
	}
}

func main() {
	// Define command line flags
	backoffStrategy := flag.String("backoff", "fibonacci", "(-b) Backoff strategy: fibonacci (f), exponential (e), linear (l), constant (c)")
	backoffStrategyShort := flag.String("b", "f", "")

	baseDelay := flag.Int("delay", 2, "(-d) Base delay in seconds for backoff")
	baseDelayShort := flag.Int("d", 2, "")

	maxAttempts := flag.Int("max-attempts", -1, "(-m) Maximum number of attempts (-1 for infinite retries)")
	maxAttemptsShort := flag.Int("m", -1, "")

	retryOnSuccess := flag.Bool("retry-on-success", false, "(-r) Retry on success (default: false)")
	retryOnSuccessShort := flag.Bool("r", false, "")

	version := flag.Bool("version", false, "(-v) Print version and exit")
	flag.BoolVar(version, "v", false, "")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Retrier usage: retrier \"command1; command2 && command3 || command4 | command5\"")
		flag.VisitAll(func(f *flag.Flag) {
			// Only show help for flags with descriptions (i.e., the long versions)
			if f.Usage != "" {
				fmt.Fprintf(os.Stderr, "  -%s\n", f.Name)
				fmt.Fprintf(os.Stderr, "        %s (default %q)\n", f.Usage, f.DefValue)
			}
		})
	}

	// Parse flags
	flag.Parse()

	if *version {
		fmt.Println("Retrier version", Version)
		os.Exit(0)
	}

	// use short versions
	if backoffStrategyShort != nil && *backoffStrategyShort != "f" {
		*backoffStrategy = *backoffStrategyShort
	}
	if baseDelayShort != nil && *baseDelayShort != 2 {
		*baseDelay = *baseDelayShort
	}
	if maxAttemptsShort != nil && *maxAttemptsShort != -1 {
		*maxAttempts = *maxAttemptsShort
	}

	// Get the command and its arguments
	if flag.NArg() < 1 {
		fmt.Println("You must provide a command to execute")
		os.Exit(1)
	}

	_, strategy := parseBackoffStrategy(*backoffStrategy, 0, *baseDelay)
	fmt.Printf("Using %s strategy for backoffs with initial delay %ds and %d max attempts\n", strategy, *baseDelay, *maxAttempts)

	command := strings.Join(flag.Args(), " ")

	// Initialize attempt counter
	attempt := 1

	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	for {
		// Run the command
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Start the command with a pty
		ptmx, err := pty.Start(cmd)
		if err != nil {
			fmt.Printf("failed to start PTY: %v\n", err)
		}
		defer func() { _ = ptmx.Close() }() // Best effort

		// Copy PTY output to real stdout
		go func() {
			_, _ = io.Copy(os.Stdout, ptmx)
		}()
		// os.Exit(1)

		err = cmd.Wait()
		if err == nil {
			// Command succeeded
			fmt.Printf("\nCommand succeeded on attempt %d\n", attempt)
			if !*retryOnSuccess && !*retryOnSuccessShort {
				os.Exit(0)
			}
		} else {
			// Command failed
			fmt.Printf("Attempt %d failed `%s` Error: %s\n", attempt, command, err)
		}

		// Check if max attempts is set and exceeded
		if *maxAttempts != -1 && attempt >= *maxAttempts {
			fmt.Printf("Finished after %d attempts\n", attempt)
			os.Exit(1)
		}

		// Handle Ctrl+D during the delay
		go func() {
			buf := make([]byte, 1)
			for {
				_, err := os.Stdin.Read(buf)
				if err == io.EOF {
					fmt.Println("\nCtrl+D detected. Exiting...")
					os.Exit(0)
				}
			}
		}()

		// Calculate backoff delay based on the selected strategy
		delay, _ := parseBackoffStrategy(*backoffStrategy, attempt, *baseDelay)
		fmt.Printf("Waiting for %v\n", delay)
		for i := int(delay.Seconds()) - 1; i > 0; i-- {
			select {
			case <-sigChan:
				// Handle Ctrl+C: Skip the current iteration
				fmt.Println("\nCtrl+C detected. Skipping current attempt...")
				attempt++
				break
			default:
				fmt.Printf("\r\033[KRetrying in %v...", time.Duration(i+1)*time.Second)
				time.Sleep(time.Second)
			}
		}
		fmt.Print("\r\033[K\n")

		// Increment attempt counter
		attempt++

	}
}
