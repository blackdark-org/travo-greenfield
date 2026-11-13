package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/blackdark-org/travo-greenfield/internal/wifi"
)

type fakeWiFi struct {
	summary    wifi.Summary
	scan       wifi.ScanResult
	summaryErr error
	scanErr    error
}

func (f fakeWiFi) Summary(context.Context) (wifi.Summary, error) {
	return f.summary, f.summaryErr
}

func (f fakeWiFi) Scan(context.Context) (wifi.ScanResult, error) {
	return f.scan, f.scanErr
}

func wifiEnv(t *testing.T, service WiFiService) *env {
	t.Helper()
	return newEnvWithDashboard(t, nil, nil, nil, func(d *Deps) { d.WiFi = service })
}

func TestWiFiRoutesRequireSession(t *testing.T) {
	e := wifiEnv(t, fakeWiFi{})
	for _, path := range []string{"/api/v1/wifi", "/api/v1/wifi/scan"} {
		if got := e.do(req{method: http.MethodGet, path: path}); got.Code != http.StatusUnauthorized {
			t.Errorf("%s: %d %s", path, got.Code, got.Body)
		}
	}
}

func TestWiFiSummaryJSONContract(t *testing.T) {
	associated, up := false, true
	e := wifiEnv(t, fakeWiFi{summary: wifi.Summary{
		Radios: []wifi.RadioInfo{{Name: "radio0", Band: "5g", Up: true}},
		STAs:   []wifi.STA{{Section: "wwan", Radio: "radio0", SSID: "HomeUplink", Ifname: "phy0-sta0", Associated: &associated}},
		APs:    []wifi.AP{{Section: "default_radio1", Radio: "radio1", SSID: "TravelAP", Ifname: "phy1-ap0", Up: &up}, {Section: "unknown", Up: nil}},
	}})
	cookie, _ := e.login()
	got := e.do(req{method: http.MethodGet, path: "/api/v1/wifi", cookie: cookie})
	assertWiFiJSON(t, got, `{
		"radios":[{"name":"radio0","band":"5g","disabled":false,"up":true}],
		"stas":[{"section":"wwan","radio":"radio0","ssid":"HomeUplink","ifname":"phy0-sta0","associated":false}],
		"aps":[
			{"section":"default_radio1","radio":"radio1","ssid":"TravelAP","ifname":"phy1-ap0","up":true},
			{"section":"unknown","radio":"","ssid":"","ifname":"","up":null}
		]
	}`)
}

func TestWiFiScanJSONContractAndEmptyCollections(t *testing.T) {
	e := wifiEnv(t, fakeWiFi{scan: wifi.ScanResult{
		Networks: []wifi.Network{{
			SSID: "HomeUplink", Bands: []string{"2g", "5g"}, BestSignalDBM: -30, Security: "wpa2",
			AccessPoints: []wifi.AccessPoint{{BSSID: "60:22:32:2f:05:2b", Band: "5g", Channel: 149, FrequencyMHz: 5745, SignalDBM: -30, Quality: 70, QualityMax: 70, Security: "wpa2"}},
		}},
		Hidden: []wifi.AccessPoint{},
		Radios: []wifi.RadioScan{{Radio: "radio0", Band: "5g", Status: "ok", APCount: 1}},
	}})
	cookie, _ := e.login()
	got := e.do(req{method: http.MethodGet, path: "/api/v1/wifi/scan", cookie: cookie})
	assertWiFiJSON(t, got, `{
		"networks":[{
			"ssid":"HomeUplink","bands":["2g","5g"],"best_signal_dbm":-30,"security":"wpa2",
			"access_points":[{"bssid":"60:22:32:2f:05:2b","band":"5g","channel":149,"frequency_mhz":5745,"signal_dbm":-30,"quality":70,"quality_max":70,"security":"wpa2"}]
		}],
		"hidden":[],
		"radios":[{"radio":"radio0","band":"5g","status":"ok","ap_count":1}]
	}`)

	e = wifiEnv(t, fakeWiFi{scan: wifi.ScanResult{Networks: []wifi.Network{}, Hidden: []wifi.AccessPoint{}, Radios: []wifi.RadioScan{}}})
	cookie, _ = e.login()
	got = e.do(req{method: http.MethodGet, path: "/api/v1/wifi/scan", cookie: cookie})
	assertWiFiJSON(t, got, `{"networks":[],"hidden":[],"radios":[]}`)
}

func TestWiFiErrorMapping(t *testing.T) {
	for _, tc := range []struct {
		name, path, code string
		service          WiFiService
	}{
		{name: "nil summary", path: "/api/v1/wifi", code: "device_unavailable"},
		{name: "nil scan", path: "/api/v1/wifi/scan", code: "device_unavailable"},
		{name: "summary unavailable", path: "/api/v1/wifi", code: "wifi_unavailable", service: fakeWiFi{summaryErr: wifi.ErrUnavailable}},
		{name: "scan unavailable", path: "/api/v1/wifi/scan", code: "wifi_unavailable", service: fakeWiFi{scanErr: wifi.ErrUnavailable}},
		{name: "all scans failed", path: "/api/v1/wifi/scan", code: "wifi_scan_failed", service: fakeWiFi{scanErr: wifi.ErrScanFailed}},
		{name: "cancelled", path: "/api/v1/wifi/scan", code: "wifi_unavailable", service: fakeWiFi{scanErr: context.Canceled}},
		{name: "other", path: "/api/v1/wifi", code: "wifi_unavailable", service: fakeWiFi{summaryErr: errTest}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := wifiEnv(t, tc.service)
			cookie, _ := e.login()
			got := e.do(req{method: http.MethodGet, path: tc.path, cookie: cookie})
			if got.Code != http.StatusServiceUnavailable || errCode(t, got) != tc.code {
				t.Fatalf("%d %s", got.Code, got.Body)
			}
			if strings.Contains(got.Body.String(), errTest.Error()) {
				t.Fatalf("raw error leaked: %s", got.Body)
			}
		})
	}
}

func assertWiFiJSON(t *testing.T, got *httptest.ResponseRecorder, want string) {
	t.Helper()
	if got.Code != http.StatusOK {
		t.Fatalf("%d %s", got.Code, got.Body)
	}
	var actual, expected any
	if err := json.Unmarshal(got.Body.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got %s\nwant %s", got.Body, want)
	}
}
