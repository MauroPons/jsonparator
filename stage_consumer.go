package main

import (
	"sync"
)

type Consumer struct {
	Exclude []string
}

func NewConsumer() *Consumer {
	return &Consumer{
		Exclude: options.Exclude,
	}
}

func (consumer Consumer) Consume(streamProducer <-chan HostsPair) <-chan StatusValidationError {
	streamConsumer := make(chan StatusValidationError)
	go func() {
		defer close(streamConsumer)
		var wg sync.WaitGroup
		wg.Add(options.Currency)
		for w := 0; w < options.Currency; w++ {
			go func() {
				defer wg.Done()
				for producerValue := range streamProducer {
					streamConsumer <- consumer.validate(producerValue)
				}
			}()
		}
		wg.Wait()
	}()
	return streamConsumer
}

func (consumer Consumer) validate(hostsPair HostsPair) StatusValidationError {
	isOk, fieldError, statusCodes := isComparisonJsonResponseOk(hostsPair, consumer.Exclude)
	return StatusValidationError{
		RelativePath:   hostsPair.RelativeURL,
		IsComparisonOk: isOk,
		FieldError:     fieldError,
		StatusCodes:    statusCodes,
	}
}

func isComparisonJsonResponseOk(hostsPair HostsPair, excludeFields []string) (bool, []string, string) {

	statusCodes := hostsPair.getStatusCodes()

	if hostsPair.Has401() {
		panic("Authorization problem")
	}
	if hostsPair.HasErrors() || !hostsPair.EqualStatusCode() {
		fieldErrorCounter.Add("diff-status-code")
		return false, []string{"diff-status-code"}, statusCodes
	}
	// Eli esta es la modificacion
	if !hostsPair.HasStatusCode200() {
		return true, []string{"ok"}, statusCodes
	}
	leftJSON, err := unmarshal(hostsPair.Left.Body)
	if err != nil {
		fieldErrorCounter.Add("error-unmarshal-left")
		return false, []string{"error-unmarshal-left"}, statusCodes
	}
	rightJSON, err := unmarshal(hostsPair.Right.Body)
	if err != nil {
		fieldErrorCounter.Add("error-unmarshal-right")
		return false, []string{"error-unmarshal-right"}, statusCodes
	}

	if len(options.Exclude) > 0 {
		for _, excludeField := range excludeFields {
			Remove(leftJSON, excludeField)
			Remove(rightJSON, excludeField)
		}
	}

	isEqual := false
	fieldError := ""
	var errorTypeFields []string

	for !isEqual {
		if fieldError != "" {
			Remove(leftJSON, fieldError)
			Remove(rightJSON, fieldError)
			fieldError = ""
		}
		isEqual, fieldError = Equal(leftJSON, rightJSON)
		if !isEqual {
			errorTypeFields = append(errorTypeFields, fieldError)
		}
	}

	if len(errorTypeFields) > 0 {
		return false, errorTypeFields, statusCodes
	}

	//isEqual, fieldError := Equal(leftJSON, rightJSON)
	//if !isEqual {
	//	fieldErrorCounter.Add(fieldError)
	//	return false, fieldError, statusCodes
	//}
	return true, []string{"ok"}, statusCodes
}

func unmarshal(b []byte) (interface{}, error) {
	j, err := Unmarshal(b)
	if err != nil {
		return nil, err
	}
	return j, nil
}

type StatusValidationError struct {
	RelativePath   string
	IsComparisonOk bool
	FieldError     []string
	StatusCodes    string
}
