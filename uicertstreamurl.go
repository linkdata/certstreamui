package certstreamui

import (
	"github.com/linkdata/jaws"
)

type uiCertStreamURL struct {
	*Settings
	jaws.String
}

func (ui *uiCertStreamURL) JawsClick(e *jaws.Element, name string) error {
	return ui.SetCertStreamURL(ui.String.Get())
}

func (s *Settings) UiCertStreamURL() jaws.ClickHandler {
	return &uiCertStreamURL{
		Settings: s,
		String:   jaws.String{Value: s.GetCertStreamURL()},
	}
}
