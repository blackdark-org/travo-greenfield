<script lang="ts">
import { type AuthMessage, passwordMessage } from '../auth-messages';
import { navigate } from '../navigate';
import { m } from '../paraglide/messages.js';
import { passwordRulesOk } from '../password-rules';
import { session } from '../session.svelte';
import { theme } from '../theme.svelte';
import Card from '../ui/Card.svelte';
import FormMessage from '../ui/FormMessage.svelte';
import PasswordField from '../ui/PasswordField.svelte';

let current = $state('');
let next = $state('');
let confirm = $state('');
let pending = $state(false);
let failure = $state<AuthMessage | null>(null);
let mismatch = $state(false);
let changed = $state(false);

async function submit(event?: SubmitEvent) {
  event?.preventDefault();
  failure = null;
  mismatch = false;
  changed = false;
  if (!passwordRulesOk(next)) {
    failure = { text: m.password_rules(), field: 'new' };
    return;
  }
  if (next !== confirm) {
    mismatch = true;
    return;
  }
  pending = true;
  const result = await session.changePassword(current, next);
  pending = false;
  if (!result.ok) {
    failure = passwordMessage(result.code);
    return;
  }
  current = next = confirm = '';
  changed = true;
}

async function logout() {
  // Leave the page first so the guard does not remember System as the place to return to.
  navigate('dashboard');
  await session.logout();
}
</script>

<div class="flex flex-col gap-4">
  <Card title={m.system_about()}>
    <div class="flex items-center justify-between gap-4">
      <span>{m.system_theme()}</span>
      <button
        type="button"
        class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg"
        aria-pressed={theme.value === 'dark'}
        onclick={() => theme.toggle()}
      >
        {m.theme_toggle()}
      </button>
    </div>
  </Card>

  <Card title={m.password_change_title()}>
    <form class="flex flex-col gap-4" onsubmit={submit}>
      <PasswordField
        id="password-current"
        label={m.password_current()}
        autocomplete="current-password"
        bind:value={current}
        error={failure?.field === 'current' ? failure.text : null}
      />
      <PasswordField
        id="password-new"
        label={m.password_new()}
        autocomplete="new-password"
        bind:value={next}
        error={failure?.field === 'new' ? failure.text : null}
      />
      <PasswordField
        id="password-confirm"
        label={m.setup_confirm()}
        autocomplete="new-password"
        bind:value={confirm}
        error={mismatch ? m.password_mismatch() : null}
      />
      {#if failure && failure.field === null}
        <FormMessage message={failure} retry={() => submit()} />
      {/if}
      {#if changed}
        <p role="status" class="status status-ok self-start">{m.password_changed()}</p>
      {/if}
      <div>
        <button
          type="submit"
          disabled={pending}
          class="rounded-md bg-accent px-4 py-2 text-accent-contrast disabled:opacity-60"
        >
          {m.password_submit()}
        </button>
      </div>
    </form>
  </Card>

  <div>
    <button
      type="button"
      class="rounded-md border border-border px-3 py-2 hover:bg-bg"
      onclick={logout}
    >
      {m.logout()}
    </button>
  </div>
</div>
