package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
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

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Retrier usage: retrier \"command1; command2 && command3 || comand4 | command5\"")
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

	for {
		// Run the command
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// os.Exit(1)

		err := cmd.Run()
		if err == nil {
			// Command succeeded
			fmt.Printf("Command succeeded on attempt %d\n", attempt)
			os.Exit(0)
		}

		// Command failed
		fmt.Printf("Attempt %d failed `%s` Error: %s\n", attempt, command, err)

		// Check if max attempts is set and exceeded
		if *maxAttempts != -1 && attempt >= *maxAttempts {
			fmt.Printf("Failed after %d attempts\n", attempt)
			os.Exit(1)
		}

		// Calculate backoff delay based on the selected strategy
		delay, _ := parseBackoffStrategy(*backoffStrategy, attempt, *baseDelay)
		fmt.Printf("Retrying in %v...\n", delay)

		// Wait for the backoff delay before retrying
		time.Sleep(delay)

		// Increment attempt counter
		attempt++
	}
}
