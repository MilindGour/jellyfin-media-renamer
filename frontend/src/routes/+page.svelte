<script lang="ts">
	import { Button, Dropdown, SourceDirectoryList } from '$lib/components';
	import type { PageProps } from './$types';
	import { formatPathString } from '$lib/stores/util';
	import type { Source, SourceDirectory } from '$lib/models';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { API } from '$lib/services/api';
	import { goto } from '$app/navigation';
	import { Log } from '$lib/services/logger';
	import { onMount } from 'svelte';

	const { data }: PageProps = $props();
	const app = new JmrApplicationStore(API.http());
	const log = new Log('Select Source Page');

	let source = $state<Source | null>(null);
	let selectedSourceDirectoryItems = $state<SourceDirectory[]>([]);
	const scanDirDisabled = $derived<boolean>(source === null);
	const searchDisabled = $derived<boolean>(
		selectedSourceDirectoryItems.length === 0 || selectedSourceDirectoryItems.some((x) => !x.type)
	);

	onMount(() => {
		app.setConfig(data.appConfig);
	});

	async function handleScanDirClick() {
		if (source !== null) {
			app.setSource(source);
		} else {
			log.error('[JMR] source is null');
		}
	}
	function handleNextButtonClick() {
		app.setSourceDirectories(selectedSourceDirectoryItems);
	}

	$effect(() => {
		// Check if sourceDirsWithMediaInfo is filled, if yes
		// navigate to next page.
		if (app.sourceDirsWithMediaInfo?.length > 0) {
			goto('/identify');
		}
	});
</script>

<svelte:head>
	<title>Select Source: JMR</title>
</svelte:head>
<section class="page flex flex-col gap-8 pb-16">
	<section
		class="form-section flex flex-col flex-wrap items-stretch gap-2 sm:flex-row sm:items-start"
	>
		<label class="basis-full" for="cfgSourceDD">Please select media source directory</label>
		<Dropdown
			bind:value={source}
			id="sourceDropdown"
			labelProp="name"
			options={data.appConfig.source}
			itemTemplate={dropdownTemplate}
		/>
		<Button type="primary" disabled={scanDirDisabled} onclick={handleScanDirClick}
			>Scan Directory</Button
		>
	</section>
	<section class="list-section">
		{#if app.sourceDirectories !== null && app.sourceDirectories.entries?.length > 0}
			<SourceDirectoryList
				name="selectedList"
				list={app.sourceDirectories.entries}
				bind:value={selectedSourceDirectoryItems}
			/>
		{:else if app.sourceDirectories?.entries === null}
			There are no files or directories in the selected source. Please select a different source.
		{:else}
			Select and scan a source to view its directories.
		{/if}
	</section>
	<section
		class="cta-section flex flex-col justify-stretch text-right sm:flex-row sm:justify-start"
	>
		<Button type="primary" onclick={handleNextButtonClick} disabled={searchDisabled}>Next</Button>
	</section>
</section>

{#snippet dropdownTemplate(item: Source)}
	<div class="item-instance">
		<div class="font-semibold">{item.name}</div>
		<div class="text-sm text-gray-500">{formatPathString(item.path)}</div>
	</div>
{/snippet}
