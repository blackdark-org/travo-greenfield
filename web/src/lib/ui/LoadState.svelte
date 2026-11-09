<script lang="ts">
import type { Snippet } from 'svelte';
import { m } from '../paraglide/messages.js';
import type { LoadStatus } from './loader.svelte';

let { status, retry, children }: { status: LoadStatus; retry: () => void; children: Snippet } =
  $props();
</script>

{#if status.state === 'loading'}
  <div role="status" aria-live="polite" aria-busy="true">
    <span class="sr-only">{m.load_loading()}</span>
    <div class="skeleton h-5 w-40"></div>
  </div>
{:else if status.state === 'error'}
  <div role="alert" class="flex flex-wrap items-center gap-3">
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
