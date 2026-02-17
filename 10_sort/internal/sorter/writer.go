package sorter

import (
	"bufio"
	"os"
)

func Write(records []Record) error {
	writer := bufio.NewWriter(os.Stdout)

	for _, record := range records {
		if _, err := writer.WriteString(record.Line + "\n"); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
