package main

import "net/http"

type Middleware []http.Handler

func (m *Middleware) Add(h http.Handler) {
	*m = append(*m, h)
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

func NewMiddlewareResponseWriter(w http.ResponseWriter) *MiddlewareResponseWriter {
	return &MiddlewareResponseWriter{
		ResponseWriter: w,
	}
}

func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(bytes)
}

func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	return w.ResponseWriter.WriteHeader(code)
}
