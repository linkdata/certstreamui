package certstreamui

import (
	"strings"

	"github.com/linkdata/jaws"
)

type uiCertStreamURL struct {
	*CertStreamUI
	jaws.String
}

func (ui *uiCertStreamURL) JawsClick(e *jaws.Element, name string) (err error) {
	urlStr := strings.TrimSpace(ui.String.Get())
	if err = ui.Settings.SetCertStreamURL(urlStr); err == nil {
		// establish connection
	}
	return
}

func (csui *CertStreamUI) UiCertStreamURL() jaws.ClickHandler {
	return &uiCertStreamURL{
		CertStreamUI: csui,
		String:       jaws.String{Value: csui.Settings.GetCertStreamURL()},
	}
}
