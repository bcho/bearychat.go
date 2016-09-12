package mock

import "testing"

func TestIncomingWithOkResponse(t *testing.T) {
	s := NewIncomingServer(IncomingWithOkResponse)
	if s.HandlerFunc == nil {
		t.Errorf("HandlerFunc should not be nil")
	}
}

func TestIncomingWithErrorResponse(t *testing.T) {
	s := NewIncomingServer(IncomingWithErrorResponse("foobar"))
	if s.HandlerFunc == nil {
		t.Errorf("HandlerFunc should not be nil")
	}
}

func TestIncomingWithWebhook(t *testing.T) {
	webhook := "http://127.0.0.1:3927/foobar"
	s := NewIncomingServer(IncomingWithWebhook(webhook))
	if s.Webhook.String() != webhook {
		t.Errorf("unexpected webhook %s", s.Webhook)
	}
	if s.Addr == nil {
		t.Errorf("Addr should not be nil")
	}
}
