<script lang="ts">
import { onDestroy } from 'svelte';
import '../api/client';
import { deleteVpnProfile, getVpnProfiles, postVpnProfile, selectVpnProfile } from '../api/gen';
import { m } from '../paraglide/messages.js';
import Card from '../ui/Card.svelte';
import LoadState from '../ui/LoadState.svelte';
import { createLoader } from '../ui/loader.svelte';

const request = new AbortController();
const profiles = createLoader(async () => {
  const result = await getVpnProfiles({ signal: request.signal });
  if (!result.data) throw new Error(m.vpn_unavailable());
  return result.data;
});
let name = $state('');
let config = $state('');
let pending = $state(false);
let error = $state('');
let removing = $state<string>();
onDestroy(() => {
  request.abort();
  config = '';
});

async function change(action: 'import' | 'select' | 'delete', id = '') {
  if (pending) return;
  pending = true;
  error = '';
  try {
    const options = { signal: request.signal };
    const response =
      action === 'import'
        ? await postVpnProfile({ ...options, body: { name, config } })
        : action === 'select'
          ? await selectVpnProfile({ ...options, path: { id } })
          : await deleteVpnProfile({ ...options, path: { id } });
    if (request.signal.aborted) return;
    if (!response.data) {
      const code = (response.error as { code?: string } | undefined)?.code;
      error =
        code === 'invalid_profile'
          ? m.vpn_invalid()
          : code === 'profile_limit'
            ? m.vpn_limit()
            : m.vpn_save_failed();
      return;
    }
    if (action === 'import') {
      name = '';
      config = '';
    }
    removing = undefined;
    await profiles.retry();
  } catch {
    if (!request.signal.aborted) error = m.vpn_save_failed();
  } finally {
    if (!request.signal.aborted) pending = false;
  }
}
</script>

<div class="grid gap-4">
  <Card title={m.vpn_saved_profiles()}>
    <LoadState status={profiles.status} retry={profiles.retry}>
      {#if profiles.data}
        {#if profiles.data.profiles.length === 0}
          <p class="text-muted">{m.vpn_no_profiles()}</p>
        {:else}
          <ul class="space-y-3">
            {#each profiles.data.profiles as profile (profile.id)}
              <li class="flex flex-wrap items-center justify-between gap-2">
                <span class="break-all">{profile.name}</span>
                <div class="flex flex-wrap items-center gap-2">
                  {#if profile.id === profiles.data.selected_id}
                    <span class="text-muted">{m.vpn_selected()}</span>
                  {:else}
                    <button type="button" class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg disabled:opacity-60" disabled={pending} onclick={() => change('select', profile.id)}>{m.vpn_select()}</button>
                  {/if}
                  {#if removing === profile.id}
                    <button type="button" class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg disabled:opacity-60" disabled={pending} onclick={() => change('delete', profile.id)}>{m.vpn_confirm_delete()}</button>
                    <button type="button" class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg disabled:opacity-60" disabled={pending} onclick={() => { removing = undefined; }}>{m.vpn_cancel()}</button>
                  {:else}
                    <button type="button" class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg disabled:opacity-60" disabled={pending} onclick={() => { removing = profile.id; }}>{m.vpn_delete()}</button>
                  {/if}
                </div>
              </li>
            {/each}
          </ul>
        {/if}
      {/if}
    </LoadState>
  </Card>

  <Card title={m.vpn_import()}>
    <form class="grid gap-3" autocomplete="off" onsubmit={(event) => { event.preventDefault(); void change('import'); }}>
      <label class="grid gap-1">
        <span>{m.vpn_name()}</span>
        <input class="rounded-md border border-border bg-bg px-3 py-2" bind:value={name} required disabled={pending} />
      </label>
      <label class="grid gap-1">
        <span>{m.vpn_config()}</span>
        <textarea class="w-full rounded-md border border-border bg-bg px-3 py-2 font-mono" rows="6" bind:value={config} required disabled={pending} spellcheck="false"></textarea>
      </label>
      <button type="submit" class="w-fit rounded-md bg-accent px-4 py-2 text-accent-contrast disabled:opacity-60" disabled={pending || !name.trim() || !config.trim()}>{m.vpn_import()}</button>
    </form>
  </Card>
  {#if error}<p role="alert" class="text-danger">{error}</p>{/if}
</div>
