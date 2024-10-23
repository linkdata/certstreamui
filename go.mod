module github.com/linkdata/certstreamui

go 1.23.2

require (
	github.com/jackc/pgx/v5 v5.7.1
	github.com/linkdata/certstream v0.0.5
	github.com/linkdata/deadlock v0.5.2
	github.com/linkdata/jaws v0.79.0
	github.com/linkdata/webserv v0.9.7
)

replace github.com/linkdata/certstream => ../certstream

require (
	github.com/coder/websocket v1.8.12 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/google/certificate-transparency-go v1.2.1 // indirect
	github.com/google/trillian v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
)
