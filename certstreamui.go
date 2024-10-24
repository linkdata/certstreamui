package certstreamui

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"sync/atomic"
	"time"

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
	cs          *certstream.CertStream
}

func New(cfg *webserv.Config, mux *http.ServeMux, jw *jaws.Jaws) (csui *CertStreamUI, err error) {
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
				}
			}
		}
	}
	return
}

func (csui *CertStreamUI) CertStream() (cs *certstream.CertStream) {
	csui.mu.RLock()
	cs = csui.cs
	csui.mu.RUnlock()
	return
}

func (csui *CertStreamUI) Run(ctx context.Context) {
	destCh := make(chan *certstream.LogEntry, 256)
	defer close(destCh)
	go csui.process(ctx, destCh)
	for {
		csui.readLogEntries(ctx, destCh)
		if ctx.Err() != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func (csui *CertStreamUI) process(_ context.Context, destCh <-chan *certstream.LogEntry) {
	when := time.Now()
	for le := range destCh {
		count := int64(len(le.DNSNames()))
		atomic.AddUint64(&csui.DomainCount, uint64(count))
		if since := time.Since(when); since >= time.Second {
			when = when.Add(since)
			csui.Jaws.Dirty(csui.UiLogEntries())
		}
	}
}

func (csui *CertStreamUI) readLogEntries(ctx context.Context, destCh chan<- *certstream.LogEntry) {
	cs := certstream.New()
	cs.BatchSize = 32
	cs.ParallelFetch = 4
	ctx, cancel := context.WithTimeout(ctx, time.Hour*24)
	defer cancel()
	if entryCh, err := cs.Start(ctx, nil); err == nil {
		var operators []string
		for opdom, op := range cs.Operators {
			operators = append(operators, fmt.Sprintf("%s*%d", opdom, len(op.Streams)))
		}
		slices.Sort(operators)
		slog.Info("certstream", "operators", operators)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		csui.mu.Lock()
		csui.cs = cs
		csui.mu.Unlock()
		for {
			select {
			case <-ticker.C:
				csui.mu.RLock()
				_, stopped := cs.CountStreams()
				csui.mu.RUnlock()
				if stopped > 1 {
					return
				}
			case le, ok := <-entryCh:
				if !ok {
					return
				}
				destCh <- le
			}
		}
	} else {
		slog.Error("certstream.Start()", "err", err)
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
