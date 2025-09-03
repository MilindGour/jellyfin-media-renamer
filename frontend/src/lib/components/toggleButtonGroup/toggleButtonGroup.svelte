<script lang="ts">
	const { id, label, options } = $props();

	let activeOption = $state(null);
</script>

<div class="toggle-button-group flex items-center gap-2">
	<div {id} class="option-list flex gap-px">
		{#each options as option}
			<button
				type="button"
				class:bg-primary-bg={activeOption === option}
				class:text-primary-fg={activeOption === option}
				class="min-w-16 cursor-pointer bg-gray-200 px-2 py-1 first:rounded-l-md last:rounded-r-md"
				onclick={() => (activeOption = option)}
			>
				{option}
			</button>
		{/each}
	</div>
</div>
