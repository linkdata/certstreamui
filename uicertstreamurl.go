package certstreamui

import (
	"github.com/linkdata/jaws"
)

type uiCertStreamURL struct {
	*CertStreamUI
	jaws.String
}

func (ui *uiCertStreamURL) JawsClick(e *jaws.Element, name string) (err error) {
	return
}

func (csui *CertStreamUI) UiCertStreamURL() jaws.ClickHandler {
	return &uiCertStreamURL{
		CertStreamUI: csui,
		/*String:       jaws.String{Value: csui.Settings.GetCertStreamURL()},*/
	}
}
