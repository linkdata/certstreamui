package certstreamui

import (
	"github.com/linkdata/certstream"
	"github.com/linkdata/jaws"
)

type uiLogEntries struct {
	*CertStreamUI
}

// JawsContains implements jaws.Container.
func (ui uiLogEntries) JawsContains(e *jaws.Element) (contents []jaws.UI) {
	ui.mu.RLock()
	defer ui.mu.RUnlock()
	ui.ring.Do(func(a any) {
		if le, ok := a.(*certstream.LogEntry); ok {
			contents = append(contents, uiLogEntry{LogEntry: le})
		}
	})
	return
}

func (csui *CertStreamUI) UiLogEntries() jaws.Container {
	return uiLogEntries{csui}
}
