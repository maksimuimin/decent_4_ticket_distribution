package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
)

var (
	file      = flag.String("file", "", "path to file with students' names")
	numbilets = flag.Uint64("numbilets", 0, "number of tickets to distribute")
	parameter = flag.Uint64("parameter", 0, "random number generator seed")
)

func processLine(r *bufio.Reader, h hash.Hash) error {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("r.ReadBytes: %w", err)
	}

	if _, err = h.Write(line); err != nil {
		return fmt.Errorf("h.Write: %w", err)
	}

	ticket, _ := binary.Uvarint(h.Sum(nil))
	ticket = ticket%*numbilets + 1
	fmt.Printf("%s: %d\n", line[:len(line)-1], ticket)

	return nil
}

func main() {
	flag.Parse()

	fh, err := os.Open(*file)
	if err != nil {
		log.Fatalf("os.Open: %s\n", err)
	}
	defer fh.Close()

	r := bufio.NewReader(fh)

	seed := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(seed, *parameter)
	h := hmac.New(sha256.New, seed)

	for {
		err = processLine(r, h)
		if err != nil {
			break
		}
	}

	if errors.Is(err, io.EOF) {
		err = nil
	}

	if err != nil {
		log.Fatalf("processLine: %s\n", err)
	}
}
