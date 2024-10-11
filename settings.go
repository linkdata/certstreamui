package certstreamui

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/linkdata/deadlock"
)

type Settings struct {
	Filename      string
	mu            deadlock.RWMutex // protects following
	CertStreamURL string
}

func (s *Settings) Load() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CertStreamURL = "http://localhost:8081"
	var b []byte
	if b, err = os.ReadFile(s.Filename); err == nil {
		err = json.Unmarshal(b, s)
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	return
}

func (s *Settings) Save() (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var b []byte
	if b, err = json.MarshalIndent(s, "", "  "); err == nil {
		err = os.WriteFile(s.Filename, b, 0640)
	}
	return
}
