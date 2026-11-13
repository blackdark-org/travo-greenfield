package connectivity

import (
	"context"
	"net/http"
	"time"
)

const ProbeURL = "http://connectivitycheck.gstatic.com/generate_204"

func NewHTTPProbe(url string, timeout time.Duration) Probe {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DisableKeepAlives = true
	client := &http.Client{
		Transport:     transport,
		Timeout:       timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return func(ctx context.Context) Internet {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return Offline
		}
		res, err := client.Do(req)
		if err != nil {
			return Offline
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusNoContent {
			return Online
		}
		if res.StatusCode >= 300 && res.StatusCode < 400 {
			target, err := res.Location()
			if err == nil && target.Hostname() != "" && (target.Scheme == "http" || target.Scheme == "https") {
				return Captive
			}
		}
		return Offline
	}
}
