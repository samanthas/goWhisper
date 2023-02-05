package main

import "testing"

func Test_NewInMemoryMessageStorage(t *testing.T) {
	storage := NewInMemoryMessageStorage()

	if storage == nil {
		t.Fatal("expected instance to not be nil")
	}
}

func Test_InMemoryMessageStore_Store(t *testing.T) {
	storage := NewInMemoryMessageStorage()

	id := storage.Store("my message")
	if id == "" {
		t.Fatal("Expected id to not be empty")
	}

	exists := storage.Exists(id)
	if !exists {
		t.Fatalf("Expected storage to contain a message with ID=%v", id)
	}
}

func Test_InMemoryMessageStorage_Store_Multiple(t *testing.T) {
	storage := NewInMemoryMessageStorage()

	id1 := storage.Store("message1")
	id2 := storage.Store("message2")

	if id1 == id2 {
		t.Fatal("Expected message ids to be different, but were not")
	}
}

func Test_InMemoryMessageStore_Exists(t *testing.T) {
	storage := NewInMemoryMessageStorage()

	id := storage.Store("my message")

	exists := storage.Exists(id)
	if !exists {
		t.Fatalf("Expected storage to contain a message with ID=%v", id)
	}
}

func Test_InMemoryMessageStore_Exists_Random(t *testing.T) {
	storage := NewInMemoryMessageStorage()

	id := "thisShouldNotExist"

	exists := storage.Exists(id)
	if exists {
		t.Fatalf("Expected message ID=%v to not exist in storage", id)
	}
}

func Test_InMemoryMessageStorage_Redeem(t *testing.T) {
	storage := NewInMemoryMessageStorage()

	// Message with non-extant ID should return error
	_, err := storage.Redeem("someNonExtantID")
	if err == nil {
		t.Fatal("Expected redeem to return error, but did not")
	}

	message := "my message"
	id := storage.Store(message)
	redeemedMessage, err := storage.Redeem(id)
	if err != nil {
		t.Fatalf("Expected redeem to not return error, got error %v", err)
	}

	if redeemedMessage != message {
		t.Fatalf("Expected %v to equal %v", redeemedMessage, message)
	}

	// Message should no longer exist after being redeemed
	if storage.Exists(id) {
		t.Fatalf("Expected message ID=%v to no longer exist", id)
	}

	redeemedMessage, err = storage.Redeem(id)
	if err == nil {
		t.Fatalf("Expected redeem ID=%v to return error, but did not", id)
	}
}
