package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	_ "net/http/pprof"

	"github.com/linkdata/certstreamui"
	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/webserv"
)

var (
	flagAddress    = flag.String("address", os.Getenv("WEBSERV_LISTEN"), "serve HTTP requests on given [address][:port]")
	flagCertDir    = flag.String("certdir", os.Getenv("WEBSERV_CERTDIR"), "where to find fullchain.pem and privkey.pem")
	flagUser       = flag.String("user", os.Getenv("WEBSERV_USER"), "switch to this user after startup (*nix only)")
	flagDataDir    = flag.String("datadir", os.Getenv("WEBSERV_DATADIR"), "where to store data files after startup")
	flagListenURL  = flag.String("listenurl", os.Getenv("WEBSERV_LISTENURL"), "manually specify URL where clients can reach us")
	flagVersion    = flag.Bool("v", false, "display version")
	flagCpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flagMemprofile = flag.String("memprofile", "", "write memory profile to file")
	flagPprof      = flag.Bool("pprof", false, "run pprof on http://localhost:6060/debug/pprof/")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(certstreamui.PkgVersion)
		return
	}

	if *flagCpuprofile != "" {
		if f, err := os.Create(*flagCpuprofile); err == nil {
			defer f.Close()
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	if *flagPprof {
		go func() {
			slog.Error(http.ListenAndServe("localhost:6060", nil).Error())
		}()
		slog.Info("pprof listening on http://localhost:6060/debug/pprof/")
	}

	cfg := &webserv.Config{
		Address:              *flagAddress,
		CertDir:              *flagCertDir,
		User:                 *flagUser,
		DataDir:              *flagDataDir,
		ListenURL:            *flagListenURL,
		DefaultDataDirSuffix: "certstreamui",
		Logger:               slog.Default(),
	}

	jw := jaws.New()
	defer jw.Close()
	jw.Debug = deadlock.Debug
	jw.Logger = slog.Default()
	http.DefaultServeMux.Handle("/jaws/", jw)
	go jw.Serve()

	l, err := cfg.Listen()
	if err == nil {
		defer l.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var csui *certstreamui.CertStreamUI
		if csui, err = certstreamui.New(cfg, http.DefaultServeMux, jw); err == nil {
			defer csui.Close()
			go csui.Run(ctx)
			if err = cfg.Serve(ctx, l, http.DefaultServeMux); err == nil {
				if *flagMemprofile != "" {
					if f, err := os.Create(*flagMemprofile); err == nil {
						defer f.Close()
						runtime.GC()
						err = pprof.WriteHeapProfile(f)
					}
				}
				return
			}
		}
	}
	slog.Error(err.Error())
}
