<script lang="ts">
import { onMount } from 'svelte';
import { postSystemUpdateCheck } from '../api/gen';
import '../api/client';
import { m } from '../paraglide/messages.js';
import Card from '../ui/Card.svelte';

type Update = NonNullable<Awaited<ReturnType<typeof postSystemUpdateCheck>>['data']>;
let result = $state<Update>();
let pending = $state(true);
let controller: AbortController | undefined;
const labels = {
  disabled: m.update_disabled,
  deferred: m.update_deferred,
  no_release: m.update_no_release,
  source_unavailable: m.update_source_unavailable,
  rate_limited: m.update_rate_limited,
  unavailable: m.update_unavailable,
  up_to_date: m.update_up_to_date,
};

async function check(mode: 'automatic' | 'manual', request: AbortController) {
  pending = true;
  try {
    const response = await postSystemUpdateCheck({ body: { mode }, signal: request.signal });
    if (!request.signal.aborted) result = response.data;
  } catch {
    if (!request.signal.aborted) result = undefined;
  } finally {
    if (!request.signal.aborted) pending = false;
  }
}

onMount(() => {
  const request = new AbortController();
  controller = request;
  void check('automatic', request);
  return () => request.abort();
});
</script>

<Card title={m.update_title()}>
  <div class="flex flex-wrap items-center gap-3">
    <p aria-live="polite" class="text-muted">
      {#if pending}
        {m.update_checking()}
      {:else if result?.status === 'available'}
        {m.update_available({ version: result.latest_version ?? '' })}
      {:else}
        {labels[result?.status ?? 'unavailable']()}
      {/if}
    </p>
    <button
      type="button"
      class="rounded-md border border-border px-3 py-2 hover:bg-bg disabled:opacity-60"
      disabled={pending || result?.status === 'disabled'}
      onclick={() => { if (controller) void check('manual', controller); }}
    >
      {m.update_check()}
    </button>
  </div>
</Card>
