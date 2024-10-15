package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/linkdata/certstreamui"
	"github.com/linkdata/certstreamui/stream"
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

func testStream() {
	ch, err := stream.Stream(context.Background(), []stream.Operator{stream.ALL}, []stream.LogStatus{stream.Usable}, -1, 256, 1)
	if err != nil {
		slog.Error("e", "err", err)
		return
	}
	for b := range ch {
		fmt.Printf("%s %d\n", b.LogSourceName, len(b.Logs))
	}
}

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
		var csui *certstreamui.CertStreamUI
		if csui, err = certstreamui.New(cfg, http.DefaultServeMux, jw); err == nil {
			defer csui.Close()
			if err = cfg.Serve(context.Background(), l, http.DefaultServeMux); err == nil {
				return
			}
		}
	}
	slog.Error(err.Error())
}
