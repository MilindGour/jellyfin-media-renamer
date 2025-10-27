<script lang="ts">
	import type { NewMediaSearchItem } from '$lib/stores/new-media-store.svelte';
	import { Button } from '..';
	import SizeTag from '../sizeTag/sizeTag.svelte';
	import Icon from 'svelte-awesome';
	import downloadIcon from 'svelte-awesome/icons/download';

	const {
		item,
		onclick
	}: { item: NewMediaSearchItem; onclick: (item: NewMediaSearchItem) => void } = $props();
</script>

<section
	class="new-media-search-item flex items-center justify-between rounded-lg border border-gray-200 bg-gray-100 p-3"
>
	<div class="left-column">
		<p class="title text-lg break-all">{item.name}</p>
		<SizeTag bytes={+item.size} />
		<span
			>S/L: {item.seeders} | {item.leechers}
			<em class="text-gray-500">- {item.username}</em></span
		>
	</div>
	<div class="right-column">
		<Button type="icon" title="Add to download queue" onclick={() => onclick(item)}>
			<Icon data={downloadIcon} />
		</Button>
	</div>
</section>
