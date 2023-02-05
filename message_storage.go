package main

import "errors"

type MessageStorage interface {
	Store(message string) string
	Exists(id string) bool
	Redeem(id string) (string, error)
}

type inMemoryMessageStorage struct {
	data map[string]string
}

// stores message in a map based on ID
// ID is returned for future message redemption
func (m *inMemoryMessageStorage) Store(message string) string {
	id := randomID()
	m.data[id] = message
	return id
}

func (m *inMemoryMessageStorage) Exists(id string) bool {
	_, exists := m.data[id]
	return exists
}

// If ID exists and message is reedemed, we need to delete the message from the map
func (m *inMemoryMessageStorage) Redeem(id string) (string, error) {
	message, exists := m.data[id]
	if exists {
		delete(m.data, id)
		return message, nil
	} else {
		return "", errors.New("ID is not valid.")
	}
}

func NewInMemoryMessageStorage() MessageStorage {
	data := make(map[string]string)
	p := inMemoryMessageStorage{data}
	return &p
}
