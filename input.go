package aocutil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Input keeps track of your session ID so it can make requests.
type Input struct {
	SessionID string
}

// NewInputFromFile makes an Input instance using the entire contents of file f
// as the session key (minus any leading or trailing whitespace).
func NewInputFromFile(f string) (*Input, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Input{SessionID: strings.TrimSpace(string(b))}, nil
}

// Reader returns an io.ReadCloser of the input. This must be closed when done.
func (i *Input) Reader(year, day int) (io.ReadCloser, error) {
	fn := fmt.Sprintf("%d_%d.txt", year, day)
	f, err := os.Open(fn)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Not cached, fetching", year, day)

		// fetch
		rc, err := i.fetch(year, day)
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		// save to disk
		data, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		log.Println("Saving to disk")
		if err = os.WriteFile(fn, data, 0644); err != nil {
			// don't fail on save fail - we still have the data
			log.Println("Could not save input to disk:", err)
		}

		// we already have data in memory
		return io.NopCloser(bytes.NewReader(data)), nil
	} else if err != nil {
		return nil, err
	} else {
		log.Println("Using cached data", year, day)
		return f, nil
	}
}

func (i *Input) fetch(year, day int) (io.ReadCloser, error) {
	r, err := http.NewRequest("GET",
		fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	if err != nil {
		return nil, err
	}

	r.AddCookie(&http.Cookie{
		Name:  "session",
		Value: i.SessionID,
	})

	s, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if s.StatusCode != 200 {
		defer s.Body.Close()
		body, err := io.ReadAll(s.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}

	return s.Body, nil
}

// Bytes returns the entire input as a byte array.
func (i *Input) Bytes(year, day int) ([]byte, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

// Strings treats each line of the input as separate and returns them in an
// array.
func (i *Input) Strings(year, day int) ([]string, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	strs := make([]string, 0)
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		strs = append(strs, sc.Text())
	}
	return strs, nil
}

// Ints treats each line of the input as an integer and returns them in an
// array.
func (i *Input) Ints(year, day int) ([]int, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	ints := make([]int, 0)
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		i, err := strconv.Atoi(sc.Text())
		if err != nil {
			return ints, fmt.Errorf("At index %d: %w", len(ints), err)
		}
		ints = append(ints, i)
	}

	return ints, nil
}

// Int64s treats each line of the input as an integer and returns them in an
// array. This allows reading numbers in different bases; treat base according
// to strconv.ParseInt.
func (i *Input) Int64s(year, day int, base int) ([]int64, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	ints := make([]int64, 0)
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		i, err := strconv.ParseInt(sc.Text(), base, 64)
		if err != nil {
			return ints, fmt.Errorf("At index %d: %w", len(ints), err)
		}
		ints = append(ints, i)
	}

	return ints, nil
}

// Floats treats each line of the input as a floating point number and returns
// them in an array.
func (i *Input) Floats(year, day int) ([]float64, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	floats := make([]float64, 0)
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		f, err := strconv.ParseFloat(sc.Text(), 64)
		if err != nil {
			return floats, fmt.Errorf("At index %d: %w", len(floats), err)
		}
		floats = append(floats, f)
	}

	return floats, nil
}

// BigInts treats each line of the input as an integer and returns them in an
// array. This allows reading numbers in different bases; treat base according
// to strconv.ParseInt.
func (i *Input) BigInts(year, day int, base int) ([]*big.Int, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	ints := make([]*big.Int, 0)
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		s := sc.Text()
		z := new(big.Int)
		_, ok := z.SetString(s, base)
		if !ok {
			return ints, fmt.Errorf("At index %d: could not parse %s", len(ints), s)
		}
		ints = append(ints, z)
	}

	return ints, nil
}

// BigFloats treats each line of the input as a floating point number  and
// returns them in an array.
func (i *Input) BigFloats(year, day int) ([]*big.Float, error) {
	rc, err := i.Reader(year, day)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	floats := make([]*big.Float, 0)
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		s := sc.Text()
		z := new(big.Float)
		_, ok := z.SetString(s)
		if !ok {
			return floats, fmt.Errorf("At index %d: could not parse %s", len(floats), s)
		}
		floats = append(floats, z)
	}

	return floats, nil
}
