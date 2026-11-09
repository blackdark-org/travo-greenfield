<script lang="ts">
import Card from '../ui/Card.svelte';

let { title, body }: { title: string; body: string } = $props();
</script>

<Card {title}>
  <p class="text-muted">{body}</p>
</Card>
