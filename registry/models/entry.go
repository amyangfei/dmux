package models

import (
	"encoding/json"
	"github.com/amyangfei/dmux/store"
)

type EntryManager struct {
	Storage store.Storage
}

// Entry represents an entry to log file.
type Entry struct {
	// Tag is used for log file identification
	Tag string
	// Path is the fullpath to find log files
	Path string
	// If Include is provided, only files provided in Include will be handled
	Include []string
	// Any log file that matches one regex item in Exclude will be ignored
	Exclude []string
}

// Get returns the Entry with the given id in the EntryManager
func (m *EntryManager) Get(id string) (*Entry, error) {
	data, err := m.Storage.Get(id)
	if err != nil {
		return nil, err
	} else if data != nil {
		entry := &Entry{}
		if err := json.Unmarshal(data, &entry); err != nil {
			return nil, err
		}
		return entry, nil
	} else {
		return nil, nil
	}
}

// Delete deletes the Entry with the given id in the EntryManager
func (m *EntryManager) Delete(id string) error {
	return m.Storage.Del(id)
}

// All returns the list of all registried Entries.
func (m *EntryManager) All() ([]*Entry, error) {
	entries := make([]*Entry, 0)
	raw, err := m.Storage.List(EntryPrefix)
	if err != nil {
		return entries, err
	}
	for _, data := range raw {
		var entry *Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Save saves the given entry in the EntryManager.
func (m *EntryManager) Save(entry *Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	key := EntryPrefix + entry.Path
	return m.Storage.Set(key, data)
}
