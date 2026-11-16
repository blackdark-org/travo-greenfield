<script lang="ts">
import type { HTMLAttributes } from 'svelte/elements';
import { cn, type WithElementRef } from '$lib/utils.js';

let {
  ref = $bindable(null),
  class: className,
  children,
  size = 'default',
  ...restProps
}: WithElementRef<HTMLAttributes<HTMLElement>> & { size?: 'default' | 'sm' } = $props();
</script>

<section
  bind:this={ref}
  data-slot="card"
  data-size={size}
  class={cn("cn-card group/card", className)}
  {...restProps}
>
  {@render children?.()}
</section>
