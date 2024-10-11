package certstreamui

import (
	"github.com/linkdata/jaws"
)

type uiCertStreamURL struct {
	*Settings
	jaws.UiString
}

func (ui uiCertStreamURL) JawsSetString(e *jaws.Element, val string) (err error) {
	if err = ui.UiString.JawsSetString(e, val); err == nil {
		err = ui.Save()
	}
	return
}

func (s *Settings) UiCertStreamURL() jaws.StringSetter {
	return uiCertStreamURL{
		Settings: s,
		UiString: jaws.UiString{
			L: &s.mu,
			P: &s.CertStreamURL,
		}}
}
