package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/linkdata/certstream"
	"github.com/linkdata/certstreamui"
	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/webserv"
)

var (
	flagAddress   = flag.String("address", os.Getenv("WEBSERV_LISTEN"), "serve HTTP requests on given [address][:port]")
	flagCertDir   = flag.String("certdir", os.Getenv("WEBSERV_CERTDIR"), "where to find fullchain.pem and privkey.pem")
	flagUser      = flag.String("user", os.Getenv("WEBSERV_USER"), "switch to this user after startup (*nix only)")
	flagDataDir   = flag.String("datadir", os.Getenv("WEBSERV_DATADIR"), "where to store data files after startup")
	flagListenURL = flag.String("listenurl", os.Getenv("WEBSERV_LISTENURL"), "manually specify URL where clients can reach us")
	flagVersion   = flag.Bool("v", false, "display version")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(certstreamui.PkgVersion)
		return
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
		var entryCh <-chan *certstream.LogEntry
		if entryCh, err = certstream.New().Start(ctx, nil); err == nil {
			var csui *certstreamui.CertStreamUI
			if csui, err = certstreamui.New(cfg, http.DefaultServeMux, jw, entryCh); err == nil {
				defer csui.Close()
				if err = cfg.Serve(ctx, l, http.DefaultServeMux); err == nil {
					return
				}
			}
		}
	}
	slog.Error(err.Error())
}
