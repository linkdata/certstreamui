package certstreamui

import (
	"html/template"
	"sync"
	"time"

	"github.com/linkdata/jaws"
)

type uiClock struct{}

var uiClockStart sync.Once

func (ui uiClock) JawsGetHtml(e *jaws.Element) template.HTML {
	uiClockStart.Do(func() {
		go func(jw *jaws.Jaws) {
			for {
				now := time.Now()
				time.Sleep(time.Second - now.Sub(now.Truncate(time.Second)))
				jw.Dirty(ui)
			}
		}(e.Jaws)
	})
	now := time.Now().Round(time.Second)
	fmt := "15:04&nbsp;MST"
	if (now.Second() % 2) == 0 {
		fmt = "15&nbsp;04&nbsp;MST"
	}
	return template.HTML(now.Format(fmt)) // #nosec G203
}

func (csui *CertStreamUI) UiClock() jaws.HtmlGetter {
	return uiClock{}
}
