package certstreamui

import (
	"encoding/json"
	"errors"
	"net/url"
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

func (s *Settings) saveLocked() (err error) {
	var b []byte
	if b, err = json.MarshalIndent(s, "", "  "); err == nil {
		err = os.WriteFile(s.Filename, b, 0640)
	}
	return
}

func (s *Settings) Save() (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.saveLocked()
}

func (s *Settings) SetCertStreamURL(val string) (err error) {
	var u *url.URL
	if u, err = url.Parse(val); err == nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.CertStreamURL = u.String()
		err = s.saveLocked()
	}
	return
}

func (s *Settings) GetCertStreamURL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CertStreamURL
}
