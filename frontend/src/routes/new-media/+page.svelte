<script lang="ts">
	import { Button, NewMediaSearchItemComponent } from '$lib/components';
	import { PopupService } from '$lib/components/popup';
	import { ToastBuilder, ToastService } from '$lib/components/toast';
	import { ToastType } from '$lib/components/toast/toast-models';
	import { HttpService } from '$lib/services/network';
	import { NewMediaStore, type NewMediaSearchItem } from '$lib/stores/new-media-store.svelte';

	const nmStore = new NewMediaStore(new HttpService());
	const ps = new PopupService();
	const ts = new ToastService();
	let searchFieldText = $state<string>('');

	function newMediaSearchItemClickHandler(item: NewMediaSearchItem) {
		ps.showConfirmation(item.name, 'Add following to download queue?').then(async (yes) => {
			if (yes) {
				const addSuccessful = await nmStore.addItemToDownloadQueue(item);
				if (addSuccessful) {
					ts.show(
						new ToastBuilder('Added the download successfully').setType(ToastType.SUCCESS).build()
					);
				} else {
					ts.show(
						new ToastBuilder('Unknown error occured while adding to download queue.')
							.setType(ToastType.ERROR)
							.build()
					);
				}
			}
		});
	}

	function searchClickHandler() {
		if (searchFieldText.length > 0) {
			nmStore.getSearchResultsForTerm(searchFieldText);
		}
	}
</script>

<section class="page flex flex-col gap-8 pb-16">
	<h1 class="text-lg">Add new media item in the downloads.</h1>
	<div class="form-control flex gap-1">
		<input
			type="text"
			class="max-w-md flex-2 rounded"
			placeholder="Enter search term"
			id="searchTerm"
			bind:value={searchFieldText}
			autocomplete="off"
		/>
		<Button type="primary" onclick={searchClickHandler}>Search</Button>
	</div>

	{#if nmStore.items.length > 0}
		<div class="search-results-total text-lg">
			Total search results: <strong>{nmStore.items.length}</strong>
		</div>
		<div class="new-media-search-list flex flex-col gap-2">
			{#each nmStore.items as item (item.id)}
				<NewMediaSearchItemComponent {item} onclick={newMediaSearchItemClickHandler} />
			{/each}
		</div>
	{/if}
</section>
