<script lang="ts">
import type { AuthMessage } from '../auth-messages';
import { m } from '../paraglide/messages.js';
import { href } from '../router';

// Messages that belong to no single input: shown as an alert with the follow-up the code allows.
let { message, retry }: { message: AuthMessage; retry?: () => void } = $props();
</script>

<div role="alert" class="flex flex-wrap items-center gap-3">
  <span class="status status-danger">{message.text}</span>
  {#if message.action === 'login'}
    <a class="text-accent underline" href={href('login')}>{m.login_title()}</a>
  {:else if message.action === 'setup'}
    <a class="text-accent underline" href={href('setup')}>{m.setup_title()}</a>
  {:else if message.action === 'retry' && retry}
    <button
      type="button"
      class="rounded-md border border-border px-3 py-1 text-sm hover:bg-bg"
      onclick={retry}
    >
      {m.load_retry()}
    </button>
  {/if}
</div>
