<script lang="ts">
	import type { DestConfig, MediaType } from '$lib/models';
	import Dropdown from '../dropdown/dropdown.svelte';

	let {
		type,
		id,
		value = $bindable(),
		destinationList
	}: { type: MediaType; id: string; value: DestConfig; destinationList: DestConfig[] } = $props();

	const options = $derived(destinationList?.filter((d) => d.type === type));
</script>

<Dropdown bind:value {id} {options} labelProp="name" itemTemplate={dropdownTemplate} />

{#snippet dropdownTemplate(item: any)}
	{item.name}
{/snippet}
