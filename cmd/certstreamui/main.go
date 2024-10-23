package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	_ "net/http/pprof"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/linkdata/certstream/certdb"
	"github.com/linkdata/certstreamui"
	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/webserv"
)

func env(key, dflt string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = dflt
	}
	return os.ExpandEnv(val)
}

var (
	flagAddress    = flag.String("address", env("WEBSERV_LISTEN", ""), "serve HTTP requests on given [address][:port]")
	flagCertDir    = flag.String("certdir", env("WEBSERV_CERTDIR", ""), "where to find fullchain.pem and privkey.pem")
	flagUser       = flag.String("user", env("WEBSERV_USER", ""), "switch to this user after startup (*nix only)")
	flagDataDir    = flag.String("datadir", env("WEBSERV_DATADIR", ""), "where to store data files after startup")
	flagListenURL  = flag.String("listenurl", env("WEBSERV_LISTENURL", ""), "manually specify URL where clients can reach us")
	flagVersion    = flag.Bool("v", false, "display version")
	flagCpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flagMemprofile = flag.String("memprofile", "", "write memory profile to file")
	flagPprof      = flag.Bool("pprof", false, "run pprof on http://localhost:6060/debug/pprof/")
	flagDbUser     = flag.String("dbuser", env("DBUSER", "certstream"), "database user")
	flagDbPass     = flag.String("dbpass", env("DBPASS", "certstream"), "database password")
	flagDbName     = flag.String("dbname", env("DBNAME", "certstream"), "database name")
	flagDbAddr     = flag.String("dbaddr", env("DBADDR", ""), "database address")
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
			if err = pprof.StartCPUProfile(f); err == nil {
				defer pprof.StopCPUProfile()
			}
		}
	}

	if *flagPprof {
		go func() {
			slog.Error(http.ListenAndServe("localhost:6060", nil).Error())
		}()
		slog.Info("pprof listening on http://localhost:6060/debug/pprof/")
	}

	if *flagDbAddr != "" {
		dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", *flagDbUser, *flagDbPass, *flagDbAddr, *flagDbName)
		db, err := sql.Open("pgx", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				var cdb *certdb.Certdb
				if cdb, err = certdb.New(context.Background(), db); err == nil {
					cdb.Close()
				}
			}
			db.Close()
		}
		if err != nil {
			fmt.Println(err)
		}
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
					var f *os.File
					if f, err = os.Create(*flagMemprofile); err == nil {
						defer f.Close()
						runtime.GC()
						_ = pprof.WriteHeapProfile(f)
					}
				}
				return
			}
		}
	}
	slog.Error(err.Error())
}
