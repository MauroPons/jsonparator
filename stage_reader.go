package main

import (
	"bufio"
	"net/url"
	"os"
	"strings"
)

type Reader struct {
	file  *os.File
	hosts []string
}

type URLPair struct {
	pathLineNumber    int
	RelativePath      string
	UrlLeft, UrlRight URL
}

type URL struct {
	URL   *url.URL
	Error error
}

func NewReader() *Reader {
	return &Reader{
		file:  fileRelativePathSource,
		hosts: options.Hosts,
	}
}

func (reader *Reader) Read() <-chan URLPair {
	streamReader := make(chan URLPair)
	go func() {
		defer close(streamReader)
		_, _ = reader.file.Seek(0, 0)
		scanner := bufio.NewScanner(reader.file)
		count := 0
		for scanner.Scan() {
			relativePath := strings.Trim(scanner.Text(), "\"")
			if !strings.Contains(relativePath, "request_ur") {
				count++
				leftUrl := URL{}
				leftUrl.URL, leftUrl.Error = getUrl(reader.hosts[0], relativePath)

				rightUrl := URL{}
				rightUrl.URL, rightUrl.Error = getUrl(reader.hosts[1], relativePath)

				streamReader <- URLPair{
					pathLineNumber: count,
					RelativePath:   relativePath,
					UrlLeft:        leftUrl,
					UrlRight:       rightUrl,
				}
			}
		}
	}()
	return streamReader
}

func getUrl(host string, relativePath string) (*url.URL, error) {
	url, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}
	url.RawQuery = url.Query().Encode()
	baseUrl, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return baseUrl.ResolveReference(url), nil
}
