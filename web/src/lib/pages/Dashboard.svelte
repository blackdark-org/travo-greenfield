<script lang="ts">
import { fetchHealth } from '../api/health';
import { fetchSystemInfo } from '../api/system';
import { fetchThroughput } from '../api/throughput';
import { formatBytes, formatLoad, formatUptime } from '../format';
import { m } from '../paraglide/messages.js';
import Card from '../ui/Card.svelte';
import LoadState from '../ui/LoadState.svelte';
import { createLoader } from '../ui/loader.svelte';
import Status from '../ui/Status.svelte';

const health = createLoader(fetchHealth);
const resources = createLoader(fetchSystemInfo);
const throughput = createLoader(fetchThroughput, 2000);

function formatRate(rate: number | null): string {
  if (rate === null) return m.dashboard_rate_unknown();
  const rounded = rate < 1024 ? Math.round(rate * 10) / 10 : rate;
  return m.dashboard_bytes_per_second({ bytes: formatBytes(rounded) });
}
</script>

<div class="flex flex-col gap-4">
  <Card title={m.dashboard_daemon()}>
    <LoadState status={health.status} retry={health.retry}>
      <dl class="grid grid-cols-[auto_1fr] gap-x-6 gap-y-2">
        <dt class="text-muted">{m.dashboard_status()}</dt>
        <dd><Status kind="ok">{m.status_ok()}</Status></dd>
      </dl>
    </LoadState>
  </Card>

  <Card title={m.dashboard_resources()}>
    <LoadState status={resources.status} retry={resources.retry} stale={resources.data !== undefined}>
      {#if resources.data}
        {@const info = resources.data}
        <dl class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2">
          <dt class="text-muted">{m.dashboard_cpu_load()}</dt>
          <dd>{formatLoad([info.load[0]])}</dd>
          <dt class="text-muted">{m.system_memory()}</dt>
          <dd>
            {#if info.memory.total_bytes === 0}
              <Status kind="neutral">{m.dashboard_unavailable()}</Status>
            {:else}
              {m.system_used_of({
                used: formatBytes(Math.max(0, info.memory.total_bytes - info.memory.available_bytes)),
                total: formatBytes(info.memory.total_bytes),
              })}
            {/if}
          </dd>
          <dt class="text-muted">{m.system_uptime()}</dt>
          <dd>{formatUptime(info.uptime_seconds)}</dd>
          <dt class="text-muted">{m.dashboard_version()}</dt>
          <dd class="break-all">{info.app.version}</dd>
        </dl>
        {#if resources.status.state === 'ready'}
          <button
            type="button"
            class="mt-3 rounded-md border border-border px-3 py-1 text-sm hover:bg-bg"
            onclick={resources.retry}
          >
            {m.dashboard_refresh()}
          </button>
        {/if}
      {/if}
    </LoadState>
  </Card>
  <Card title={m.dashboard_throughput()}>
    <LoadState status={throughput.status} retry={throughput.retry} stale={throughput.data !== undefined}>
      {#if throughput.data?.interfaces.length}
        <table class="w-full table-fixed text-left">
          <thead class="text-muted">
            <tr>
              <th scope="col" class="font-normal">{m.dashboard_interface()}</th>
              <th scope="col" class="text-right font-normal">{m.dashboard_rx()}</th>
              <th scope="col" class="text-right font-normal">{m.dashboard_tx()}</th>
            </tr>
          </thead>
          <tbody>
            {#each throughput.data.interfaces as row (row.name)}
              <tr class="border-t border-border">
                <th scope="row" class="break-all py-2 font-medium">{row.name}</th>
                <td class="py-2 text-right tabular-nums">{formatRate(row.rx_bytes_per_second)}</td>
                <td class="py-2 text-right tabular-nums">{formatRate(row.tx_bytes_per_second)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {:else}
        <p class="text-muted">{m.dashboard_no_interfaces()}</p>
      {/if}
    </LoadState>
  </Card>
</div>
