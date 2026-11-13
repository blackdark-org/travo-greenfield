import { render } from 'svelte/server';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { ScanResult } from '../api/gen';
import WiFi from './WiFi.svelte';

const { createLoader, fetchWifi, fetchWifiScan } = vi.hoisted(() => ({
  createLoader: vi.fn(),
  fetchWifi: vi.fn(),
  fetchWifiScan: vi.fn(),
}));
vi.mock('../ui/loader.svelte', () => ({ createLoader }));
vi.mock('../api/wifi', () => ({ fetchWifi, fetchWifiScan }));

const summary = { radios: [], stas: [], aps: [] };
const scan: ScanResult = {
  networks: [
    {
      ssid: 'HomeUplink',
      bands: ['5g'],
      best_signal_dbm: -30,
      security: 'wpa2',
      access_points: [
        {
          bssid: '60:22:32:2f:05:2b',
          band: '5g',
          channel: 149,
          frequency_mhz: 5745,
          signal_dbm: -30,
          quality: 70,
          quality_max: 70,
          security: 'wpa2',
        },
      ],
    },
  ],
  hidden: [],
  radios: [{ radio: 'radio0', band: '5g', status: 'ok', ap_count: 1 }],
};

describe('WiFi page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    createLoader.mockReturnValue({ status: { state: 'ready' }, data: summary, retry: vi.fn() });
  });

  it('starts idle and never scans automatically', () => {
    const body = render(WiFi).body;
    expect(body).toContain('Scan to see nearby networks.');
    expect(body).toContain('Scan');
    expect(fetchWifiScan).not.toHaveBeenCalled();
    expect(createLoader).toHaveBeenCalledWith(fetchWifi);
  });

  it('announces a slow scan and disables the button', () => {
    const body = render(WiFi, { props: { initialScanStatus: { state: 'loading' } } }).body;
    expect(body).toContain('Scanning, this can take several seconds');
    expect(body).toContain('aria-busy="true"');
    expect(body).toContain('disabled');
  });

  it('keeps prior scan data dimmed and stale on error with Retry', () => {
    const body = render(WiFi, {
      props: {
        initialScanStatus: {
          state: 'error',
          message: 'No radio could complete the scan. Try again.',
        },
        initialScanData: scan,
      },
    }).body;
    expect(body).toContain('HomeUplink');
    expect(body).toContain('opacity-60');
    expect(body).toContain('Last known value');
    expect(body).toContain('Retry');
  });
});
