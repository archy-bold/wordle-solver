# Wordle Command Line

A Wordle implementation for the command line, written in go.

## Get It

Download the latest version for your system from the [releases](https://github.com/archy-bold/wordle-go/releases) page.

Or clone the repo and run from there. You may need to re-build for your system first.

```bash
git clone git@github.com:archy-bold/wordle-go.git
cd wordle-go
./wordle
```

## Build

```bash
go build -o wordle
```

## Run

Run the command with no arguments to play a game with today's word.

```bash
./wordle
```

### Options

- `-cheat` - Runs in solve mode to work out an existing wordle. Follow the instructions to enter your results and receive suggested words to play.
- `-random` - Choose a random word instead of today's
- `-auto` - Automatically completes the puzzle
- `-date=[2021-12-31]` - Set the winning word from a specific date
- `-word=[answer]` - Set the winning word with this argument.
- `-all` - Runs the auto-solver through every permutation, giving results when complete.
- `-starter=[word]` - Specify the starter word for strategies
