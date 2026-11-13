<script lang="ts">
import { onDestroy } from 'svelte';
import type { ScanResult } from '../api/gen';
import { fetchWifi, fetchWifiScan } from '../api/wifi';
import { m } from '../paraglide/messages.js';
import Card from '../ui/Card.svelte';
import LoadState from '../ui/LoadState.svelte';
import { createLoader, type LoadStatus } from '../ui/loader.svelte';
import { wifiScanErrorMessage } from '../wifi';
import WiFiScanList from './WiFiScanList.svelte';
import WiFiSummary from './WiFiSummary.svelte';

const idle: LoadStatus = { state: 'ready' };
let {
  initialScanStatus = idle,
  initialScanData,
}: { initialScanStatus?: LoadStatus; initialScanData?: ScanResult } = $props();

const summary = createLoader(fetchWifi);
const request = new AbortController();
// The optional initial values exist for deterministic SSR state tests; normal routes start idle.
// svelte-ignore state_referenced_locally
let scanStatus = $state<LoadStatus>(initialScanStatus);
// svelte-ignore state_referenced_locally
let scanData = $state<ScanResult | undefined>(initialScanData);

onDestroy(() => request.abort());

async function scan() {
  if (scanStatus.state === 'loading') return;
  scanStatus = { state: 'loading' };
  try {
    const result = await fetchWifiScan(request.signal);
    if (request.signal.aborted) return;
    scanData = result;
    scanStatus = { state: 'ready' };
  } catch (error) {
    if (!request.signal.aborted) {
      scanStatus = { state: 'error', message: wifiScanErrorMessage(error) };
    }
  }
}
</script>

<div class="grid gap-4">
  <Card title={m.wifi_summary_title()}>
    <LoadState status={summary.status} retry={summary.retry} stale={summary.data !== undefined}>
      {#if summary.data}
        <WiFiSummary data={summary.data} />
      {/if}
    </LoadState>
    {#if summary.status.state === 'ready'}
      <button
        type="button"
        class="mt-3 rounded-md border border-border px-3 py-1 text-sm hover:bg-bg"
        onclick={summary.retry}
      >
        {m.wifi_refresh()}
      </button>
    {/if}
  </Card>

  <Card title={m.wifi_scan_title()}>
    <LoadState status={scanStatus} retry={scan} stale={scanData !== undefined}>
      {#if scanData}
        <WiFiScanList result={scanData} />
      {:else}
        <p class="text-muted">{m.wifi_scan_idle()}</p>
      {/if}
    </LoadState>
    {#if scanStatus.state !== 'error'}
      <button
        type="button"
        class="mt-3 rounded-md border border-border px-3 py-1 text-sm hover:bg-bg disabled:opacity-60"
        disabled={scanStatus.state === 'loading'}
        aria-busy={scanStatus.state === 'loading'}
        onclick={scan}
      >
        {scanStatus.state === 'loading'
          ? m.wifi_scanning()
          : scanData
            ? m.wifi_scan_again()
            : m.wifi_scan()}
      </button>
    {/if}
  </Card>
</div>
