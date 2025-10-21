<script lang="ts">
	import Icon from 'svelte-awesome';
	import pencil from 'svelte-awesome/icons/pencil';
	import trashIcon from 'svelte-awesome/icons/trash';
	import plusIcon from 'svelte-awesome/icons/plus';
	import { v4 as uuidv4 } from 'uuid';

	import type {
		AllowedExtensions,
		DestConfig,
		DirEntry,
		RenameEntry,
		RenameMediaResponseItem
	} from '$lib/models';
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
	import { PopupService } from '../popup';
	import MediaTag from '../mediaTag/media-tag.svelte';
	import { DestinationDirectoryDropdown } from '$lib/components';

	const log = new Log('RenameSelectionListItem');
	const toastService = new ToastService();
	const popupService = new PopupService();
	const id = `renameItem_${uuidv4()}`;

	let {
		item = $bindable(),
		allowedExtensions,
		dest = $bindable(),
		destinations
	}: {
		item: RenameMediaResponseItem;
		allowedExtensions: AllowedExtensions;
		dest: DestConfig;
		destinations: DestConfig[];
	} = $props();

	const totalSelectedSize = $derived(
		item.selected?.reduce((acc, cur) => acc + cur?.media?.size + (cur?.subtitle?.size || 0), 0)
	);
	const totalSelectedItems = $derived(
		item.selected?.map<number>((i) => (!!i.subtitle ? 2 : 1)).reduce((ac, cr) => ac + cr)
	);
	const renameAccordionTitle = $derived(
		`Files selected for rename [${totalSelectedItems} items totalling ${convertToSizeString(totalSelectedSize)}]`
	);

	const totalIgnoredItems = $derived(item.ignored?.length || 0);
	const totalIgnoredSize = $derived(
		item.ignored?.map((i) => i.size).reduce((acc, cur) => acc + cur, 0)
	);
	const ignoredAccordionTitle = $derived(
		`Ignored files [${totalIgnoredItems} items totalling ${convertToSizeString(totalIgnoredSize)}]`
	);

	function moveToIgnored(selectedItem: RenameEntry, onlySubtitle = false) {
		let msg = '';
		if (onlySubtitle) {
			item.ignored = [...item.ignored, selectedItem.subtitle!];
			delete selectedItem.subtitle;
			msg = 'Moved the selected subtitle to ignored list';
		} else {
			item.ignored = [
				...item.ignored,
				...[selectedItem.media, selectedItem.subtitle].filter((x) => !!x)
			];
			item.selected = item.selected.filter((i) => i !== selectedItem);
			msg = 'Moved the selected media and subtitle to ignored list';
		}

		if (msg.length > 0) {
			toastService.show(ToastFactory.createSuccessToast(msg));
		}
	}

	async function moveToSelected(ignoredItem: DirEntry) {
		const ft = getFiletype(ignoredItem.path, allowedExtensions);
		const movingSubtitle = ft === 'SUBTITLE';
		const movingMedia = ft === 'MEDIA';

		if (item.type === 'MOVIE') {
			let success = true;
			const mediaPresent = item.selected.length > 0;
			const subtitlePresent = mediaPresent && !!item.selected[0].subtitle;

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
		} else if (item.type === 'TV') {
			let success = true;

			const result = await popupService.showTVEpisodeEdit();
			if (result !== null) {
				const { season, episode } = result;
				const existingMedia = item.selected.find(
					(i) => i.season === season && i.episode === episode
				);
				const subtitlePresent = !!existingMedia?.subtitle;
				if (existingMedia && movingSubtitle) {
					if (subtitlePresent) {
						toastService.show(
							ToastFactory.createErrorToast(
								'Cannot move subtitle because there is another subtitle present for that media. Please remove that first before adding this one.'
							)
						);
						success = false;
					} else {
						existingMedia.subtitle = ignoredItem;
					}
				} else if (!existingMedia && movingMedia) {
					const newList = [...item.selected, { media: ignoredItem, season, episode }].sort(
						sortBySeasonAndEpisode
					);
					item.selected = newList;
				} else if (!existingMedia && movingSubtitle) {
					toastService.show(
						ToastFactory.createErrorToast(
							'Cannot find the media for this subtitle in the selected list. Please add media first.'
						)
					);
					success = false;
				}
			} else {
				success = false;
			}

			if (success) {
				const msg = `Moved the item to selected list.`;
				toastService.show(ToastFactory.createSuccessToast(msg));
				item.ignored = item.ignored.filter((x) => x !== ignoredItem);
			}
		} else {
			toastService.show(
				ToastFactory.createErrorToast('Operation not supported for type ' + item.type)
			);
		}
	}

	async function selectedItemEditHandler(entry: RenameEntry) {
		if (item.type === 'TV') {
			const { season: oldSeason, episode: oldEpisode } = entry;
			const result = await popupService.showTVEpisodeEdit(oldSeason, oldEpisode);

			if (result !== null) {
				entry.season = result.season;
				entry.episode = result.episode;
				const msg = `Changed item to Season ${result.season} Episode ${result.episode}`;
				toastService.show(ToastFactory.createSuccessToast(msg));
			} else {
				log.info('Edit operation cancelled by user');
			}
		} else {
			const msg = 'Edit not supported for movies';
			log.error(msg);
			toastService.show(ToastFactory.createErrorToast(msg));
		}
	}

	function sortBySeasonAndEpisode(a: RenameEntry, b: RenameEntry) {
		if (a.season! > b.season!) return 1;
		if (a.season! < b.season!) return -1;
		if (a.episode! > b.episode!) return 1;
		if (a.episode! < b.episode!) return -1;
		return 0;
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
				<MediaTag type={item.type} />
				<h3 class="text-md font-medium">{item.info.name} ({item.info.yearOfRelease})</h3>
				<p class="line-clamp-2 text-sm text-gray-500">
					{getBasename(item.entry.path)}
				</p>
			</section>
		</section>
	</section>
	<section class="media-target-select flex flex-col gap-2">
		<label for={`dd_${id}`}>Select destination:</label>
		<DestinationDirectoryDropdown
			bind:value={dest}
			destinationList={destinations}
			sourceSize={totalSelectedSize}
			id={`dd_${id}`}
			type={item.type}
		/>
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
									<Button
										type="mini-icon"
										title="Edit"
										onclick={() => selectedItemEditHandler(selectedItem)}
									>
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
