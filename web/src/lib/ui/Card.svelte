<script lang="ts">
import type { Snippet } from 'svelte';

let { title, children }: { title?: string; children: Snippet } = $props();
</script>

<section class="rounded-lg border border-border bg-surface p-4">
  {#if title}
    <h2 class="mb-3 text-base font-semibold">{title}</h2>
  {/if}
  {@render children()}
</section>
