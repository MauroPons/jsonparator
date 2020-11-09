package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type Writer struct {
	fileError   *os.File
	progressBar *ProgressBar
}

func NewWriter() *Writer {
	return &Writer{
		fileError:   fileRelativePathError,
		progressBar: progressBar,
	}
}

func (writer Writer) Write(consumerStream <-chan StatusValidationError, context context.Context) {
	countError := 0
	for producerValue := range consumerStream {
		if producerValue.IsComparisonOk {
			progressBar.IncrementOk()
		} else {
			countError++
			progressBar.IncrementError()
			addRelativePathToFileError(writer, producerValue.RelativePath)
			addRelativePathToFileParam(producerValue.RelativePath, "error")
			addRelativePathToFileTypeErrorArray(producerValue.FieldError, producerValue.RelativePath)
			if producerValue.StatusCodes != "200-200" {
				addRelativePathToFileTypeError(producerValue.StatusCodes, producerValue.RelativePath)
			}
		}
	}

	options.FilePathTotalLinesError = countError

	select {
	case <-context.Done():
	case <-consumerStream:
	}

}

func addRelativePathToFileTypeErrorArray(fieldError []string, relativePath string) {
	for _, value := range fieldError {
		addRelativePathToFileTypeError(value, relativePath)
	}
}

func addRelativePathToFileTypeError(fieldError string, relativePath string) {
	file := mapFileParams[fieldError]
	if file == nil {
		file = createFilesByTypeError(fieldError)
	}
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, relativePath)
	_ = w.Flush()
}

func addRelativePathToFileError(writer Writer, relativePath string) {
	w := bufio.NewWriter(writer.fileError)
	fmt.Fprintln(w, relativePath)
	_ = w.Flush()
}

func createFilesByTypeError(fieldError string) *os.File {
	relativePaths := options.BasePath + "type-error/"
	_ = os.Mkdir(relativePaths, 0777)
	pathFileParam := relativePaths + fieldError + ".txt"
	file, _ := os.Create(pathFileParam)
	mapFileParams[fieldError] = file
	return file
}
