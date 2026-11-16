<script lang="ts">
import type { Snippet } from 'svelte';
import Root from '$lib/components/ui/card/card.svelte';
import Title from '$lib/components/ui/card/card-title.svelte';

let { title, children }: { title?: string; children: Snippet } = $props();
</script>

<Root>
  {#if title}
    <Title>{title}</Title>
  {/if}
  {@render children()}
</Root>
