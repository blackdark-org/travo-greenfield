<script lang="ts">
import { m } from '../paraglide/messages.js';
import { theme } from '../theme.svelte';
import Card from '../ui/Card.svelte';
</script>

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
