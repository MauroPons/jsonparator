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
	isOk, fieldError := isComparisonJsonResponseOk(hostsPair, consumer.Exclude)
	return StatusValidationError{
		RelativePath:   hostsPair.RelativeURL,
		IsComparisonOk: isOk,
		FieldError:     fieldError,
	}
}

func isComparisonJsonResponseOk(hostsPair HostsPair, excludeFields []string) (bool, string) {
	if hostsPair.Has401() {
		panic("Authorization problem")
	}
	if hostsPair.HasErrors() || !hostsPair.EqualStatusCode() {
		return false, "diff-status-code"
	}
	leftJSON, err := unmarshal(hostsPair.Left.Body)
	if err != nil {
		return false, "error-unmarshal-left"
	}
	rightJSON, err := unmarshal(hostsPair.Right.Body)
	if err != nil {
		return false, "error-unmarshal-right"
	}

	if len(options.Exclude) > 0 {
		for _, excludeField := range excludeFields {
			Remove(leftJSON, excludeField)
			Remove(rightJSON, excludeField)
		}
	}
	isEqual, FieldError := Equal(leftJSON, rightJSON)
	if !isEqual {
		fieldErrorCounter.Add(FieldError)
		return false, FieldError
	}
	return true, "ok"
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
	FieldError     string
}
