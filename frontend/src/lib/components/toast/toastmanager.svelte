<script lang="ts">
	import { onMount } from 'svelte';
	import type { Toast, ToastCloseReason } from './toast-models.ts';
	import ToastComponent from './toast.svelte';
	import { ToastService } from './toast-service.svelte';
	import { ToastManagerStore } from './toast-mgr-store.svelte';
	import { Log } from '$lib/services/logger';

	const { id } = $props();
	const store = new ToastManagerStore(id);
	const service = new ToastService();
	const log = new Log('ToastManager');

	onMount(() => {
		service.registerManager(store);
	});

	function toastOnCloseHandler(toast: Toast, closeReason: ToastCloseReason) {
		log.info('Toast closed. Reason:', closeReason);
		store.removeToast(toast);
	}
</script>

<div
	class="toast-manager pointer-events-none fixed top-0 right-0 bottom-0 flex w-full flex-col gap-2 p-3 sm:w-[50vw] lg:w-[30vw]"
	{id}
>
	{#each store.toasts as toast (toast.id)}
		<ToastComponent {toast} onclose={(closeReason) => toastOnCloseHandler(toast, closeReason)} />
	{/each}
</div>
