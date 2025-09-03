<script lang="ts">
	import { Button, Dropdown, SourceDirectoryList } from '$lib/components';
	import type { PageProps } from './$types';
	import { formatPathString } from '$lib/stores/util';
	import type { Source } from '$lib/models/models';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { API } from '$lib/services/api';

	const { data }: PageProps = $props();
	const app = new JmrApplicationStore(API.http());

	let sourceValue = $state<Source | null>(null);
	let sourceValid = $state(false);
	let srcDirsValid = $state(false);
	let srcDirsValue = $state([]);

	async function handleScanDirClick() {
		if (sourceValue !== null) {
			app.setSource(sourceValue);
		}
	}
	function handleSearchClick() {
		// get the value of csdl component using the service
	}
</script>

<section class="page flex flex-col gap-8 pb-16">
	<section
		class="form-section flex flex-col flex-wrap items-stretch gap-2 sm:flex-row sm:items-start"
	>
		<label class="basis-full" for="cfgSourceDD">Please select media source directory</label>
		<Dropdown
			bind:value={sourceValue}
			bind:valid={sourceValid}
			required
			id="sourceDropdown"
			labelProp="name"
			options={data.sourcesResponse.sources}
			itemTemplate={dropdownTemplate}
		/>
		<Button type="primary" disabled={!sourceValid} onclick={handleScanDirClick}
			>Scan Directory</Button
		>
	</section>
	<section class="list-section">
		{#await app.sourceDirectories then sourceDirectories}
			{#if sourceDirectories !== null}
				<SourceDirectoryList
					name="selectedList"
					list={sourceDirectories.entries}
					bind:valid={srcDirsValid}
					bind:value={srcDirsValue}
				/>
			{:else}
				Select and scan a source to view its directories...
			{/if}
		{/await}
	</section>
	<section class="cta-section flex flex-col items-stretch text-right sm:flex-row sm:items-end">
		<Button type="primary" onclick={handleSearchClick} disabled={!srcDirsValid}
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
