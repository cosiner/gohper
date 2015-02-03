package filters

import (
	"compress/gzip"
	"strings"

	. "github.com/cosiner/gomodule/server"
)

type (
	GzipResponseWrapper struct {
		gzipWriter *gzip.Writer
		Response
	}
)

func (grw GzipResponseWrapper) Write(data []byte) (int, error) {
	return grw.gzipWriter.Write(data)
}

func GzipFilter(req Request, resp Response, chain FilterChain) {
	if enc := req.ContentEncoding(); strings.Contains(enc, ENCODING_GZIP) {
		resp.SetContentEncoding(ENCODING_GZIP)
		grw := GzipResponseWrapper{
			gzipWriter: gzip.NewWriter(resp),
			Response:   resp,
		}
		chain.Filter(req, grw)
		grw.gzipWriter.Close()
	} else {
		chain.Filter(req, resp)
	}
}
