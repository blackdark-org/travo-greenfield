<script lang="ts">
import LayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
import LogOut from '@lucide/svelte/icons/log-out';
import Menu from '@lucide/svelte/icons/menu';
import MoonStar from '@lucide/svelte/icons/moon-star';
import Network from '@lucide/svelte/icons/network';
import Server from '@lucide/svelte/icons/server';
import Settings from '@lucide/svelte/icons/settings';
import Shield from '@lucide/svelte/icons/shield';
import Sun from '@lucide/svelte/icons/sun';
import Wifi from '@lucide/svelte/icons/wifi';
import X from '@lucide/svelte/icons/x';
import { dashboardEvents } from './lib/dashboard-events.svelte';
import { navigate } from './lib/navigate';
import { pages } from './lib/pages/lazy';
import PageHost from './lib/pages/PageHost.svelte';
import { m } from './lib/paraglide/messages.js';
import { guard, href, type Route } from './lib/router';
import { router } from './lib/router.svelte';
import { session } from './lib/session.svelte';
import { theme } from './lib/theme.svelte';
import LoadState from './lib/ui/LoadState.svelte';

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

// Groups per spec. Sections without a route; the rest render the "not yet" card.
const nav: { route: Route | null; label: () => string; icon: typeof Wifi }[] = [
  { route: 'dashboard', label: () => m.nav_dashboard(), icon: LayoutDashboard },
  { route: 'wifi', label: () => m.nav_wifi(), icon: Wifi },
  { route: null, label: () => m.nav_network(), icon: Network },
  { route: 'vpn', label: () => m.nav_vpn(), icon: Shield },
  { route: null, label: () => m.nav_services(), icon: Server },
  { route: 'system', label: () => m.nav_system(), icon: Settings },
];

void session.bootstrap();

function logout() {
  dashboardEvents.clear();
  // Leave the page first so the guard does not remember it as the place to return to.
  navigate('dashboard');
  void session.logout();
}

const authenticated = $derived(session.state === 'authenticated');
// The shell, not a page, owns the stream. Navigation does not restart it; login,
// expiry and a replacement session after password change each key a new owner.
$effect(() => {
  const generation = session.generation;
  if (!authenticated) {
    dashboardEvents.clear();
    return;
  }
  dashboardEvents.start(generation);
  return () => dashboardEvents.stop(generation);
});
// What is shown: the requested route after the guard. While the state is unknown nothing is shown
// but the loader (or the bootstrap error), so a reload never flashes the wrong page.
const shown = $derived(guard(session.state, router.route));
const bootstrapStatus = $derived(
  session.error === null
    ? ({ state: 'loading' } as const)
    : ({ state: 'error', message: session.error } as const),
);

// Apply the guard to the address bar too, and remember where an expired or anonymous user was going.
// Leaving login/setup after authenticating goes to returnTo here, in one step, so the dashboard
// never mounts in between.
$effect(() => {
  const wanted = router.route;
  if (shown === wanted) return;
  if (shown === 'login') session.returnTo = wanted === 'setup' ? null : wanted;
  if (session.state === 'authenticated' && (wanted === 'login' || wanted === 'setup')) {
    const next = session.returnTo ?? 'dashboard';
    session.returnTo = null;
    navigate(next, { replace: true });
    return;
  }
  navigate(shown, { replace: true });
});

// Close the drawer after navigating.
$effect(() => {
  void router.route;
  drawerOpen = false;
});
</script>

<svelte:window
  onkeydown={(e) => e.key === 'Escape' && closeDrawer()}
  onpagehide={() => dashboardEvents.stop()}
/>

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

  {#if authenticated}
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
                aria-current={shown === item.route ? 'page' : undefined}
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
  {/if}

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
        class:invisible={!authenticated}
        onclick={openDrawer}
      >
        <Menu aria-hidden="true" size={20} />
      </button>
      <div class="flex items-center gap-1">
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
        {#if authenticated}
          <button
            type="button"
            class="rounded-md p-2 hover:bg-bg"
            aria-label={m.logout()}
            onclick={logout}
          >
            <LogOut aria-hidden="true" size={20} />
          </button>
        {/if}
      </div>
    </header>

    <main class="flex-1 overflow-y-auto p-4">
      <div class="mx-auto max-w-3xl">
        {#key shown}
          {#if session.state === 'unknown'}
            <LoadState status={bootstrapStatus} retry={() => session.bootstrap()}>
              <span></span>
            </LoadState>
          {:else if shown === 'dashboard'}
            <PageHost load={pages.dashboard} props={{ events: dashboardEvents }} />
          {:else if shown === 'vpn'}
            <PageHost load={pages.vpn} />
          {:else if shown === 'system'}
            <PageHost load={pages.system} />
          {:else if shown === 'login'}
            <PageHost load={pages.login} />
          {:else if shown === 'setup'}
            <PageHost load={pages.setup} />
          {:else}
            <PageHost load={pages.notYet} />
          {/if}
        {/key}
      </div>
    </main>
  </div>
</div>
