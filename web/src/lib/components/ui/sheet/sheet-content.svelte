<script lang="ts" module>
export type Side = 'top' | 'right' | 'bottom' | 'left';
</script>

<script lang="ts">
  import { Dialog as SheetPrimitive } from "bits-ui";
  import X from "@lucide/svelte/icons/x";
import { m } from "$lib/paraglide/messages.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { cn, type WithoutChildrenOrChild } from "$lib/utils.js";
  import SheetOverlay from "./sheet-overlay.svelte";
  import SheetPortal from "./sheet-portal.svelte";
  import type { Snippet } from "svelte";
  import type { ComponentProps } from "svelte";

  let {
    ref = $bindable(null),
    class: className,
    side = "left",
    showCloseButton = true,
    portalProps,
    children,
    ...restProps
  }: WithoutChildrenOrChild<SheetPrimitive.ContentProps> & {
    portalProps?: WithoutChildrenOrChild<ComponentProps<typeof SheetPortal>>;
    side?: Side;
    showCloseButton?: boolean;
    children: Snippet;
  } = $props();
</script>

<SheetPortal {...portalProps}>
  <SheetOverlay />
  <SheetPrimitive.Content
    bind:ref
    data-slot="sheet-content"
    data-side={side}
    class={cn(
      "cn-sheet-content",
      className
    )}
    {...restProps}
  >
    {@render children?.()}
    {#if showCloseButton}
      <SheetPrimitive.Close data-slot="sheet-close">
        {#snippet child({ props })}
          <Button variant="ghost" class="cn-sheet-close" size="icon-sm" {...props}>
            <X />
            <span class="sr-only">{m.nav_close()}</span>
          </Button>
        {/snippet}
      </SheetPrimitive.Close>
    {/if}
  </SheetPrimitive.Content>
</SheetPortal>
