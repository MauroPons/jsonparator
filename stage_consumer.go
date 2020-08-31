package main

import "sync"

type Consumer struct {
	Exclude        string
}

func NewConsumer() *Consumer {
	return &Consumer{
		Exclude:   options.Exclude,
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
	return StatusValidationError {
		RelativePath:   hostsPair.RelativeURL,
		IsComparisonOk: isComparisonJsonResponseOk(hostsPair, consumer.Exclude),
	}
}

func isComparisonJsonResponseOk(hostsPair HostsPair, excludeField string) bool {
	if hostsPair.Has401() {
		panic("Authorization problem")
	}
	if hostsPair.HasErrors() || !hostsPair.EqualStatusCode() {
		return false
	}
	leftJSON, err := unmarshal(hostsPair.Left.Body)
	if err != nil {
		return false
	}
	rightJSON, err := unmarshal(hostsPair.Left.Body)
	if err != nil {
		return false
	}
	if options.Exclude != "" {
		Remove(leftJSON, excludeField)
		Remove(rightJSON, excludeField)
	}
	if !Equal(leftJSON, rightJSON) {
		return false
	}
	return true
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
}