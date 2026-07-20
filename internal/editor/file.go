package editor

import (
	"bufio"
	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"io"
	"log"
	"os"
	"strings"
)

// this is a test comment for my editor

func OpenFile(path string) (*buffer.Buffer, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lines := make([]string, 0)

	for {
		line, err := reader.ReadString('\n')

		if len(line) > 0 {
			lines = append(lines, strings.TrimSuffix(line, "\n"))
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
	}

	log.Println(lines)

	buf := buffer.InitBuffer(lines)
	return buf, nil
}

func (e *Editor) SaveFile() error {
	file, err := os.OpenFile(e.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 0; i < e.Buffer.LineCount(); i++ {
		line, _ := e.Buffer.GetLine(i)

		_, err := writer.WriteString(line)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Don't write an extra newline after the last line.
		if i != e.Buffer.LineCount()-1 {
			_, err = writer.WriteString("\n")
			if err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	if err := writer.Flush(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
