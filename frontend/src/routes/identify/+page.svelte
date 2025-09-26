<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, SourceDirectoryInfoItem } from '$lib/components';
	import { API } from '$lib/services/api';
	import { Log } from '$lib/services/logger';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { onMount } from 'svelte';

	const app = new JmrApplicationStore(API.http());
	const log = new Log('Identify page');
	const nextButtonDisabled = $derived(
		app.sourceDirsWithMediaInfo.some((x) => !x.identifiedMediaId)
	);

	onMount(() => {
		// check if data is present, otherwise navigate to home page.
		if (!app.sourceDirsWithMediaInfo.length) {
			setTimeout(() => {
				goto('/');
			}, 3000);
		}
	});

	function searchMediaInfoClickHandler() {
		app.searchMediaInfoProvider();
	}
	function restartClickHandler() {
		window.location.assign('/');
	}
	function nextButtonClickHandler() {
		app.getMediaSelectionForRenames();
	}

	$effect(() => {
		// Check if sourceDirsWithMediaInfo is filled, if yes
		// navigate to next page.
		if (app.mediaSelectionForRenames?.length > 0) {
			goto('/rename');
		}
	});
</script>

<svelte:head>
	<title>Identify: JMR</title>
</svelte:head>

<div class="identify-page flex flex-col gap-4">
	<section class="cta text-right">
		<Button onclick={searchMediaInfoClickHandler}>Search Media Info</Button>
		<Button onclick={restartClickHandler}>Restart</Button>
	</section>
	<section class="source-directory-info-list flex flex-col gap-6">
		{#each app.sourceDirsWithMediaInfo as _, itemIndex}
			<SourceDirectoryInfoItem bind:item={app.sourceDirsWithMediaInfo[itemIndex]} />
		{:else}
			Session refreshed. Redirecting to home page...
		{/each}
	</section>
	<section class="cta">
		<Button disabled={nextButtonDisabled} type="primary" onclick={nextButtonClickHandler}
			>Next</Button
		>
	</section>
</div>
