package certstreamui

import (
	"io"

	"github.com/linkdata/certstream"
	"github.com/linkdata/jaws"
)

type uiLogEntry struct {
	*certstream.LogEntry
}

// JawsRender implements jaws.UI.
func (ui uiLogEntry) JawsRender(e *jaws.Element, w io.Writer, params []any) error {
	return jaws.NewTemplate("tr_logentry.html", ui).JawsRender(e, w, params)
}

// JawsUpdate implements jaws.UI.
func (ui uiLogEntry) JawsUpdate(e *jaws.Element) {
}
