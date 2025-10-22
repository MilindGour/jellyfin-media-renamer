<script lang="ts">
	import Icon from 'svelte-awesome';
	import chevronDown from 'svelte-awesome/icons/chevronDown';
	import { slide } from 'svelte/transition';

	let { title, body, open = false } = $props();
</script>

<section class="accordion">
	<button
		onclick={() => (open = !open)}
		class="accordion-header flex cursor-pointer items-center gap-2 text-xl"
	>
		<div class="icon flex items-center transition-transform" class:rotate-180={open}>
			<Icon data={chevronDown} />
		</div>
		<h3>{title}</h3>
	</button>
	{#if open}
		<section class="accordion-body mt-2 pl-0 sm:pl-6" transition:slide>
			{@render body()}
		</section>
	{/if}
</section>
