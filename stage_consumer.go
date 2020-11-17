package main

import (
	"sync"
)

const (
	diffStatusCodeError      = "diff-status-code"
	errorUnmarshalLeftError  = "error-unmarshal-left"
	errorUnmarshalRightError = "error-unmarshal-right"
	ok                       = "ok"
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
					streamConsumer <- validate(producerValue, consumer.Exclude)
				}
			}()
		}
		wg.Wait()
	}()
	return streamConsumer
}

func validate(hostsPair HostsPair, fieldsToExclude []string) StatusValidationError {
	var fieldErrorArray []string

	isOk, fieldError, statusCodes := isComparisonJsonResponseOk(hostsPair, fieldsToExclude)
	for !isOk {
		if !isFieldErrorBasic(fieldError) {
			fieldErrorArray = append(fieldErrorArray, fieldError)
			fieldsToExclude = append(fieldsToExclude, fieldError)
			isOk, fieldError, statusCodes = isComparisonJsonResponseOk(hostsPair, fieldsToExclude)
		} else {
			fieldErrorArray = append(fieldErrorArray, fieldError)
			break
		}
	}

	if len(fieldErrorArray) > 0 && fieldErrorArray[0] != ok {
		isOk = false
	} else {
		isOk = true
	}

	result := StatusValidationError{
		RelativePath:   hostsPair.RelativeURL,
		IsComparisonOk: isOk,
		FieldError:     fieldErrorArray,
		StatusCodes:    statusCodes,
	}
	return result
}

func isFieldErrorBasic(fieldError string) bool {
	switch fieldError {
	case
		diffStatusCodeError,
		errorUnmarshalLeftError,
		errorUnmarshalRightError:
		return true
	}
	return false
}

func isComparisonJsonResponseOk(hostsPair HostsPair, excludeFields []string) (bool, string, string) {

	statusCodes := hostsPair.getStatusCodes()

	if hostsPair.Has401() {
		panic("Authorization problem")
	}
	if hostsPair.HasErrors() || !hostsPair.EqualStatusCode() {
		fieldErrorCounter.Add("diff-status-code")
		return false, "diff-status-code", statusCodes
	}
	// Eli esta es la modificacion
	if !hostsPair.HasStatusCode200() {
		return true, "ok", statusCodes
	}
	leftJSON, err := unmarshal(hostsPair.Left.Body)
	if err != nil {
		fieldErrorCounter.Add("error-unmarshal-left")
		return false, "error-unmarshal-left", statusCodes
	}
	rightJSON, err := unmarshal(hostsPair.Right.Body)
	if err != nil {
		fieldErrorCounter.Add("error-unmarshal-right")
		return false, "error-unmarshal-right", statusCodes
	}

	if len(options.Exclude) > 0 {
		for _, excludeField := range excludeFields {
			Remove(leftJSON, excludeField)
			Remove(rightJSON, excludeField)
		}
	}
	isEqual, fieldError := Equal(leftJSON, rightJSON)
	if !isEqual {
		fieldErrorCounter.Add(fieldError)
		return false, fieldError, statusCodes
	}
	return true, "ok", statusCodes
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
