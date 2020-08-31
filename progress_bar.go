package main

import (
	"gopkg.in/cheggaaa/pb.v1"
	"sync"
)

type ProgressBar struct {
	totalPb *pb.ProgressBar
	okPb    *pb.ProgressBar
	errorPb *pb.ProgressBar
	pool    *pb.Pool
}

func NewProgressBar() *ProgressBar {
	totalPb := makeProgressBar(options.FilePathTotalLines, "TOTAL")
	okPb := makeProgressBar(options.FilePathTotalLines, "OK")
	errorPb := makeProgressBar(options.FilePathTotalLines, "ERROR")
	return &ProgressBar{
		totalPb,okPb,errorPb,nil,
	}
}

func makeProgressBar(total int, prefix string) *pb.ProgressBar {
	progressBar := pb.New(total)
	progressBar.Prefix(prefix)
	progressBar.SetMaxWidth(120)
	progressBar.SetRefreshRate(1000)
	progressBar.ShowElapsedTime = true
	progressBar.ShowTimeLeft = false
	return progressBar
}

func (p *ProgressBar) IncrementOk() {
	p.okPb.Add(1)
}

func (p *ProgressBar) IncrementError() {
	p.errorPb.Add(1)
}

func (p *ProgressBar) IncrementTotal() {
	p.totalPb.Add(1)
}

func (p *ProgressBar) Start() {
	pool, err := pb.StartPool(p.totalPb, p.okPb, p.errorPb)
	if err != nil {
		panic(err)
	}
	p.pool = pool
	p.okPb.Start()
}

func (p *ProgressBar) Stop() {
	wg := new(sync.WaitGroup)
	for _, bar := range []*pb.ProgressBar{p.totalPb, p.okPb, p.errorPb} {
		wg.Add(1)
		go func(cb *pb.ProgressBar) {
			cb.Finish()
			wg.Done()
		}(bar)
	}
	wg.Wait()
	// close pool
	_ = p.pool.Stop()
}