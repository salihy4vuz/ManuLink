package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// DBStore represents our persistent data structure
type DBStore struct {
	Accounts map[string]interface{} `json:"accounts"` // Used interface{} to prevent circular dependency, will map via JSON
	Tickets  []interface{}          `json:"tickets"`
}

var (
	dbFilePath = "database/data.json"
	dbMutex    sync.Mutex
)

// Save writes data to the JSON file
func Save(accounts interface{}, tickets interface{}) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	store := DBStore{
		Accounts: accounts.(map[string]interface{}),
		Tickets:  tickets.([]interface{}),
	}

	bytes, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dbFilePath, bytes, 0644)
}

// Load reads data from the JSON file
func Load(accountsTarget interface{}, ticketsTarget interface{}) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return nil // No db file yet, use defaults
	}

	bytes, err := ioutil.ReadFile(dbFilePath)
	if err != nil {
		return err
	}

	// We decode the entire file into a temporary struct, then re-encode/decode parts to the targets
	var tempStore map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &tempStore); err != nil {
		return err
	}

	if accountsTarget != nil && tempStore["accounts"] != nil {
		if err := json.Unmarshal(tempStore["accounts"], accountsTarget); err != nil {
			return err
		}
	}

	if ticketsTarget != nil && tempStore["tickets"] != nil {
		if err := json.Unmarshal(tempStore["tickets"], ticketsTarget); err != nil {
			return err
		}
	}

	return nil
}

// EnsureDB calls Save immediately to ensure file exists
func EnsureDB(accounts interface{}, tickets interface{}) {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		_ = Save(accounts, tickets)
	}
}
