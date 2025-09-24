<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { Icon } from 'svelte-awesome';
	import closeIcon from 'svelte-awesome/icons/close';
	import {
		slideFromRight,
		ToastType,
		type CloseFunctionHandler,
		type Toast
	} from './toast-models.ts';

	const { toast, onclose }: { toast: Toast; onclose: CloseFunctionHandler } = $props();

	let toastTimeoutID: number | null = null;

	onMount(() => {
		if (toast.closeinMS > 0) {
			startCloseTimer();
		}
	});
	onDestroy(() => {
		if (toastTimeoutID !== null) {
			clearTimeout(toastTimeoutID);
		}
	});

	function startCloseTimer() {
		toastTimeoutID = setTimeout(() => {
			onclose('CLOSE_TIMER');
		}, toast.closeinMS);
	}
	function closeButtonClickHandler() {
		if (toastTimeoutID !== null) {
			clearTimeout(toastTimeoutID);
		}
		onclose('CLOSE_BUTTON');
	}
</script>

<div
	class="toast pointer-events-auto relative rounded border-l-8 bg-amber-200 p-4 shadow-md {toast.type ===
	ToastType.INFO
		? 'border-blue-400'
		: toast.type === ToastType.WARNING
			? 'border-amber-400'
			: 'border-red-400'}"
	transition:slideFromRight
>
	<button class="absolute right-4 cursor-pointer" onclick={closeButtonClickHandler}
		><Icon data={closeIcon} /></button
	>
	<h4 class="pr-4 text-lg font-semibold">{toast.title}</h4>
	<p class="pr-4">{toast.message}</p>
</div>
