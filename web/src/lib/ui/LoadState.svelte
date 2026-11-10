<script lang="ts">
import type { Snippet } from 'svelte';
import { m } from '../paraglide/messages.js';
import type { LoadStatus } from './loader.svelte';
import Status from './Status.svelte';

// stale: a value is already known, so a failed refresh dims it instead of blanking it.
let {
  status,
  retry,
  stale = false,
  children,
}: { status: LoadStatus; retry: () => void; stale?: boolean; children: Snippet } = $props();
</script>

{#if status.state === 'loading' && stale}
  <!-- a retry of a known value: keep it, dimmed, instead of a skeleton -->
  <div class="opacity-60" role="status" aria-live="polite" aria-busy="true">
    <span class="sr-only">{m.load_loading()}</span>
    {@render children()}
  </div>
{:else if status.state === 'loading'}
  <div role="status" aria-live="polite" aria-busy="true">
    <span class="sr-only">{m.load_loading()}</span>
    <div class="skeleton h-5 w-40"></div>
  </div>
{:else if status.state === 'error'}
  {#if stale}
    <div class="opacity-60">{@render children()}</div>
  {/if}
  <div role="alert" class="flex flex-wrap items-center gap-3">
    {#if stale}
      <Status kind="stale">{m.load_stale()}</Status>
    {/if}
    <span class="status status-danger">{m.load_error({ cause: status.message })}</span>
    <button
      type="button"
      class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg"
      onclick={retry}
    >
      {m.load_retry()}
    </button>
  </div>
{:else}
  {@render children()}
{/if}
