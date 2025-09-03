<script lang="ts">
	const { children, type = 'secondary', disabled = false, ...restProps } = $props();
</script>

<button
	class="button min-w-32 cursor-pointer rounded px-3 py-2 disabled:pointer-events-none disabled:bg-gray-200 disabled:text-gray-50"
	class:bg-primary-bg={type === 'primary'}
	class:text-primary-fg={type === 'primary'}
	{disabled}
	{...restProps}
>
	{@render children?.()}
</button>
