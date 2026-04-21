package sse_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/ajianaz/gofin-full/api/internal/sse"
)

func TestHubSubscribeUnsubscribe(t *testing.T) {
	log := zerolog.Nop()
	hub := sse.NewHub(log)

	client := &sse.Client{
		ID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID: uuid.MustParse("00000000-0000-0000-0000-00000000000a"),
		Ch:     make(chan sse.Event, 16),
		Done:   make(chan struct{}),
	}

	hub.Subscribe(client)
	if got := hub.ClientCount(); got != 1 {
		t.Errorf("expected 1 client, got %d", got)
	}

	hub.Unsubscribe(client)
	if got := hub.ClientCount(); got != 0 {
		t.Errorf("expected 0 clients, got %d", got)
	}
}

func TestHubSendToUser(t *testing.T) {
	log := zerolog.Nop()
	hub := sse.NewHub(log)

	testUserID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	client := &sse.Client{
		ID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		UserID: testUserID,
		Ch:     make(chan sse.Event, 16),
		Done:   make(chan struct{}),
	}

	hub.Subscribe(client)
	defer hub.Unsubscribe(client)

	event := sse.Event{Type: "test", Data: map[string]string{"msg": "hello"}}
	hub.SendToUser(testUserID, event)

	select {
	case received := <-client.Ch:
		if received.Type != "test" {
			t.Errorf("expected event type test, got %s", received.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestHubBroadcast(t *testing.T) {
	log := zerolog.Nop()
	hub := sse.NewHub(log)

	userID1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID2 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	c1 := &sse.Client{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), UserID: userID1, Ch: make(chan sse.Event, 16), Done: make(chan struct{})}
	c2 := &sse.Client{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), UserID: userID2, Ch: make(chan sse.Event, 16), Done: make(chan struct{})}

	hub.Subscribe(c1)
	hub.Subscribe(c2)
	defer func() {
		hub.Unsubscribe(c1)
		hub.Unsubscribe(c2)
	}()

	event := sse.Event{Type: "broadcast", Data: "hello all"}
	hub.Broadcast(event)

	for i, c := range []*sse.Client{c1, c2} {
		select {
		case received := <-c.Ch:
			if received.Type != "broadcast" {
				t.Errorf("client %d: expected broadcast, got %s", i, received.Type)
			}
		case <-time.After(time.Second):
			t.Fatalf("client %d: timed out", i)
		}
	}
}

func TestHubSendToNonexistentUser(t *testing.T) {
	log := zerolog.Nop()
	hub := sse.NewHub(log)

	done := make(chan struct{})
	go func() {
		hub.SendToUser(uuid.MustParse("00000000-0000-0000-0000-ffffffffffff"), sse.Event{Type: "test", Data: "noop"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SendToUser blocked on non-existent user")
	}
}

func TestHubConcurrentAccess(t *testing.T) {
	log := zerolog.Nop()
	hub := sse.NewHub(log)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			clientID := uuid.New()
			userID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
			userID[15] = byte(id % 10)
			client := &sse.Client{ID: clientID, UserID: userID, Ch: make(chan sse.Event, 16), Done: make(chan struct{})}
			hub.Subscribe(client)
			hub.SendToUser(userID, sse.Event{Type: "test", Data: id})
			hub.Unsubscribe(client)
		}(i)
	}
	wg.Wait()

	if got := hub.ClientCount(); got != 0 {
		t.Errorf("expected 0 clients after cleanup, got %d", got)
	}
}

func TestMarshalEvent(t *testing.T) {
	event := sse.Event{Type: "notification", Data: map[string]string{"title": "test"}}
	data := sse.MarshalEvent(event)

	if string(data) != `event: notification`+"\n"+`data: {"title":"test"}`+"\n\n" {
		t.Errorf("unexpected marshal output: %q", string(data))
	}
}
