package certstreamui

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"

	"embed"

	"github.com/linkdata/certstream"
	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/staticserve"
	"github.com/linkdata/webserv"
)

//go:embed assets
var assetsFS embed.FS

//go:generate go run github.com/cparta/makeversion/v2/cmd/mkver@latest -name CertStreamUI -out version.gen.go

type CertStreamUI struct {
	Config      *webserv.Config
	Jaws        *jaws.Jaws
	FaviconURI  string
	PkgName     string
	PkgVersion  string
	Settings    Settings
	DomainCount uint64
	mu          deadlock.RWMutex // protects following
	domainCh    <-chan string
}

func New(cfg *webserv.Config, mux *http.ServeMux, jw *jaws.Jaws, entryCh <-chan *certstream.LogEntry) (csui *CertStreamUI, err error) {
	var tmpl *template.Template
	var faviconuri string
	if err = os.MkdirAll(cfg.DataDir, 0750); err == nil {
		if tmpl, err = template.New("").ParseFS(assetsFS, "assets/ui/*.html"); err == nil {
			jw.AddTemplateLookuper(tmpl)
			var extraFiles []string
			addStaticFiles := func(filename string, ss *staticserve.StaticServe) (err error) {
				uri := path.Join("/static", ss.Name)
				if strings.HasSuffix(filename, "favicon.png") {
					faviconuri = uri
				}
				extraFiles = append(extraFiles, uri)
				mux.Handle(uri, ss)
				return
			}
			if err = staticserve.WalkDir(assetsFS, "assets/static", addStaticFiles); err == nil {
				if err = jw.GenerateHeadHTML(extraFiles...); err == nil {
					csui = &CertStreamUI{
						Config:     cfg,
						Jaws:       jw,
						FaviconURI: faviconuri,
						PkgName:    PkgName,
						PkgVersion: PkgVersion,
					}
					csui.AddRoutes(mux)
					csui.Settings.filename = path.Join(csui.Config.DataDir, "settings.json")
					err = csui.Settings.Load()
					go csui.readLogEntries(entryCh)
				}
			}
		}
	}
	return
}

func (csui *CertStreamUI) readLogEntries(entryCh <-chan *certstream.LogEntry) {
	for le := range entryCh {
		atomic.AddUint64(&csui.DomainCount, uint64(len(le.DNSNames())))
	}
}

func (csui *CertStreamUI) AddRoutes(mux *http.ServeMux) {
	mux.Handle("GET /{$}", csui.Jaws.Handler("index.html", csui))
	mux.Handle("GET /setup/{$}", csui.Jaws.Handler("setup.html", csui))
	mux.Handle("GET /about/{$}", csui.Jaws.Handler("about.html", csui))
}

func (csui *CertStreamUI) SettingsLoad() (err error) {
	return csui.Settings.Load()
}

func (csui *CertStreamUI) SettingsSave() (err error) {
	return csui.Settings.Save()
}

func (csui *CertStreamUI) Close() (err error) {
	return
}
