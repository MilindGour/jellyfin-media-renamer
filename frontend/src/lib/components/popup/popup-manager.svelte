<script lang="ts">
	import { onMount } from 'svelte';
	import { Popup, PopupStore } from './popup-store.svelte';
	import { PopupService } from './popup-service.svelte';

	const store = new PopupStore();
	const service = new PopupService();

	onMount(() => {
		service.registerPopupStore(store);
	});
	function onPopupClose(p: Popup, result: any) {
		store.removePopup(p, result);
	}
</script>

{#if store.popups.length > 0}
	<section class="popup-manager fixed inset-0 bg-[#ffffff30]">
		{#each store.popups as p}
			<p.component data={p.data} onclose={(result: any) => onPopupClose(p, result)}></p.component>
		{/each}
	</section>
{/if}
