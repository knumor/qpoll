package handlers

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// MyWriteCloser is a custom writer that implements the io.WriteCloser interface
type MyWriteCloser struct {
    *bufio.Writer
}

// Close flushes the buffer
func (mwc *MyWriteCloser) Close() error {
    return mwc.Flush()
}

// GenQRForPoll generates a QR code as image data
func (hc *HandlerContext) GenQRForPoll(rw http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	qr, err := qrcode.New(fmt.Sprintf("http://10.0.0.33:8080/vote/%s", id))
	if err != nil {
		http.Error(rw, "failed to create QR code", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "image/png")

	bw := bufio.NewWriter(rw)
	mwc := &MyWriteCloser{bw}
	qrw := standard.NewWithWriter(mwc,
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithBgTransparent(),
	)

	qr.Save(qrw)
}
