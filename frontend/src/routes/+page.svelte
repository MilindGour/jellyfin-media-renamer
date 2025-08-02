<script lang="ts">
	import { Dropdown, DropdownService } from '$lib/components/dropdown';
	import { Button } from '$lib/components/button';
	import type { PageProps } from './$types';
	import { onMount } from 'svelte';
	import type { ConfigSource, DirEntry } from '$lib/models/config-models';
	import { JNetworkClient } from '$lib/services/network';
	import { JmrStore } from '$lib/stores/app-store.svelte';

	const { data }: PageProps = $props();
	const netClient: JNetworkClient = new JNetworkClient();

	const jmrStore = new JmrStore(netClient);

	onMount(() => {
		console.log('[dbg] load data:', data);
	});
	function formatPath(path: string): string {
		if (path.startsWith('/')) {
			path = path.substring(1);
		}
		return path.split('/').join(' > ');
	}

	async function handleScanDirClick() {
		const cs = DropdownService.getValueOf<ConfigSource>('cfgSourceDD');
		jmrStore.setConfigSource(cs);
	}
</script>

<section class="page flex flex-col gap-8">
	<section
		class="form-section flex flex-col flex-wrap items-stretch gap-2 sm:flex-row sm:items-start"
	>
		<label class="basis-full" for="cfgSourceDD">Please select media source directory</label>
		<Dropdown
			id="cfgSourceDD"
			labelProp="name"
			options={data.configSources}
			itemTemplate={dropdownTemplate}
		/>
		<Button type="primary" onclick={handleScanDirClick}>Scan Directory</Button>
	</section>
	<section class="list-section">
		{#await jmrStore.configSourceDetails}
			Loading subdirectories of selected source...
		{:then cfd}
			{#if cfd !== null}
				{#each cfd.directoryEntries as cfdItem (cfdItem.id)}
					{@render configSourceDetailListItem(cfdItem)}
				{/each}
			{:else}
				Please select a media source directory.
			{/if}
		{/await}
	</section>
</section>

{#snippet dropdownTemplate(item: ConfigSource)}
	<div class="item-instance">
		<div class="font-semibold">{item.name}</div>
		<div class="text-sm text-gray-500">{formatPath(item.path)}</div>
	</div>
{/snippet}

{#snippet configSourceDetailListItem(item: DirEntry)}
	<div class="config-source-details-item">
		{item.name}, {item.size}, {item.path}
	</div>
{/snippet}
