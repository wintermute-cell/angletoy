package saving

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *SaveState
	once     sync.Once
)

// SaveState struct now public to allow external access
type SaveState struct {
	Version             int32   `json:"version"`
	TutorialSeen        bool    `json:"tutorial"`
	UnlockedRewardIndex int32   `json:"ridx"`
	CurrentDialogIndex  int32   `json:"didx"`
	IceMeterValue       float32 `json:"ice"`
}

// GetInstance returns the singleton instance of SaveState
func GetInstance() *SaveState {
	once.Do(func() {
		instance = &SaveState{
			Version:             1,
			TutorialSeen:        false,
			UnlockedRewardIndex: -1,
			CurrentDialogIndex:  0,
			IceMeterValue:       100,
		}
	})
	return instance
}

// SaveGame saves the current state to a file
func SaveGame() error {
	state := GetInstance()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	saveDir := filepath.Join(homeDir, ".snowbound")
	os.MkdirAll(saveDir, os.ModePerm)

	savePath := filepath.Join(saveDir, "savefile.json")

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(state)
}

// SaveExists checks if the save file exists
func SaveExists() bool {
	homeDir, _ := os.UserHomeDir()
	savePath := filepath.Join(homeDir, ".snowbound", "savefile.json")

	_, err := os.Stat(savePath)
	return err == nil
}

// LoadGame loads the game state from a file
func LoadGame() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	savePath := filepath.Join(homeDir, ".snowbound", "savefile.json")

	file, err := os.Open(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var state SaveState
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return err
	}

	*GetInstance() = state // Update the singleton instance
	return nil
}
