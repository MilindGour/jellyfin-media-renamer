<script lang="ts">
	import { PopupComponent } from '.';
	import Button from '../button/button.svelte';

	const { data, onclose } = $props();
	let season = $state<number>(data.season);
	let episode = $state<number>(data.episode);

	function closePopup(success: boolean) {
		onclose(success ? { season, episode } : null);
	}
</script>

<PopupComponent title={data.title || 'Confirmation'}>
	{#snippet body()}
		<div class="form flex gap-2">
			<fieldset class="flex w-fit items-center overflow-hidden rounded border border-gray-200">
				<label for="season" class="bg-gray-200 p-2">Season</label>
				<input type="number" id="season" class="w-[8ch] border-0" bind:value={season} />
			</fieldset>
			<fieldset class="flex w-fit items-center overflow-hidden rounded border border-gray-200">
				<label for="episode" class="bg-gray-200 p-2">Episode</label>
				<input type="number" id="episode" class="w-[8ch] border-0" bind:value={episode} />
			</fieldset>
		</div>
	{/snippet}
	{#snippet footer()}
		<Button onclick={() => closePopup(false)}>Cancel</Button>
		<Button onclick={() => closePopup(true)}>Save</Button>
	{/snippet}
</PopupComponent>
