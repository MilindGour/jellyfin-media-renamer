<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { PopupComponent } from '.';
	import { Button } from '$lib/components';
	import type { FileTransferData, ProgressData, WebSocketService } from '$lib/services/network';
	import { Log } from '$lib/services/logger';
	import { formatTimeString, removeCommonSubstring } from '$lib/stores/util';
	import SizeTag from '../sizeTag/sizeTag.svelte';

	const {
		data,
		onclose
	}: { data: { ws: WebSocketService; title?: string }; onclose: (success: boolean) => void } =
		$props();

	const ws = data.ws;
	const log = new Log('CopyStatusPopup');

	let progressStore = $state<ProgressData>([]);
	let autoScrollEnabled = $state<boolean>(true);

	const isComplete = $derived(
		progressStore.length > 0 && progressStore.every((p) => p.percent_complete === 100)
	);
	const completedString = $derived(
		`${progressStore.filter((p) => p.percent_complete === 100).length} / ${progressStore.length}`
	);

	function closePopup(success: boolean) {
		onclose(success);
	}

	onMount(() => {
		ws.connect();
		ws.addListener<ProgressData>('progress', onProgressMessage);
	});

	onDestroy(() => {
		log.info('destroying copy status popup.');
		ws.removeListener('progress', onProgressMessage);
	});

	function onProgressMessage(progressData: ProgressData) {
		log.info('message received by popup:', progressData);
		progressStore = progressData;
		if (autoScrollEnabled) {
			scrollToCurrentFile();
		}
	}

	function scrollToCurrentFile() {
		queueMicrotask(() => {
			const el: HTMLElement | null = document.querySelector('.current-transfer-item');
			if (el) {
				el.scrollIntoView({
					behavior: 'smooth',
					block: 'end'
				});
			} else {
				const allCompleted = document.querySelectorAll('.completed-transfer-item');
				const lastCompleted =
					allCompleted?.length > 0 ? allCompleted.item(allCompleted.length - 1) : null;
				if (lastCompleted) {
					lastCompleted.scrollIntoView({
						behavior: 'smooth',
						block: 'end'
					});
				}
			}
		});
	}
</script>

<PopupComponent title={data.title || '📑 File transfer status'}>
	{#snippet body()}
		<div>
			Completed: {completedString}
		</div>
		<div class="mt-2 flex max-h-100 flex-col gap-2 overflow-auto">
			{#each progressStore as item, itemIndex}
				{@render progressItemComponent(item, itemIndex)}
			{/each}
		</div>
		{#if isComplete}
			<div class="mt-2 text-green-600">
				All files have been transferred successfully! Click close to go back to home page.
			</div>
		{/if}
	{/snippet}
	{#snippet footer()}
		<label class="mr-6 inline-flex items-center gap-2">
			<input type="checkbox" bind:checked={autoScrollEnabled} />
			Auto scroll
		</label>
		<Button disabled={!isComplete} type="primary" onclick={() => closePopup(true)}>Close</Button>
	{/snippet}
</PopupComponent>

{#snippet progressItemComponent(item: FileTransferData, index: number)}
	{@const x = removeCommonSubstring(item.files.old_path, item.files.new_path)}
	<div
		class="progress-file-item relative flex gap-4 rounded-md p-4 {item.percent_complete === 100
			? 'bg-green-100'
			: 'bg-gray-100'}"
		class:current-transfer-item={item.percent_complete > 0 && item.percent_complete < 100}
		class:completed-transfer-item={item.percent_complete === 100}
	>
		<div class="number-col">{index + 1}</div>
		<div class="info-col">
			<div>{x.second}</div>
			<div class="text-sm text-gray-500"><SizeTag bytes={item.total_bytes} /> {x.first}</div>
			{#if item.percent_complete > 0 && item.percent_complete < 100}
				<div class="progres-data mt-2">
					<div
						class="progres-bar absolute inset-0 bg-amber-400 opacity-20"
						style:width={`${item.percent_complete}%`}
					></div>
					<span>{item.transfer_speed}, {formatTimeString(item.time_remaining)} remaining</span>
				</div>
			{/if}
		</div>
	</div>
{/snippet}
