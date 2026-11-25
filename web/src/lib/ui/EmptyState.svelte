<script lang="ts">
import type { Component, Snippet } from 'svelte';

// "Nothing here yet" content. A failed load is never empty: LoadState shows that.
let {
  icon: Icon,
  title,
  text,
  action,
}: {
  icon: Component<{ size?: number; class?: string; 'aria-hidden'?: 'true' }>;
  title: string;
  text?: string;
  action?: Snippet;
} = $props();
</script>

<div data-slot="empty-state" class="flex flex-col items-center gap-1.5 py-4 text-center">
  <Icon aria-hidden="true" size={20} class="text-muted-foreground" />
  <p class="max-w-prose text-sm font-medium">{title}</p>
  {#if text}<p class="max-w-prose text-sm text-muted-foreground">{text}</p>{/if}
  {#if action}<div class="mt-1">{@render action()}</div>{/if}
</div>
