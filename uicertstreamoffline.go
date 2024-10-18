package certstreamui

import (
	"fmt"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiCertStreamOffline struct {
	*CertStreamUI
}

func (ui uiCertStreamOffline) JawsGetHtml(e *jaws.Element) template.HTML {
	ui.mu.RLock()
	running, stopped := ui.running, ui.stopped
	ui.mu.RUnlock()
	var warn string
	if stopped > 0 {
		warn = "&#9888;&nbsp;"
	}
	return template.HTML(fmt.Sprintf("%s%d/%d&nbsp;", warn, running, running+stopped))
}

func (csui *CertStreamUI) UiCertStreamOffline() jaws.HtmlGetter {
	return uiCertStreamOffline{csui}
}
