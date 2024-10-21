package certstreamui

import (
	"github.com/linkdata/jaws"
)

type uiLogEntries struct {
	*CertStreamUI
}

// JawsContains implements jaws.Container.
func (ui uiLogEntries) JawsContains(e *jaws.Element) (contents []jaws.UI) {
	ui.mu.RLock()
	defer ui.mu.RUnlock()
	return
}

func (csui *CertStreamUI) UiLogEntries() jaws.Container {
	return uiLogEntries{csui}
}
