<script lang="ts">
	import type { ConfigSourceByID, DirEntry } from '$lib/models/config-models';
	import { CSDLStore } from './csdl-store.svelte';
	import { convertToSizeString, joinStrings } from '$lib/stores/util';

	const { id, data }: { id: string; data: ConfigSourceByID | null } = $props();
	let store: CSDLStore | null = $state(null);

	function handleListCheckChange(e: Event, item: DirEntry) {
		const checked = (e.target as HTMLInputElement).checked;
		if (checked) {
			store?.selectItem(item);
		} else {
			store?.unselectItem(item);
		}
	}
	function handleSelectAllChange(e: Event) {
		const checked = (e.target as HTMLInputElement).checked;
		if (checked) {
			store?.selectAll();
		} else {
			store?.unselectAll();
		}
	}

	$effect(() => {
		if (data && data?.directoryEntries !== null) {
			store = new CSDLStore(data.directoryEntries);
		}
	});
</script>

<div class="config-source-list list_{id} flex flex-col gap-3">
	{#if data !== null}
		<section class="list-actions-section">
			<div class="flex items-center gap-3 px-3">
				<input
					checked={store?.selectAllChecked}
					indeterminate={store?.selectAllIndeterminate}
					type="checkbox"
					id={joinStrings(id, 'cbSelectAll')}
					onchange={handleSelectAllChange}
				/>
				<label for={joinStrings(id, 'cbSelectAll')}>Select All</label>
			</div>
		</section>
		<section class="list flex flex-col gap-1">
			{#each data.directoryEntries as item (item.id)}
				<label
					for={joinStrings('checkbox', id, item.id.toString())}
					class="config-source-details-item flex cursor-pointer items-center gap-3 rounded bg-gray-50 p-3 hover:bg-gray-100"
				>
					<div class="checkbox-wrapper">
						<input
							onchange={(e: Event) => handleListCheckChange(e, item)}
							checked={store?.isItemSelected(item)}
							type="checkbox"
							id={joinStrings('checkbox', id, item.id.toString())}
						/>
					</div>
					<div class="title-wrapper">
						<p class="text-lg">{item.name}</p>
						<p class="text-sm text-gray-500">{convertToSizeString(item.size)}</p>
					</div>
				</label>
			{/each}
		</section>
	{:else}
		Please select some input source.
	{/if}
</div>
