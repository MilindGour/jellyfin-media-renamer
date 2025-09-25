<script lang="ts">
	import { PopupComponent } from '.';
	import Button from '../button/button.svelte';

	const { data, onclose } = $props();
</script>

<PopupComponent title={data.title || 'Confirmation'}>
	{#snippet body()}
		<p>{data.bodyText}</p>
	{/snippet}
	{#snippet footer()}
		<Button onclick={() => onclose(true)}>Yes</Button>
		<Button onclick={() => onclose(false)}>No</Button>
	{/snippet}
</PopupComponent>
