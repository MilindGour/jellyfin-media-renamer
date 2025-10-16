<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components';
	import RenameSelectionListItem from '$lib/components/renameSelectionList/renameSelectionListItem.svelte';
	import { API } from '$lib/services/api';
	import { Log } from '$lib/services/logger';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { onMount } from 'svelte';

	const app = new JmrApplicationStore(API.http());
	const nextButtonDisabled = false;

	onMount(() => {
		// check if data is present, otherwise navigate to home page.
		if (!app.sourceDirsWithMediaInfo.length) {
			setTimeout(() => {
				goto('/');
			}, 3000);
		}
	});
	function restartClickHandler() {
		window.location.assign('/');
	}

	function nextButtonClickHandler() {
		app.getMediaConfirmPreview();
	}
</script>

<svelte:head>
	<title>Renames: JMR</title>
</svelte:head>

<section class="page flex flex-col gap-4">
	<section class="cta text-right">
		<Button onclick={restartClickHandler}>Restart</Button>
	</section>
	<h1>Please adjust any selections that are not correct in following list</h1>
	<section class="rename-selection-list mt-3 flex flex-col gap-4">
		{#each app.mediaSelectionForRenames as _, itemIndex}
			<RenameSelectionListItem
				bind:item={app.mediaSelectionForRenames[itemIndex]}
				allowedExtensions={app.config?.allowedExtensions!}
			/>
		{:else}
			Session refreshed. Redirecting to home page...
		{/each}
	</section>
	<section class="cta">
		<Button disabled={nextButtonDisabled} type="primary" onclick={nextButtonClickHandler}
			>Next</Button
		>
	</section>
</section>
