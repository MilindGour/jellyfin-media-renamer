<script lang="ts">
	import Icon from 'svelte-awesome';
	import pencil from 'svelte-awesome/icons/pencil';
	import trashIcon from 'svelte-awesome/icons/trash';
	import plusIcon from 'svelte-awesome/icons/plus';
	import yesIcon from 'svelte-awesome/icons/check';
	import noIcon from 'svelte-awesome/icons/remove';

	import type {
		AllowedExtensions,
		DirEntry,
		RenameEntry,
		RenameMediaResponseItem
	} from '$lib/models/models';
	import {
		convertToSizeString,
		getBasename,
		getFiletype,
		getRelativePath,
		getSeasonEpisodeShortString
	} from '$lib/stores/util';
	import Accordion from '../accordion/accordion.svelte';
	import Button from '../button/button.svelte';
	import { Log } from '$lib/services/logger';
	import { ToastService, ToastFactory } from '$lib/components/toast';
	import EpisodeInput from '../episodeInput/episodeInput.svelte';

	const log = new Log('RenameSelectionListItem');
	const toastService = new ToastService();

	let {
		item,
		allowedExtensions
	}: { item: RenameMediaResponseItem; allowedExtensions: AllowedExtensions } = $props();

	const renameAccordionTitle = $derived(
		`Files selected for rename (${item.selected?.length || 0})`
	);
	const ignoredAccordionTitle = $derived(`Ignored files (${item.ignored?.length || 0})`);

	function moveToIgnored(selectedItem: RenameEntry, onlySubtitle = false) {
		if (onlySubtitle) {
			item.ignored = [...item.ignored, selectedItem.subtitle!];
			selectedItem.subtitle = undefined;
		} else {
			item.ignored = [
				...item.ignored,
				...[selectedItem.media, selectedItem.subtitle].filter((x) => !!x)
			];
			item.selected = item.selected.filter((i) => i !== selectedItem);
		}
	}

	function moveToSelected(ignoredItem: DirEntry) {
		const ft = getFiletype(ignoredItem.path, allowedExtensions);

		if (item.type === 'MOVIE') {
			let success = true;
			const mediaPresent = item.selected.length > 0;
			const subtitlePresent = mediaPresent && !!item.selected[0].subtitle;
			const movingSubtitle = ft === 'SUBTITLE';
			const movingMedia = ft === 'MEDIA';

			if (mediaPresent && movingSubtitle) {
				if (subtitlePresent) {
					toastService.show(
						ToastFactory.createWarningToast(
							'Cannot move the subtitle because there is another subtitle selected. Please remove that first before adding this one.'
						)
					);
					success = false;
				} else {
					item.selected[0].subtitle = ignoredItem;
				}
			} else if (!mediaPresent && movingMedia) {
				item.selected = [{ media: ignoredItem }];
			} else {
				success = false;
			}

			if (success) {
				item.ignored = item.ignored.filter((x) => x !== ignoredItem);
			}
		} else {
			toastService.show(
				ToastFactory.createErrorToast('Operation not supported for type ' + item.type)
			);
		}
	}
</script>

<section class="rename-selection-list-item flex flex-col gap-6 rounded border border-gray-200 p-3">
	<section class="media-searchlist-item flex items-center gap-2 rounded">
		<section class="info flex gap-2">
			<img
				class="aspect-[2/3] max-h-16 rounded"
				src={item.info.thumbnailUrl || 'noimg.svg'}
				alt="poster"
			/>
			<section class="text">
				<h3 class="text-md font-medium">{item.info.name} ({item.info.yearOfRelease})</h3>
				<p class="line-clamp-2 text-sm text-gray-500">
					{getBasename(item.entry.path)}
				</p>
			</section>
		</section>
	</section>
	<section class="renames-and-ignores flex flex-col gap-4">
		<!-- Selected files Accordion -->
		<Accordion title={renameAccordionTitle} open>
			{#snippet body()}
				<section class="accordion-body flex flex-col gap-2">
					{#each item.selected as selectedItem}
						<section class="rename-item flex items-center gap-1 rounded bg-gray-50 p-3">
							<div class="first-column min-w-20">
								{#if typeof selectedItem.season === 'number' && typeof selectedItem.episode === 'number'}
									<div class="episode-info">
										{getSeasonEpisodeShortString(selectedItem.season, selectedItem.episode)}
									</div>
								{/if}
								<div class="filesize text-xs">({convertToSizeString(selectedItem.media.size)})</div>
							</div>
							<div class="second-column grow-1">
								<div class="media-name">
									{getRelativePath(selectedItem.media.path, item.entry.path)}
								</div>
								<div class="subtitle-name items-basline gap-1 text-xs text-gray-500">
									Subtitle:
									{#if selectedItem.subtitle}
										{getRelativePath(selectedItem.subtitle.path, item.entry.path)}
										<Button
											type="mini-icon"
											title="Move subtitle to ignored"
											onclick={() => moveToIgnored(selectedItem, true)}
											><Icon data={trashIcon} class="ml-1" /></Button
										>
									{:else}
										<span class="text-red-300">No subtitle</span>
									{/if}
								</div>
							</div>
							<div class="third-column relative flex gap-4">
								<Button
									type="mini-icon"
									title="Move to ignored"
									onclick={() => moveToIgnored(selectedItem)}
								>
									<Icon data={trashIcon} />
								</Button>
								{#if item.type !== 'MOVIE'}
									<Button type="mini-icon" title="Edit">
										<Icon data={pencil} />
									</Button>
								{/if}
							</div>
						</section>
					{:else}
						<i>No selected items.</i>
					{/each}
				</section>
			{/snippet}
		</Accordion>

		<!-- Ignored files Accordion -->
		<Accordion title={ignoredAccordionTitle}>
			{#snippet body()}
				<section class="accordion-body flex flex-col gap-2">
					{#each item.ignored as ignoredItem}
						<section class="rename-item flex items-center gap-1 rounded bg-gray-50 p-3">
							<div class="first-column min-w-20">
								<div class="filesize text-xs">({convertToSizeString(ignoredItem.size)})</div>
							</div>
							<div class="second-column grow-1">
								<div class="media-name">
									{getRelativePath(ignoredItem.path, item.entry.path)}
								</div>
							</div>
							<div class="third-column">
								<Button
									type="mini-icon"
									title="Add to selected"
									onclick={() => moveToSelected(ignoredItem)}
								>
									<Icon data={plusIcon} />
								</Button>
							</div>
						</section>
					{/each}
				</section>
			{/snippet}
		</Accordion>
	</section>
</section>
