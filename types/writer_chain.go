package types

import "io"

// WriterChain is a chain of writer, the last writer never removed
// it especially useful for wrap response, also it can be used in other
// similar scene, if need, just EMBED IT ANONYMOUSLY to get this ability
// data write to WriterChain is default write to the top writer of chain
type WriterChain struct {
	io.Writer
	base io.Writer
	next *WriterChain
}

// NewWriterChain create a new writer chain based on given writer
// parameter must not be nil, otherwise, nil was returned
func NewWriterChain(wr io.Writer) *WriterChain {
	if wr == nil {
		return nil
	}
	return &WriterChain{Writer: wr, base: wr, next: nil}
}

// Wrap add a writer to top of the writer chain
func (wc *WriterChain) Wrap(wr io.Writer) {
	newChain := &WriterChain{
		Writer: wc.Writer,
		next:   wc.next,
	}
	wc.Writer = wr
	wc.next = newChain
}

// Unwrap remove a writer on the top of writer chain and return it,
// if there is only the base one, nil was returned
func (wc *WriterChain) Unwrap() (w io.Writer) {
	if next := wc.next; next != nil {
		w = wc.Writer
		wc.Writer = next.Writer
		wc.next = next.next
	}
	return
}

// IsWrapped check whether writer chain is wrapped
func (wc *WriterChain) IsWrapped() bool {
	return wc.next != nil
}

// Writer return the base writer of WriterChain
func (wc *WriterChain) BaseWriter() io.Writer {
	return wc.base
}
