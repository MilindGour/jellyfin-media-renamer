<script lang="ts">
	import type { DirEntry } from '$lib/models/models';
	import SrcDirListItem from './src-dir-list-item.svelte';

	let {
		name,
		list,
		valid = $bindable(false),
		value = $bindable([])
	}: { name: string; list: DirEntry[]; valid: boolean; value: any[] } = $props();

	let itemStates = $state<{ valid: boolean; value: any }[]>([]);

	$effect(() => {
		if (list?.length > 0) {
			itemStates = new Array(list.length).fill({ value: null, valid: true });
		} else {
			itemStates = [];
		}
	});

	$effect(() => {
		if (itemStates.length > 0) {
			const allValids = itemStates.every((i) => i.valid);
			if (allValids) {
				const nonNullItems = itemStates.filter((i) => i && i.value);
				valid = nonNullItems.length > 0;
				value = nonNullItems;
			} else {
				valid = false;
				value = [];
			}
		} else {
			valid = false;
			value = [];
		}
	});
</script>

<div class="source-directory-list flex flex-col gap-1">
	{#if itemStates.length === list.length}
		{#each list as listItem, index}
			<SrcDirListItem
				entry={listItem}
				{name}
				bind:valid={itemStates[index].valid}
				bind:value={itemStates[index].value}
			/>
		{/each}
	{/if}
</div>
