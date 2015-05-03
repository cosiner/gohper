package io2

import "io"

type WriterChain interface {
	io.Writer
	// Wrap add a new writer to top of chain
	Wrap(io.Writer)
	// Unwrap remove writer from top of chain, if there is only the base writer
	// nil was returned
	Unwrap() io.Writer
	// IsWrapped check whether a chain has wrapped writer
	IsWrapped() bool
	// Writer return the top writer
	Writer() io.Writer
	// BaseWriter return the base writer
	BaseWriter() io.Writer
}

// NewWriterChain create a new writer chain based on given writer
// parameter must not be nil, otherwise, nil was returned
func NewWriterChain(wr io.Writer) WriterChain {
	if wr == nil {
		return nil
	}
	return &writerChain{top: wr, base: wr, next: nil}
}

// WriterChain is a chain of writer, the last writer never removed
// it especially useful for wrap response, also it can be used in other
// similar scene, if need, just EMBED IT ANONYMOUSLY to get this ability
// data write to WriterChain is default write to the top writer of chain
type writerChain struct {
	top  io.Writer
	base io.Writer
	next *writerChain
}

func (wc *writerChain) Write(data []byte) (int, error) {
	return wc.top.Write(data)
}

// Wrap add a writer to top of the writer chain
func (wc *writerChain) Wrap(wr io.Writer) {
	newChain := &writerChain{
		top:  wc.top,
		next: wc.next,
	}
	wc.top = wr
	wc.next = newChain
}

// Unwrap remove a writer on the top of writer chain and return it,
// if there is only the base one, nil was returned
func (wc *writerChain) Unwrap() (w io.Writer) {
	if next := wc.next; next != nil {
		w = wc.top
		wc.top = next.top
		wc.next = next.next
	}
	return
}

// IsWrapped check whether writer chain is wrapped
func (wc *writerChain) IsWrapped() bool {
	return wc.next != nil
}

// Writer return the top writer of writerChain
func (wc *writerChain) Writer() io.Writer {
	return wc.top
}

// BaseWriter return the base writer of writer chain
func (wc *writerChain) BaseWriter() io.Writer {
	return wc.base
}
