package certstreamui

import (
	"fmt"
	"html/template"
	"net/http"
)

func (csui *CertStreamUI) UiNavLink(rq *http.Request, url, title string) template.HTML {
	if rq != nil && rq.URL.Path == url {
		return template.HTML(fmt.Sprintf(`<a href="#" class="nav-link active">%s</a>`, title)) // #nosec G203
	}
	return template.HTML(fmt.Sprintf(`<a href="%s" class="nav-link">%s</a>`, url, title)) // #nosec G203
}
