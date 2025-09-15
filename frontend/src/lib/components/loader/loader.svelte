<script lang="ts">
	import { loaderService } from './loader-store.svelte';

	const DELAY_BEFORE_LOADER_MS = 500;

	let loaderVisible = $state<boolean>(false);

	$effect(() => {
		if (loaderService.visible === false) {
			loaderVisible = false;
		} else {
			setTimeout(() => {
				// check if it is still true
				loaderVisible = loaderService.visible;
			}, DELAY_BEFORE_LOADER_MS);
		}
	});
</script>

{#if loaderVisible}
	<div class="loader-wrapper fixed inset-0 flex items-center justify-center bg-white opacity-50">
		<div class="loader"></div>
	</div>
{/if}

<style>
	.loader {
		width: 50px;
		padding: 8px;
		aspect-ratio: 1;
		border-radius: 50%;
		background: #000000;
		--_m: conic-gradient(#0000 10%, #000), linear-gradient(#000 0 0) content-box;
		-webkit-mask: var(--_m);
		mask: var(--_m);
		-webkit-mask-composite: source-out;
		mask-composite: subtract;
		animation: l3 1s infinite linear;
	}
	@keyframes l3 {
		to {
			transform: rotate(1turn);
		}
	}
</style>
