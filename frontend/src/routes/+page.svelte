<script lang="ts">
	import { Button } from '$lib/components/button';
	import { Dropdown, DropdownService } from '$lib/components/dropdown';
	import type { ConfigSource } from '$lib/models/config-models';
	import { JNetworkClient } from '$lib/services/network';
	import { JmrStore } from '$lib/stores/app-store.svelte';
	import type { PageProps } from './$types';
	import { formatPathString } from '$lib/stores/util';
	import { ConfigSourceDetailList } from '$lib/components/configSourceDetailList';

	const { data }: PageProps = $props();
	const netClient: JNetworkClient = new JNetworkClient();
	const jmrStore = new JmrStore(netClient);

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
			options={data.configSources.sources}
			itemTemplate={dropdownTemplate}
		/>
		<Button type="primary" onclick={handleScanDirClick}>Scan Directory</Button>
	</section>
	<section class="list-section">
		{#await jmrStore.configSourceDetails}
			Loading subdirectories of selected source...
		{:then cfd}
			<ConfigSourceDetailList id="directoryList" data={cfd} />
		{/await}
	</section>
</section>

{#snippet dropdownTemplate(item: ConfigSource)}
	<div class="item-instance">
		<div class="font-semibold">{item.name}</div>
		<div class="text-sm text-gray-500">{formatPathString(item.path)}</div>
	</div>
{/snippet}
