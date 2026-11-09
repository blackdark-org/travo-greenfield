<script lang="ts">
import LayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
import Menu from '@lucide/svelte/icons/menu';
import MoonStar from '@lucide/svelte/icons/moon-star';
import Network from '@lucide/svelte/icons/network';
import Server from '@lucide/svelte/icons/server';
import Settings from '@lucide/svelte/icons/settings';
import Shield from '@lucide/svelte/icons/shield';
import Sun from '@lucide/svelte/icons/sun';
import Wifi from '@lucide/svelte/icons/wifi';
import X from '@lucide/svelte/icons/x';
import { pages } from './lib/pages/lazy';
import PageHost from './lib/pages/PageHost.svelte';
import { m } from './lib/paraglide/messages.js';
import { href, type Route } from './lib/router';
import { router } from './lib/router.svelte';
import { theme } from './lib/theme.svelte';

let drawerOpen = $state(false);
let menuButton: HTMLButtonElement | undefined = $state();
let closeButton: HTMLButtonElement | undefined = $state();
// Below md the sidebar is a drawer: closed means inert (not tabbable, hidden from AT).
const mobile = $state({ value: false });
$effect(() => {
  const mq = matchMedia('(max-width: 767.98px)');
  const sync = () => (mobile.value = mq.matches);
  sync();
  mq.addEventListener('change', sync);
  return () => mq.removeEventListener('change', sync);
});

function openDrawer() {
  drawerOpen = true;
  queueMicrotask(() => closeButton?.focus());
}

function closeDrawer() {
  if (!drawerOpen) return;
  drawerOpen = false;
  queueMicrotask(() => menuButton?.focus());
}

// Groups per spec. Only Dashboard and System have routes yet; the rest render the "not yet" card.
const nav: { route: Route | null; label: () => string; icon: typeof Wifi }[] = [
  { route: 'dashboard', label: () => m.nav_dashboard(), icon: LayoutDashboard },
  { route: 'wifi', label: () => m.nav_wifi(), icon: Wifi },
  { route: null, label: () => m.nav_network(), icon: Network },
  { route: null, label: () => m.nav_vpn(), icon: Shield },
  { route: null, label: () => m.nav_services(), icon: Server },
  { route: 'system', label: () => m.nav_system(), icon: Settings },
];

// Close the drawer after navigating.
$effect(() => {
  void router.route;
  drawerOpen = false;
});
</script>

<svelte:window onkeydown={(e) => e.key === 'Escape' && closeDrawer()} />

<div class="flex h-full w-full overflow-hidden">
  {#if drawerOpen}
    <!-- pointer-only backdrop: keyboard users close via the Close button or Escape -->
    <button
      type="button"
      class="fixed inset-0 z-20 bg-black/50 md:hidden"
      aria-hidden="true"
      tabindex="-1"
      onclick={closeDrawer}
    ></button>
  {/if}

  <aside
    id="sidebar"
    inert={mobile.value && !drawerOpen}
    class="fixed inset-y-0 left-0 z-30 flex w-60 shrink-0 flex-col border-r border-border bg-surface transition-transform md:static md:translate-x-0 {drawerOpen
      ? 'translate-x-0'
      : '-translate-x-full'}"
  >
    <div class="flex items-center justify-between px-4 py-3">
      <span class="text-lg font-semibold">{m.app_name()}</span>
      <button
        type="button"
        class="rounded-md p-1 md:hidden"
        aria-label={m.nav_close()}
        bind:this={closeButton}
        onclick={closeDrawer}
      >
        <X aria-hidden="true" size={20} />
      </button>
    </div>
    <nav aria-label={m.nav_label()} class="flex-1 overflow-y-auto px-2 pb-4">
      <ul class="flex flex-col gap-1">
        {#each nav as item (item.label())}
          {@const Icon = item.icon}
          <li>
            {#if item.route}
              <a
                href={href(item.route)}
                aria-current={router.route === item.route ? 'page' : undefined}
                class="flex items-center gap-3 rounded-md px-3 py-2 hover:bg-bg aria-[current=page]:bg-bg aria-[current=page]:font-semibold"
                ><Icon aria-hidden="true" size={18} />{item.label()}</a
              >
            {:else}
              <span class="flex items-center gap-3 rounded-md px-3 py-2 text-muted"
                ><Icon aria-hidden="true" size={18} />{item.label()}</span
              >
            {/if}
          </li>
        {/each}
      </ul>
    </nav>
  </aside>

  <!-- inert while the drawer is open so Tab cannot leave the drawer -->
  <div class="flex min-w-0 flex-1 flex-col" inert={mobile.value && drawerOpen}>
    <header class="flex items-center justify-between border-b border-border bg-surface px-4 py-2">
      <button
        type="button"
        class="rounded-md p-2 md:invisible"
        aria-label={m.nav_open()}
        aria-controls="sidebar"
        aria-expanded={drawerOpen}
        bind:this={menuButton}
        onclick={openDrawer}
      >
        <Menu aria-hidden="true" size={20} />
      </button>
      <button
        type="button"
        class="rounded-md p-2 hover:bg-bg"
        aria-label={m.theme_toggle()}
        onclick={() => theme.toggle()}
      >
        {#if theme.value === 'dark'}
          <Sun aria-hidden="true" size={20} />
        {:else}
          <MoonStar aria-hidden="true" size={20} />
        {/if}
      </button>
    </header>

    <main class="flex-1 overflow-y-auto p-4">
      <div class="mx-auto max-w-3xl">
        {#key router.route}
          {#if router.route === 'dashboard'}
            <PageHost load={pages.dashboard} />
          {:else if router.route === 'system'}
            <PageHost load={pages.system} />
          {:else if router.route === 'login'}
            <PageHost
              load={pages.info}
              props={{ title: m.login_title(), body: m.login_body() }}
            />
          {:else if router.route === 'setup'}
            <PageHost
              load={pages.info}
              props={{ title: m.setup_title(), body: m.setup_body() }}
            />
          {:else}
            <PageHost load={pages.notYet} />
          {/if}
        {/key}
      </div>
    </main>
  </div>
</div>
