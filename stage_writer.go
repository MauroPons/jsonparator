package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type Writer struct {
	fileError	*os.File
	progressBar *ProgressBar
}

func NewWriter() *Writer {
	return &Writer{
		fileError:   fileError,
		progressBar: progressBar,
	}
}

func (writer Writer) Write(streamProducer <-chan StatusValidationError, context context.Context) {
	for producerValue := range streamProducer {
		progressBar.IncrementTotal()
		if producerValue.IsComparisonOk {
			progressBar.IncrementOk()
		}else{
			progressBar.IncrementError()
			w := bufio.NewWriter(writer.fileError)
			fmt.Fprintln(w, producerValue.RelativePath)
			_ = w.Flush()
		}
	}
	select {
		case <-context.Done():
	}
}
