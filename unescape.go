package main

import (
	"fmt"
	"io"
)

func unescape(input io.Reader, output io.Writer) error {
	c := make(chan byte)

	var inputErr error

	go func() {
		defer close(c)

		inputDone := false
		inputBuffer := make([]byte, 4096)
		for !inputDone {
			n, err := input.Read(inputBuffer[:])
			if err != nil {
				inputDone = true
				if err != io.EOF {
					inputErr = err
				}
			}

			for _, b := range inputBuffer[:n] {
				c <- b
			}
		}
	}()

	backslash := false
	for b := range c {
		if backslash {
			switch b {
			case 'a', 'b', 'e', 'f', 'r', 'v':
			case 'n':
				if _, err := output.Write([]byte("\n")); err != nil {
					return fmt.Errorf("output error: %s", err)
				}
			case 't':
				if _, err := output.Write([]byte("\t")); err != nil {
					return fmt.Errorf("output error: %s", err)
				}
			case '\\', '\'', '"', '?':
				if _, err := output.Write([]byte{b}); err != nil {
					return fmt.Errorf("output error: %s", err)
				}
			default:
				if _, err := output.Write([]byte{'\\', b}); err != nil {
					return fmt.Errorf("output error: %s", err)
				}
			}

			backslash = false
			continue
		} else {
			if b == '\\' {
				backslash = true
				continue
			}
			if _, err := output.Write([]byte{b}); err != nil {
				return fmt.Errorf("output error: %s", err)
			}
		}
	}

	if inputErr != nil {
		return fmt.Errorf("input error: %s", inputErr)
	}

	return nil
}
