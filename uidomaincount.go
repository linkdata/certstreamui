package certstreamui

import (
	"html/template"
	"strconv"
	"sync/atomic"

	"github.com/linkdata/jaws"
)

type uiDomainCount struct {
	*CertStreamUI
}

func (ui uiDomainCount) JawsGetHtml(e *jaws.Element) template.HTML {
	return template.HTML(strconv.FormatUint(atomic.LoadUint64(&ui.DomainCount), 10))
}

func (csui *CertStreamUI) UiDomainCount() jaws.HtmlGetter {
	return uiDomainCount{csui}
}
