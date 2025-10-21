<script lang="ts">
	import type { DestConfig, MediaType } from '$lib/models';
	import { convertToSizeString } from '$lib/stores/util';
	import Dropdown from '../dropdown/dropdown.svelte';

	let {
		type,
		id,
		value = $bindable(),
		sourceSize,
		destinationList
	}: {
		type: MediaType;
		id: string;
		value: DestConfig;
		sourceSize: number;
		destinationList: DestConfig[];
	} = $props();

	const options = $derived(destinationList?.filter((d) => d.type === type));

	function getSizeString(item: DestConfig): string {
		return `${convertToSizeString(item.free_size_kb * 1000)} / ${convertToSizeString(item.total_size_kb * 1000)} available`;
	}

	function isOptionDisabled(destConfig: DestConfig, sourceSize: number): boolean {
		const destSizeInBytes = destConfig.free_size_kb * 1000;
		const result = destSizeInBytes <= sourceSize;
		return result;
	}
</script>

<Dropdown
	bind:value
	{id}
	{options}
	disabledItemFn={(i: DestConfig) => isOptionDisabled(i, sourceSize)}
	labelProp="name"
	itemTemplate={dropdownTemplate}
/>

{#snippet dropdownTemplate(item: DestConfig)}
	<div class="text-lg">{item.name}</div>
	<p class="text-sm font-semibold text-gray-500">{getSizeString(item)}</p>
	<p class="text-sm text-gray-500">{item.path}</p>
{/snippet}
