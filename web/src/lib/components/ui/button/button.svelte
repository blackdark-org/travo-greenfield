<script lang="ts" module>
import type { HTMLAnchorAttributes, HTMLButtonAttributes } from 'svelte/elements';
import { cn, type WithElementRef } from '$lib/utils.js';

const styles = {
  base: 'cn-button group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap transition-all outline-none select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0',
  variants: {
    variant: {
      default: 'cn-button-variant-default',
      outline: 'cn-button-variant-outline',
      secondary: 'cn-button-variant-secondary',
      ghost: 'cn-button-variant-ghost',
      destructive: 'cn-button-variant-destructive',
      link: 'cn-button-variant-link',
    },
    size: {
      default: 'cn-button-size-default',
      sm: 'cn-button-size-sm',
      icon: 'cn-button-size-icon',
      'icon-sm': 'cn-button-size-icon-sm',
    },
  },
} as const;

export function buttonVariants({
  variant = 'default',
  size = 'default',
}: {
  variant?: ButtonVariant;
  size?: ButtonSize;
} = {}) {
  return cn(styles.base, styles.variants.variant[variant], styles.variants.size[size]);
}

export type ButtonVariant = keyof typeof styles.variants.variant;
export type ButtonSize = keyof typeof styles.variants.size;

export type ButtonProps = WithElementRef<HTMLButtonAttributes> &
  WithElementRef<HTMLAnchorAttributes> & {
    variant?: ButtonVariant;
    size?: ButtonSize;
  };
</script>

<script lang="ts">
  let {
    class: className,
    variant = "default",
    size = "default",
    ref = $bindable(null),
    href = undefined,
    type = "button",
    disabled,
    children,
    ...restProps
  }: ButtonProps = $props();
</script>

{#if href}
  <a
    bind:this={ref}
    data-slot="button"
    class={cn(buttonVariants({ variant, size }), className)}
    href={disabled ? undefined : href}
    aria-disabled={disabled}
    role={disabled ? "link" : undefined}
    tabindex={disabled ? -1 : undefined}
    {...restProps}
  >
    {@render children?.()}
  </a>
{:else}
  <button
    bind:this={ref}
    data-slot="button"
    class={cn(buttonVariants({ variant, size }), className)}
    {type}
    {disabled}
    {...restProps}
  >
    {@render children?.()}
  </button>
{/if}
