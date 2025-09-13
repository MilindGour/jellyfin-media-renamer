<script lang="ts">
	import type { DirEntry, MediaType, SourceDirectory } from '$lib/models/models';
	import { Button } from '$lib/components';
	import SrcDirListItem from './src-dir-list-item.svelte';

	let {
		name,
		list,
		value = $bindable([])
	}: {
		name: string;
		list: DirEntry[];
		value: SourceDirectory[];
	} = $props();

	let items = $state<SourceDirectory[]>([]);
	let selectAllChecked = $derived(value.length === items.length);
	let selectAllIndeterminate = $derived(value.length < items.length && value.length > 0);

	$effect(() => {
		if (list?.length > 0) {
			items = list.map((l) => ({ entry: l, selected: false, type: null }) as SourceDirectory);
		} else {
			items = [];
		}
	});

	$effect(() => {
		value = items.filter((i) => i.selected);
	});

	function selectAllChangeHandler(e: Event) {
		const checked = (e.currentTarget as HTMLInputElement).checked;
		items.forEach((i) => (i.selected = checked));
	}
	function allType(newType: MediaType) {
		items.forEach((i) => (i.type = newType));
	}
</script>

<div class="source-directory-list-wrapper">
	<div class="actions flex flex-col justify-between p-3 sm:flex-row">
		<div class="action flex items-center gap-2">
			<input
				type="checkbox"
				id="selectAll"
				indeterminate={selectAllIndeterminate}
				checked={selectAllChecked}
				onchange={selectAllChangeHandler}
			/>
			<label for="selectAll">Select All</label>
		</div>
		<div class="action pt-1">
			<Button onclick={() => allType('MOVIE')}>All Movies</Button>
			<Button onclick={() => allType('TV')}>All TVs</Button>
		</div>
	</div>
	<div class="source-directory-list flex flex-col gap-1">
		{#each items as _, itemIndex}
			<SrcDirListItem {name} bind:value={items[itemIndex]} />
		{/each}
	</div>
</div>
