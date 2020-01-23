package io

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

type readCloser struct {
	Reader io.Reader
	Closer io.Closer
}

func (r readCloser) Close() error {
	return r.Closer.Close()
}

func (r readCloser) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func OptionalGZIPReader(r io.Reader) io.ReadCloser {
	var rc readCloser
	rc.Closer = ioutil.NopCloser(r)
	if c, ok := r.(io.Closer); ok {
		rc.Closer = c
	}
	buf := new(bytes.Buffer)
	gz, err := gzip.NewReader(io.TeeReader(r, buf))
	if err != nil {
		rc.Reader = io.MultiReader(buf, r)
		return rc
	}
	rc.Reader = gz
	rc.Closer = MultiCloser(gz, rc.Closer)
	return rc
}

func ResponseReader(resp *http.Response) io.ReadCloser {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		buf := new(bytes.Buffer)
		gz, err := gzip.NewReader(io.TeeReader(resp.Body, buf))
		if err != nil {
			return readCloser{
				Reader: io.MultiReader(buf, resp.Body),
				Closer: resp.Body,
			}
		}
		return readCloser{
			Reader: gz,
			Closer: MultiCloser(gz, resp.Body),
		}
	}
	return readCloser{
		Reader: resp.Body,
		Closer: resp.Body,
	}
}
