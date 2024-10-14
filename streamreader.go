package certstreamui

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
)

type CertStreamReader struct {
	domainCh chan string
}

func NewCertStreamReader(ctx context.Context, wsUrl string) (csr *CertStreamReader, err error) {
	var wsConn *websocket.Conn
	var resp *http.Response
	if wsConn, resp, err = websocket.Dial(ctx, wsUrl, nil); err == nil {
		csr = &CertStreamReader{
			domainCh: make(chan string),
		}
		_ = wsConn
		_ = resp
	}
	return
}

func (csr *CertStreamReader) consume() {
	defer close(csr.domainCh)

}
