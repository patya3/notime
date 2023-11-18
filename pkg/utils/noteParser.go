package utils

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Hit struct {
	Note       string
	LineNumber int
	Path       string
	Context    string
}

func ParseVimNotes() ([]Hit, error) {

	args := []string{
		"// NOTE: ",
		"/home/pmeszaros/aensys_work/",
		"-A",
		"10",
		"--no-ignore",
		"-g",
		"!node_modules",
		"-g",
		"!build",
		"-g",
		"!public",
		"-g",
		"!android",
		"-g",
		"!ios",
		"-g",
		"!dist",
		"-g",
		"!.git",
		"--no-messages",
		"-n",
	}

	cmd := exec.Command("rg", args...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()

	out := outb.String()

	hits := strings.Split(out, "--\n")

	var parsedHits []Hit

	for _, hit := range hits {
		lines := strings.Split(hit, "\n")

		firstLine := strings.Split(lines[0], ":")

		var path string
		var lineNumber int
		var note string
		var err error
		var context bytes.Buffer

		path, note = firstLine[0], strings.TrimSpace(firstLine[3])
		lineNumber, err = strconv.Atoi(firstLine[1])
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		for _, line := range lines[1:] {
			lineContent := strings.Split(line, " ")

			context.WriteString(strings.Join(lineContent[1:], " "))
			context.WriteString("\n")
		}

		parsedHits = append(parsedHits, Hit{Path: path, Note: note, LineNumber: lineNumber, Context: context.String()})

	}

	return parsedHits, nil

}
