<script lang="ts">
	import JmrHeader from '$lib/jmrHeader.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<div class="page-container mx-auto max-w-5xl">
	<JmrHeader />
	{@render children()}
</div>
