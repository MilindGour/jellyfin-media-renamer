<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, SourceDirectoryInfoItem } from '$lib/components';
	import { API } from '$lib/services/api';
	import { JmrApplicationStore } from '$lib/stores/app-store.svelte';
	import { onMount } from 'svelte';

	const app = new JmrApplicationStore(API.http());

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
</script>

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
</div>
