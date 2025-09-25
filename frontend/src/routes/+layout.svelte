<script lang="ts">
	import Loader from '$lib/components/loader/loader.svelte';
	import { PopupManager } from '$lib/components/popup';
	import { ToastManagerComponent } from '$lib/components/toast';
	import JmrHeader from '$lib/jmrHeader.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<div class="page-container mx-auto max-w-5xl px-4">
	<JmrHeader />
	{@render children()}
	<Loader />
	<ToastManagerComponent id="appToastManager" />
	<PopupManager />
</div>
