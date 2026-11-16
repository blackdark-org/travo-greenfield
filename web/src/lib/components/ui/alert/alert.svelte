<script lang="ts" module>
import { cn } from '$lib/utils.js';

const styles = {
  base: 'cn-alert group/alert relative w-full',
  variants: {
    variant: {
      default: 'cn-alert-variant-default',
      destructive: 'cn-alert-variant-destructive',
    },
  },
  defaultVariants: {
    variant: 'default',
  },
} as const;
export function alertVariants({ variant = 'default' }: { variant?: AlertVariant } = {}) {
  return cn(styles.base, styles.variants.variant[variant]);
}

export type AlertVariant = keyof typeof styles.variants.variant;
</script>

<script lang="ts">
  import { type WithElementRef } from "$lib/utils.js";
  import type { HTMLAttributes } from "svelte/elements";

  let {
    ref = $bindable(null),
    class: className,
    variant = "default",
    children,
    ...restProps
  }: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
    variant?: AlertVariant;
  } = $props();
</script>

<div
  bind:this={ref}
  data-slot="alert"
  role="alert"
  class={cn(alertVariants({ variant }), className)}
  {...restProps}
>
  {@render children?.()}
</div>
