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
	readerStream := pipeline.reader.Read()
	producerStream := pipeline.producer.Produce(readerStream)
	consumerStream := pipeline.consumer.Consume(producerStream)
	pipeline.writer.Write(consumerStream, ctx)
}
