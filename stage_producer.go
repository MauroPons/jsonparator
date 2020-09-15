package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Producer struct {
	client *http.Client
}

type HostsPair struct {
	RelativeURL string
	Errors      []error
	Left, Right Host
}

type Host struct {
	StatusCode int
	Body       []byte
	URL        *url.URL
	Error      error
}

func NewProducer() *Producer {
	return &Producer{
		client: getHttpClient(),
	}
}

func getHttpClient() *http.Client {
	return &http.Client{}
}

func (producer Producer) Produce(streamReader <-chan URLPair) <-chan HostsPair {
	streamProducer := make(chan HostsPair)
	go func() {
		defer close(streamProducer)
		var wg sync.WaitGroup
		wg.Add(options.Currency)
		velocityRPS := 60000 / (options.Velocity)
		limiter := time.Tick(time.Duration(velocityRPS) * time.Millisecond)
		for w := 0; w < options.Currency; w++ {
			go func() {
				defer wg.Done()
				for readerValue := range streamReader {
					<-limiter
					if !strings.Contains(readerValue.RelativePath, "request_ur") {
						streamProducer <- producer.process(readerValue)
					}
				}
			}()
		}
		wg.Wait()
	}()
	return streamProducer
}

func (producer *Producer) process(pair URLPair) HostsPair {
	work := func(url URL) <-chan Host {
		channelHost := make(chan Host)
		go func() {
			defer close(channelHost)
			channelHost <- producer.fetch(url)
		}()
		return channelHost
	}

	leftCh := work(pair.UrlLeft)
	rightCh := work(pair.UrlRight)

	lHost := <-leftCh
	rHost := <-rightCh

	response := HostsPair{
		RelativeURL: pair.RelativePath,
		Left:        lHost,
		Right:       rHost,
	}

	if lHost.Error != nil {
		response.Errors = append(response.Errors, lHost.Error)
	}

	if rHost.Error != nil {
		response.Errors = append(response.Errors, rHost.Error)
	}

	return response
}

func (producer Producer) fetch(url URL) Host {
	host := Host{}
	request, _ := http.NewRequest("GET", url.URL.String(), nil)
	for _, value := range options.Headers {
		splitKey := strings.Split(value, ":")
		request.Header.Add(splitKey[0], splitKey[1])
	}
	// try 1
	response, err := producer.client.Do(request)
	if err != nil {
		host.Error = err
		return host
	}
	if response.StatusCode != 200 {
		// retry 2
		response, err = producer.client.Do(request)
		if err != nil {
			host.Error = err
			return host
		}
	}

	host.URL = url.URL
	host.Body, _ = ioutil.ReadAll(response.Body)
	host.StatusCode = response.StatusCode
	return host
}

func (hostsPair HostsPair) HasErrors() bool {
	return len(hostsPair.Errors) > 0
}

func (hostsPair HostsPair) EqualStatusCode() bool {
	return hostsPair.Left.StatusCode == hostsPair.Right.StatusCode
}

func (hostsPair HostsPair) Has401() bool {
	return hostsPair.Left.StatusCode == 401 || hostsPair.Right.StatusCode == 401
}

func (hostsPair HostsPair) getStatusCodes() string {
	return fmt.Sprintf("%d-%d", hostsPair.Left.StatusCode, hostsPair.Right.StatusCode)
}
