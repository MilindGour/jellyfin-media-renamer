<script lang="ts">
	import { goto } from '$app/navigation';
	import RenameSelectionListItem from '$lib/components/renameSelectionList/renameSelectionListItem.svelte';
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
</script>

<svelte:head>
	<title>Renames: JMR</title>
</svelte:head>

<section class="page">
	<h1>Please adjust any selections that are not correct in following list</h1>
	<section class="rename-selection-list mt-3 flex flex-col gap-4">
		{#each app.mediaSelectionForRenames as item}
			<RenameSelectionListItem {item} allowedExtensions={app.config?.allowedExtensions!} />
		{:else}
			Session refreshed. Redirecting to home page...
		{/each}
	</section>
</section>
