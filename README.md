Import the inputs for [Advent of Code](https://adventofcode.com/) challenges
straight into your code.

Fetched data is cached by default.

## Usage

1. Log into Advent of Code from your browser.
2. Find your session cookie and save the value to a file (e.g. `session_id`).
3. Init a reader and get your input:

```go
i, err := aocutil.NewInputFromFile("session_id")
if err != nil {
  log.Fatal(err)
}

lines, err := i.Strings(2018, 2)
if err != nil {
  log.Fatal(err)
}

// use lines
```

## Types

```go
func (i *Input) BigFloats(year, day int) ([]*big.Float, error)
func (i *Input) BigInts(year, day int, base int) ([]*big.Int, error)
func (i *Input) Bytes(year, day int) ([]byte, error)
func (i *Input) Floats(year, day int) ([]float64, error)
func (i *Input) Int64s(year, day int, base int) ([]int64, error)
func (i *Input) Ints(year, day int) ([]int, error)
func (i *Input) Reader(year, day int) (io.ReadCloser, error)
func (i *Input) Strings(year, day int) ([]string, error)
```

See godoc for more info.
