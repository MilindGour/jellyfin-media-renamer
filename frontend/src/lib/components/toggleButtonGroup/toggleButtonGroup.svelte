<script lang="ts">
	let {
		id,
		options,
		required = false,
		valid = $bindable(!required),
		value = $bindable(null)
	} = $props();

	$effect(() => {
		if (value?.length > 0) {
			valid = true;
		} else {
			valid = !required;
		}
	});
</script>

<div class="toggle-button-group flex items-center gap-2">
	<div {id} class="option-list flex gap-px">
		{#each options as option}
			<button
				type="button"
				class:bg-primary-bg={value === option}
				class:text-primary-fg={value === option}
				class="min-w-16 cursor-pointer bg-gray-200 px-2 py-1 first:rounded-l-md last:rounded-r-md"
				onclick={() => (value = option)}
			>
				{option}
			</button>
		{/each}
	</div>
</div>
