<script lang="ts">
	import { Button, Dropdown, SourceDirectoryList } from '$lib/components';
	import type { PageProps } from './$types';
	import { formatPathString } from '$lib/stores/util';
	import type { Source, SourceDirectoryListItemValue } from '$lib/models/models';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { API } from '$lib/services/api';

	const { data }: PageProps = $props();
	const app = new JmrApplicationStore(API.http());

	let source = $state<Source | null>(null);
	let selectedSourceDirectoryItems = $state<SourceDirectoryListItemValue[]>([]);
	const scanDirDisabled = $derived<boolean>(source === null);
	const searchDisabled = $derived<boolean>(
		selectedSourceDirectoryItems.length === 0 || selectedSourceDirectoryItems.some((x) => !x.type)
	);

	async function handleScanDirClick() {
		if (source !== null) {
			app.setSource(source);
		}
	}
	function handleSearchClick() {
		app.setSourceDirectoryListItems(selectedSourceDirectoryItems);
	}
</script>

<section class="page flex flex-col gap-8 pb-16">
	<section
		class="form-section flex flex-col flex-wrap items-stretch gap-2 sm:flex-row sm:items-start"
	>
		<label class="basis-full" for="cfgSourceDD">Please select media source directory</label>
		<Dropdown
			bind:value={source}
			id="sourceDropdown"
			labelProp="name"
			options={data.sourcesResponse.sources}
			itemTemplate={dropdownTemplate}
		/>
		<Button type="primary" disabled={scanDirDisabled} onclick={handleScanDirClick}
			>Scan Directory</Button
		>
	</section>
	<section class="list-section">
		{#await app.sourceDirectories then sourceDirectories}
			{#if sourceDirectories !== null}
				<SourceDirectoryList
					name="selectedList"
					list={sourceDirectories.entries}
					bind:value={selectedSourceDirectoryItems}
				/>
			{:else}
				Select and scan a source to view its directories...
			{/if}
		{/await}
	</section>
	<section class="cta-section flex flex-col items-stretch text-right sm:flex-row sm:items-end">
		<Button type="primary" onclick={handleSearchClick} disabled={searchDisabled}
			>Search Media Online</Button
		>
	</section>
</section>

{#snippet dropdownTemplate(item: Source)}
	<div class="item-instance">
		<div class="font-semibold">{item.name}</div>
		<div class="text-sm text-gray-500">{formatPathString(item.path)}</div>
	</div>
{/snippet}
