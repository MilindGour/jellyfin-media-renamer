<script lang="ts">
	import { onMount } from 'svelte';
	import { slide } from 'svelte/transition';
	const { title = '', message = '', body, footer, data = {} } = $props();
	let popupEl: HTMLElement;

	onMount(() => {
		const allFocusables = popupEl.querySelectorAll('input, button');
		if (allFocusables.length > 0) {
			// focus first element
			(allFocusables[0] as HTMLInputElement).focus();
		}
	});
</script>

<section
	class="backdrop fixed inset-0 flex items-start justify-center bg-[#ffffffaa]"
	bind:this={popupEl}
>
	<section
		class="popup-component mt-24 min-w-lg rounded-xl border border-gray-200 bg-white px-4 pt-8 pb-4 shadow-xl"
		transition:slide
	>
		{#if title.length > 0}
			<h2 class="mb-4 text-xl font-medium">{title}</h2>
		{/if}

		{#if body}
			{@render body(data)}
		{:else if message?.length > 0}
			<p>{message}</p>
		{/if}

		<section class="footer mt-8 text-right">
			{#if footer}
				{@render footer()}
			{/if}
		</section>
	</section>
</section>
