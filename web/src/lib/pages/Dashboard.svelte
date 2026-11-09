<script lang="ts">
import { fetchHealth } from '../api/health';
import { m } from '../paraglide/messages.js';
import Card from '../ui/Card.svelte';
import LoadState from '../ui/LoadState.svelte';
import { createLoader } from '../ui/loader.svelte';
import Status from '../ui/Status.svelte';

const health = createLoader(fetchHealth);
</script>

<Card title={m.dashboard_daemon()}>
  <LoadState status={health.status} retry={health.retry}>
    <dl class="grid grid-cols-[auto_1fr] gap-x-6 gap-y-2">
      <dt class="text-muted">{m.dashboard_status()}</dt>
      <dd><Status kind="ok">{m.status_ok()}</Status></dd>
      <dt class="text-muted">{m.dashboard_version()}</dt>
      <dd>{health.data?.version}</dd>
    </dl>
  </LoadState>
</Card>
