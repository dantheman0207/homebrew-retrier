# Retrier

`retrier` is a command-line tool that retries a given command with a specified backoff strategy, allowing you to control the delay between retries and the maximum number of attempts.

## Installation

You can install `retrier` using [Homebrew](https://brew.sh/). Follow these instructions to install it via a Homebrew tap:

### Install via Homebrew

To install `retrier`, first add the Homebrew tap for the repository:

```bash
brew tap dantheman0207/retrier
```

Once the tap is added, you can install the tool using:

```bash
brew install retrier
```

After installation, you can check if `retrier` is available by running:

```bash
retrier --help
```

### Install via git

```bash
git clone git@github.com:dantheman0207/homebrew-retrier.git
cd homebrew-retrier
go build
./retrier -h
```

## Usage

The `retrier` command accepts a command to run and provides various options for configuring the backoff strategy, delay, and maximum number of attempts.

### Example Usage

Here are some examples of how you can use `retrier`:

#### Basic Usage

Retry a command using the default Fibonacci backoff strategy:

```bash
retrier -b f "echo 'Retrying...'"
```

#### Exponential Backoff with Custom Delay

Retry a command using exponential backoff with a 3-second base delay:

```bash
retrier -b e -d 3 "echo 'Retrying...'"
```

#### Infinite Retries

Retry a command infinitely until it succeeds:

```bash
retrier -m -1 "echo 'Retrying...'"
```

### Command-Line Flags

Below is the help output for the tool, which shows all the available flags and their descriptions.

```
Retrier usage: retrier "command1; command2 && command3 || comand4 | command5"
  -backoff
        (-b) Backoff strategy: fibonacci (f), exponential (e), linear (l), constant (c) (default "fibonacci")
  -delay
        (-d) Base delay in seconds for backoff (default "2")
  -max-attempts
        (-m) Maximum number of attempts (-1 for infinite retries) (default "-1")
```

### Backoff Strategies

You can specify different backoff strategies using the `-b` or `-backoff` flags:

- `fibonacci` (alias `f`): Backoff delay follows the Fibonacci sequence.
- `exponential` (alias `e`): Backoff delay doubles with each retry.
- `linear` (alias `l`): Backoff delay increases linearly.
- `constant` (alias `c`): Backoff delay remains constant between retries.

### Delay

The base delay can be set using `-d` or `-delay`. This is the initial delay (in seconds) before retries, and how it changes depends on the backoff strategy.

### Maximum Attempts

You can limit the number of retry attempts using `-m` or `-max-attempts`. By default, the tool retries indefinitely (`-1`), but you can set a specific limit.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.