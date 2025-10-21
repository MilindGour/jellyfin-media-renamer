<script lang="ts">
	import type { DestConfig, MediaType } from '$lib/models';
	import { convertToSizeString } from '$lib/stores/util';
	import Dropdown from '../dropdown/dropdown.svelte';

	let {
		type,
		id,
		value = $bindable(),
		destinationList
	}: { type: MediaType; id: string; value: DestConfig; destinationList: DestConfig[] } = $props();

	const options = $derived(destinationList?.filter((d) => d.type === type));

	function getSizeString(item: DestConfig): string {
		return `${convertToSizeString(item.free_size_kb * 1000)} / ${convertToSizeString(item.total_size_kb * 1000)} available`;
	}
</script>

<Dropdown bind:value {id} {options} labelProp="name" itemTemplate={dropdownTemplate} />

{#snippet dropdownTemplate(item: any)}
	<div class="text-lg">{item.name}</div>
	<p class="text-sm text-gray-500">{getSizeString(item)}</p>
{/snippet}
