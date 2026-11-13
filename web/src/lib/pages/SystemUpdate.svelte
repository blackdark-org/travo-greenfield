<script lang="ts">
import { onMount } from 'svelte';
import { getSystemUpdateStage, postSystemUpdateCheck, postSystemUpdateStage } from '../api/gen';
import '../api/client';
import { m } from '../paraglide/messages.js';
import Card from '../ui/Card.svelte';

type Update = NonNullable<Awaited<ReturnType<typeof postSystemUpdateCheck>>['data']>;
let result = $state<Update>();
type Stage = NonNullable<Awaited<ReturnType<typeof getSystemUpdateStage>>['data']>;
let stage = $state<Stage>();
let verifying = $state(false);
let verifyFailed = $state(false);
let stageRequest = 0;
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

async function readStage(request: AbortController) {
  const generation = ++stageRequest;
  try {
    const response = await getSystemUpdateStage({ signal: request.signal });
    if (!request.signal.aborted && generation === stageRequest) stage = response.data;
  } catch {
    if (!request.signal.aborted && generation === stageRequest) stage = undefined;
  }
}

async function verify(version: string, request: AbortController) {
  stageRequest++;
  verifying = true;
  verifyFailed = false;
  try {
    const response = await postSystemUpdateStage({ body: { version }, signal: request.signal });
    if (request.signal.aborted) return;
    if (response.data) stage = response.data;
    else {
      verifyFailed = true;
      await readStage(request);
    }
  } catch {
    if (!request.signal.aborted) {
      verifyFailed = true;
      await readStage(request);
    }
  } finally {
    if (!request.signal.aborted) verifying = false;
  }
}

onMount(() => {
  const request = new AbortController();
  controller = request;
  void check('automatic', request);
  void readStage(request);
  const timer = setInterval(() => {
    if (!verifying && stage?.status === 'verifying') void readStage(request);
  }, 2000);
  return () => {
    clearInterval(timer);
    request.abort();
  };
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
      disabled={pending || verifying || result?.status === 'disabled'}
      onclick={() => { if (controller) void check('manual', controller); }}
    >
      {m.update_check()}
    </button>
    {#if result?.status === 'available' && result.latest_version}
      <button
        type="button"
        class="rounded-md border border-border px-3 py-2 hover:bg-bg disabled:opacity-60"
        disabled={pending || verifying || stage?.status === 'verifying'}
        onclick={() => {
          if (controller && result?.latest_version) void verify(result.latest_version, controller);
        }}
      >
        {m.update_verify()}
      </button>
    {/if}
  </div>
  <div aria-live="polite">
    {#if verifying || stage?.status === 'verifying'}
      <p class="mt-2 text-muted">{m.update_verifying()}</p>
    {:else if stage?.status === 'verified' && stage.version}
      <p class="mt-2">{m.update_verified({ version: stage.version })}</p>
    {/if}
    {#if verifyFailed}
      <p class="mt-2 text-danger">{m.update_verify_failed()}</p>
    {/if}
  </div>
</Card>
