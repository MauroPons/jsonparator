package main

import "context"

type Pipeline struct {
	reader   *Reader
	producer *Producer
	consumer *Consumer
	writer   *Writer
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		reader:   NewReader(),
		producer: NewProducer(),
		consumer: NewConsumer(),
		writer:   NewWriter(),
	}
}

func (pipeline *Pipeline) Run(ctx context.Context) {
	streamReader := pipeline.reader.Read()
	streamProducer := pipeline.producer.Produce(streamReader)
	streamConsumer := pipeline.consumer.Consume(streamProducer)
	pipeline.writer.Write(streamConsumer, ctx)
}
