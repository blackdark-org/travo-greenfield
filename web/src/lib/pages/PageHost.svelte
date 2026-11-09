<script lang="ts">
import LoadState from '../ui/LoadState.svelte';
import { createLoader } from '../ui/loader.svelte';
import type { PageImport } from './lazy';

// Keyed by route in App, so each route gets a fresh loader for its chunk.
let { load, props = {} }: { load: PageImport; props?: Record<string, string> } = $props();

const page = createLoader(() => load());
const Page = $derived(page.data);
</script>

<LoadState status={page.status} retry={page.retry}>
  {#if Page}
    <Page {...props} />
  {/if}
</LoadState>
