<script module>
	let _id = 0;

	function getNextId() {
		return (++_id).toString();
	}
</script>

<script lang="ts">
	import type { MediaInfo } from '$lib/models';
	import { joinStrings } from '$lib/stores/util';
	import MediaSearchlistItem from './media-searchlist-item.svelte';

	let { list, value = $bindable<string>() }: { list: MediaInfo[]; value: string } = $props();

	const id = joinStrings('medialist', getNextId());
</script>

<div class="media-search-list flex flex-col gap-2">
	{#if list.length > 0}
		<section class="text">Please select the correct media from the list:</section>
		<section class="list flex flex-col gap-1">
			{#each list as listItem}
				<MediaSearchlistItem
					item={listItem}
					name={id}
					onchange={(e) => (value = e.currentTarget.value)}
				/>
			{/each}
		</section>
	{:else}
		<section class="text">Check the identified info, then click search.</section>
	{/if}
</div>
