<script lang="ts">
	const { children, type = 'secondary', disabled = false, ...restProps } = $props();
</script>

<button
	class="button cursor-pointer rounded disabled:pointer-events-none disabled:bg-gray-200 disabled:text-gray-50 {type ===
	'primary'
		? 'min-w-32 px-3 py-2'
		: type === 'mini-icon'
			? 'aspect-square rounded-full p-1'
			: 'px-3 py-2 hover:bg-gray-100'}"
	class:bg-primary-bg={type === 'primary'}
	class:text-primary-fg={type === 'primary'}
	{disabled}
	{...restProps}
>
	{@render children?.()}
</button>
