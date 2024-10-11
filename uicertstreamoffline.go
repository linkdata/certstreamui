package certstreamui

import (
	"html/template"

	"github.com/linkdata/jaws"
)

type uiCertStreamOffline struct {
	*CertStreamUI
}

func (ui uiCertStreamOffline) JawsGetHtml(e *jaws.Element) template.HTML {
	var s string
	if s == "" {
		e.SetAttr("title", s)
		return "&#9888;&nbsp;"
	}
	return ""
}

func (csui *CertStreamUI) UiCertStreamOffline() jaws.HtmlGetter {
	return uiCertStreamOffline{csui}
}
