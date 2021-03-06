package crawler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {

	// context and fetch information
	URL         string
	ContextStr  string
	ContextTime time.Time
	FetchedTime time.Time

	// response meta information
	Status               string
	StatusCode           int
	Proto                string
	ContentLength        int64
	TransferEncodingJson []byte
	HeaderJson           []byte
	TrailerJson          []byte
	RequestJson          []byte
	TlsJson              []byte

	// response body
	Body []byte
}

func (r *Response) Equal(r2 Response) bool {
	return r.URL == r2.URL &&
		r.ContextStr == r2.ContextStr &&
		r.ContextTime.Unix() == r2.ContextTime.Unix() &&
		r.FetchedTime.Unix() == r2.FetchedTime.Unix() &&
		r.Status == r2.Status &&
		r.StatusCode == r2.StatusCode &&
		r.Proto == r2.Proto &&
		r.ContentLength == r2.ContentLength &&
		string(r.TransferEncodingJson) == string(r2.TransferEncodingJson) &&
		string(r.HeaderJson) == string(r2.HeaderJson) &&
		string(r.TrailerJson) == string(r2.TrailerJson) &&
		string(r.RequestJson) == string(r2.RequestJson) &&
		string(r.TlsJson) == string(r2.TlsJson) &&
		string(r.Body) == string(r2.Body)
}

func (r *Response) ReadRaw(raw *http.Response) (err error) {

	// read basic information
	r.Status = raw.Status
	r.StatusCode = raw.StatusCode
	r.Proto = raw.Proto
	r.ContentLength = raw.ContentLength

	// read complex fields to JSON
	r.TransferEncodingJson, err = json.Marshal(raw.TransferEncoding)
	if err != nil {
		return
	}
	r.HeaderJson, err = json.Marshal(raw.Header)
	if err != nil {
		return
	}
	r.TrailerJson, err = json.Marshal(raw.Trailer)
	if err != nil {
		return
	}
	r.RequestJson, err = json.Marshal(*raw.Request)
	if err != nil {
		return
	}
	/*
		// no such field or method in Golang 1.2
		r.TlsJson, err = json.Marshal(raw.TLS)
		if err != nil {
			return
		}
	*/

	// read body
	defer raw.Body.Close()
	r.Body, err = ioutil.ReadAll(raw.Body)
	return
}

func (r *Response) InContext(ctx *Context) bool {
	return r.ContextStr == ctx.Str &&
		r.ContextTime.Equal(ctx.Time) &&
		r.FetchedTime.Equal(ctx.Fetched)
}

func (r *Response) SetContext(ctx *Context) {
	r.ContextStr = ctx.Str
	r.ContextTime = ctx.Time
	r.FetchedTime = ctx.Fetched
}

func (r *Response) GetContext() (ctx *Context) {
	return &Context{
		Str:     r.ContextStr,
		Time:    r.ContextTime,
		Fetched: r.FetchedTime,
	}
}
