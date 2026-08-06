package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type WAL struct {
	file *os.File
}

func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{file: file}, nil
}

func (w *WAL) WriteSet(key, value string) error {
	line := fmt.Sprintf("SET %s %s\n", key, value)
	_, err := w.file.WriteString(line)
	return err
}

func (w *WAL) WriteDelete(key string) error {
	line := fmt.Sprintf("DEL %s\n", key)
	_, err := w.file.WriteString(line)
	return err
}

func (w *WAL) Replay() (map[string]Entry, error) {
	data := make(map[string]Entry)

	scanner := bufio.NewScanner(w.file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "SET":
			if len(parts) != 3 {
				continue
			}
			data[parts[1]] = Entry{Key: parts[1], Value: parts[2]}
		case "DEL":
			if len(parts) != 2 {
				continue
			}
			data[parts[1]] = Entry{Key: parts[1], Deleted: true}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (w *WAL) Close() error {
	if w == nil || w.file == nil {
		return nil
	}
	return w.file.Close()
}
