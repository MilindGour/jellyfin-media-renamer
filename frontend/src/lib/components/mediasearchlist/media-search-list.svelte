<script module>
	let _id = 0;

	function getNextId() {
		return (++_id).toString();
	}
</script>

<script lang="ts">
	import type { MediaInfo } from '$lib/models/models';
	import { joinStrings } from '$lib/stores/util';
	import MediaSearchlistItem from './media-searchlist-item.svelte';

	let { value = $bindable<string>() }: { value: string } = $props();

	const id = joinStrings('medialist', getNextId());
	const td0Item: MediaInfo[] = [];
	const td1Item: MediaInfo[] = [
		{
			description: 'Test description',
			mediaId: '1234',
			name: 'First result',
			thumbnailUrl:
				'https://media.themoviedb.org/t/p/w440_and_h660_face/tg9I5pOY4M9CKj8U0cxVBTsm5eh.jpg',
			yearOfRelease: 2025
		}
	];
	const tdNItem: MediaInfo[] = [
		{
			description: 'Test description',
			mediaId: '1234',
			name: 'First result',
			thumbnailUrl:
				'https://media.themoviedb.org/t/p/w440_and_h660_face/tg9I5pOY4M9CKj8U0cxVBTsm5eh.jpg',
			yearOfRelease: 2025
		},
		{
			description: 'Test description 2',
			mediaId: '5678',
			name: 'Second result',
			thumbnailUrl:
				'https://media.themoviedb.org/t/p/w440_and_h660_face/yb4F1Oocq8GfQt6iIuAgYEBokhG.jpg',
			yearOfRelease: 2022
		}
	];

	const itemList = tdNItem;
</script>

<div class="media-search-list flex flex-col gap-2">
	{#if itemList.length > 0}
		<section class="text">Please select the correct media from the list:</section>
		<section class="list flex flex-col gap-1">
			{#each itemList as item, itemIndex}
				<MediaSearchlistItem {item} name={id} onchange={(e) => (value = e.currentTarget.value)} />
			{/each}
		</section>
	{:else}
		<section class="text">0 results found. Please improve the search text above.</section>
	{/if}
</div>
