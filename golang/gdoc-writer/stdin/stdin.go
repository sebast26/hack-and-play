package stdin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func readStdin() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	output := make([]byte, 0)

	for {
		n, err := reader.Read(buf[:cap(buf)])
		buf = buf[:n]

		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return nil, errors.New(fmt.Sprintf("Unexpected error when reading from stdio: %v", err))
		}

		output = append(output, buf...)

		if err != nil && err != io.EOF {
			return nil, errors.New(fmt.Sprintf("Unexpected error when reading from stdio: %v", err))
		}
	}

	return output, nil
}
