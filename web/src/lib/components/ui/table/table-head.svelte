<script lang="ts">
import type { HTMLThAttributes } from 'svelte/elements';
import { cn, type WithElementRef } from '$lib/utils.js';

let {
  ref = $bindable(null),
  class: className,
  children,
  ...restProps
}: WithElementRef<HTMLThAttributes> = $props();
</script>

<th bind:this={ref} data-slot="table-head" class={cn("cn-table-head", className)} {...restProps}>
  {@render children?.()}
</th>
