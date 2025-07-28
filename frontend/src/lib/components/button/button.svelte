<script lang="ts">
	const { children, type = 'secondary', ...restProps } = $props();
</script>

<button
	class="button cursor-pointer rounded px-3 py-2"
	class:bg-primary-bg={type === 'primary'}
	class:text-primary-fg={type === 'primary'}
	{...restProps}
>
	{@render children?.()}
</button>
