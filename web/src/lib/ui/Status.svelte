<script lang="ts">
import type { Snippet } from 'svelte';
import { type StatusKind, statusClass } from './status';

let { kind, children }: { kind: StatusKind; children: Snippet } = $props();
</script>

<span class={statusClass(kind)}>{@render children()}</span>
