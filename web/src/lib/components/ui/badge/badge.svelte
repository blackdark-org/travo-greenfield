<script lang="ts" module>
import { cn } from '$lib/utils.js';

const styles = {
  base: 'cn-badge inline-flex w-fit items-center justify-center',
  variants: {
    variant: {
      default: 'cn-badge-variant-default',
      secondary: 'cn-badge-variant-secondary',
      destructive: 'cn-badge-variant-destructive',
      outline: 'cn-badge-variant-outline',
      ghost: 'cn-badge-variant-ghost',
      link: 'cn-badge-variant-link',
    },
  },
  defaultVariants: {
    variant: 'default',
  },
} as const;
export function badgeVariants({ variant = 'default' }: { variant?: BadgeVariant } = {}) {
  return cn(styles.base, styles.variants.variant[variant]);
}

export type BadgeVariant = keyof typeof styles.variants.variant;
</script>

<script lang="ts">
  import { type WithElementRef } from "$lib/utils.js";
  import type { HTMLAnchorAttributes } from "svelte/elements";

  let {
    ref = $bindable(null),
    href,
    class: className,
    variant = "default",
    children,
    ...restProps
  }: WithElementRef<HTMLAnchorAttributes> & {
    variant?: BadgeVariant;
  } = $props();
</script>

<svelte:element
  this={href ? "a" : "span"}
  bind:this={ref}
  data-slot="badge"
  {href}
  class={cn(badgeVariants({ variant }), className)}
  {...restProps}
>
  {@render children?.()}
</svelte:element>
