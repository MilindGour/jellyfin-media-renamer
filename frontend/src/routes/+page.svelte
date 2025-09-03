<script lang="ts">
	import { Button, Dropdown, SourceDirectoryList } from '$lib/components';
	import type { PageProps } from './$types';
	import { formatPathString } from '$lib/stores/util';
	import type { Source } from '$lib/models/models';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { API } from '$lib/services/api';

	const { data }: PageProps = $props();
	let sourceValue = $state<Source | null>(null);
	let sourceValid = $state(false);
	const app = new JmrApplicationStore(API.http());

	async function handleScanDirClick() {
		if (sourceValue !== null) {
			app.setSource(sourceValue);
		}
	}
	function handleSearchClick() {
		// get the value of csdl component using the service
	}
</script>

<section class="page flex flex-col gap-8">
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
			{console.log(sourceDirectories)}
			{#if sourceDirectories !== null}
				<SourceDirectoryList />
				<section class="cta-section">
					<Button type="primary" onclick={handleSearchClick}>Search Media Online</Button>
				</section>
			{:else}
				Select and scan a source to view its directories...
			{/if}
		{/await}
	</section>
</section>

{#snippet dropdownTemplate(item: Source)}
	<div class="item-instance">
		<div class="font-semibold">{item.name}</div>
		<div class="text-sm text-gray-500">{formatPathString(item.path)}</div>
	</div>
{/snippet}
