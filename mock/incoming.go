package mock

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type IncomingServer struct {
	Addr        *url.URL // listen address
	Webhook     *url.URL // webhook url
	HandlerFunc http.HandlerFunc

	ln net.Listener
}

type incomingServerSetter func(*IncomingServer)

// Set successful response.
func IncomingWithOkResponse(s *IncomingServer) {
	s.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"code": 0, "result": null}`))
	}
}

// Set error response.
func IncomingWithErrorResponse(err string) incomingServerSetter {
	resp := fmt.Sprintf(`{"code": 1, "result": null, "error": "%s"}`, err)

	return func(s *IncomingServer) {
		s.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(resp))
		}
	}
}

// Set webhook.
//
// Panic when webhook is invalid url.
func IncomingWithWebhook(rawwebhook string) incomingServerSetter {
	webhook, err := url.Parse(rawwebhook)
	if err != nil {
		panic(err)
	}
	addr := &url.URL{
		Scheme: webhook.Scheme,
		Host:   webhook.Host,
	}

	return func(s *IncomingServer) {
		s.Addr = addr
		s.Webhook = webhook
	}
}

// IncomingServer mocks incoming webhook response.
//
//      go NewIncomingServer(
//              IncomingWithOkResponse,
//              IncomingWithWebhook("http://localhost:3927"),
//      ).ListenAndServe()
func NewIncomingServer(setters ...incomingServerSetter) *IncomingServer {
	server := &IncomingServer{}

	for _, setter := range setters {
		setter(server)
	}

	return server
}

// ListenAndServe starts accepting incoming connection.
func (s *IncomingServer) ListenAndServe() error {
	srv := &http.Server{Handler: s}
	srv.SetKeepAlivesEnabled(false)

	ln, err := net.Listen("tcp", s.Addr.Host)
	if err != nil {
		return err
	}
	s.ln = ln

	return srv.Serve(ln)
}

// Stop mock server.
func (s *IncomingServer) Stop() {
	if s.ln != nil {
		// TODO wait all requests done
		s.ln.Close()
	}

	time.Sleep(1 * time.Millisecond)
}

// ServeHTTP serves a http request.
func (s *IncomingServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != s.Webhook.Path {
		// invalid webhook address
		w.WriteHeader(404)
		w.Write([]byte(`{"code": 6, "error": "资源未找到", "result": null}`))
		return
	}

	if s.HandlerFunc != nil {
		s.HandlerFunc(w, req)
	}
}
