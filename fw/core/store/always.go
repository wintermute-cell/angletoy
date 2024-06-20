package store

// This file defines some premade types of stored data that is used by the framework.

type AppState struct {
	ShouldQuit bool
}

func NewAppState() *AppState {
	return &AppState{}
}
