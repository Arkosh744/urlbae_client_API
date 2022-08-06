package urlbae

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type loggingRoundTripper struct {
	logger io.Writer
	next   http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Fprintf(l.logger, "[%s] %s %s \n", time.Now().Format(time.ANSIC), req.Method, req.URL.String())
	return l.next.RoundTrip(req)
}
