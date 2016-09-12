package bearychat

import (
	"net/http"
	"testing"

	"github.com/bcho/bearychat.go/mock"
)

const (
	testWebhook = "http://127.0.0.1:3927/=bwaaa/incoming/deadbeef"
)

func TestWebhookResponse_IsOk(t *testing.T) {
	var resp WebhookResponse

	resp = WebhookResponse{Code: 0}
	if !resp.IsOk() {
		t.Errorf("response should be ok when code is 0")
	}

	resp = WebhookResponse{Code: 1}
	if resp.IsOk() {
		t.Errorf("response should not be ok when code is not 0")
	}
}

func TestIncomingWebhookClient_SetWebhook(t *testing.T) {
	h := NewIncomingWebhookClient("")
	if h.SetWebhook(testWebhook) == nil {
		t.Errorf("should return webhook client")
	}

	if h.Webhook != testWebhook {
		t.Errorf("should set webhook")
	}
}

func TestIncomingWebhookClient_SetHTTPClient(t *testing.T) {
	h := NewIncomingWebhookClient(testWebhook)

	if h.httpClient != http.DefaultClient {
		t.Errorf("should use `http.DefaultClient` by default")
	}

	testHTTPClient := &http.Client{}
	if h.SetHTTPClient(testHTTPClient) == nil {
		t.Errorf("should return webhook client")
	}

	if h.httpClient != testHTTPClient {
		t.Errorf("should set http client")
	}
}

func TestIncomingWebhookClient_Send_WithoutWebhook(t *testing.T) {
	h := NewIncomingWebhookClient("")
	_, err := h.Send(nil)
	if err == nil {
		t.Errorf("should not send when webhook is not set")
	}
}

func TestIncomingWebhookClient_Send_WithoutHTTPClient(t *testing.T) {
	h := NewIncomingWebhookClient(testWebhook)
	h.SetHTTPClient(nil)
	_, err := h.Send(nil)
	if err == nil {
		t.Errorf("should not send when http client is not set")
	}
}

func TestIncomingWebhookClient_Send_Ok(t *testing.T) {
	s := mock.NewIncomingServer(
		mock.IncomingWithOkResponse,
		mock.IncomingWithWebhook(testWebhook),
	)
	go s.ListenAndServe()

	m := Incoming{Text: "Hello, BearyChat"}
	payload, err := m.Build()
	if err != nil {
		t.Errorf("build failed: %+v", err)
	}

	resp, err := NewIncomingWebhookClient(testWebhook).Send(payload)
	if err != nil {
		t.Errorf("send failed: %+v", err)
	}

	if !resp.IsOk() {
		t.Errorf("expect send ok")
	}
	s.Stop()
}

func TestIncomingWebhookClient_Send_Failed(t *testing.T) {
	expectedError := "failed"
	s := mock.NewIncomingServer(
		mock.IncomingWithErrorResponse(expectedError),
		mock.IncomingWithWebhook(testWebhook),
	)
	go s.ListenAndServe()

	m := Incoming{Text: "Hello, BearyChat"}
	payload, err := m.Build()
	if err != nil {
		t.Errorf("build failed: %+v", err)
	}

	resp, err := NewIncomingWebhookClient(testWebhook).Send(payload)
	if err != nil {
		t.Errorf("send failed: %+v", err)
	}

	if resp.IsOk() {
		t.Errorf("expect send failed")
	}
	if resp.Error != expectedError {
		t.Errorf("unexpected error: %s", resp.Error)
	}
	s.Stop()
}

func ExampleNewIncomingWebhookClient() {
	m := Incoming{Text: "Hello, BearyChat"}
	payload, _ := m.Build()
	resp, _ := NewIncomingWebhookClient("YOUR WEBHOOK URL").Send(payload)
	if resp.IsOk() {
		// parse resp result
	} else {
		// parse resp error
	}
}
