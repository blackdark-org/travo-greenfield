<script lang="ts">
let {
  id,
  label,
  value = $bindable(''),
  autocomplete,
  error = null,
}: {
  id: string;
  label: string;
  value?: string;
  autocomplete: 'current-password' | 'new-password';
  error?: string | null;
} = $props();
</script>

<div class="flex flex-col gap-1">
  <label for={id} class="text-sm font-medium">{label}</label>
  <input
    {id}
    type="password"
    name={id}
    {autocomplete}
    required
    bind:value
    aria-invalid={error ? 'true' : undefined}
    aria-describedby={error ? `${id}-error` : undefined}
    class="rounded-md border border-border bg-bg px-3 py-2"
  />
  {#if error}
    <p id="{id}-error" class="text-sm text-(--danger)">{error}</p>
  {/if}
</div>
