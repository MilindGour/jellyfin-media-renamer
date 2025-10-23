<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components';
	import { PopupService } from '$lib/components/popup';
	import { RenameSelectionListItem } from '$lib/components';
	import { API } from '$lib/services/api';
	import { Log } from '$lib/services/logger';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { onMount } from 'svelte';

	const { data } = $props();
	const log = new Log('Rename Page');
	const app = new JmrApplicationStore(API.http());
	const popupService = new PopupService();
	const confirmAndSyncDisabled = $derived(app.mediaDestinationSelections.some((ds) => ds === null));

	onMount(() => {
		// check if data is present, otherwise navigate to home page.
		if (!app.sourceDirsWithMediaInfo.length) {
			setTimeout(() => {
				goto('/');
			}, 3000);
		}

		log.info('Destinations:', data.destinations);
	});
	function restartClickHandler() {
		window.location.assign('/');
	}

	async function confirmRenameClickHandler() {
		const userConfirmed = await popupService.showConfirmation(
			'Are you sure you want to go ahead with the renaming?'
		);
		if (userConfirmed) {
			await app.confirmMediaRequest();
			queueMicrotask(async () => {
				await popupService.showFileTransferStatusPopup();
				window.location.assign('/');
			});
		}
	}
</script>

<svelte:head>
	<title>Renames: JMR</title>
</svelte:head>

<section class="page flex flex-col gap-4 pb-16">
	<section class="cta text-right">
		<Button onclick={restartClickHandler}>Restart</Button>
	</section>
	<h1>Please adjust any selections that are not correct in following list</h1>
	<section class="rename-selection-list mt-3 flex flex-col gap-4">
		{#each app.mediaSelectionForRenames as _, itemIndex}
			<RenameSelectionListItem
				bind:item={app.mediaSelectionForRenames[itemIndex]}
				bind:dest={app.mediaDestinationSelections[itemIndex]}
				allowedExtensions={app.config?.allowedExtensions!}
				destinations={data?.destinations}
			/>
		{:else}
			Session refreshed. Redirecting to home page...
		{/each}
	</section>
	<section class="cta flex flex-col md:flex-row">
		<Button disabled={confirmAndSyncDisabled} type="primary" onclick={confirmRenameClickHandler}
			>Confirm Rename & Sync</Button
		>
	</section>
</section>
