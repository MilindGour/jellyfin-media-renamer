<script module>
	let _id = 0;

	function getNextId() {
		return ++_id;
	}
</script>

<script lang="ts">
	import type { SourceDirWithInfo } from '$lib/models/models';

	import { formatPathString, joinStrings } from '$lib/stores/util';
	import MediaSearchList from '../mediasearchlist/media-search-list.svelte';

	import MediaTag from '../mediaTag/media-tag.svelte';

	let { item = $bindable() }: { item: SourceDirWithInfo } = $props();
	const id = getNextId();
	const identifiedNameId = joinStrings('srcdirinfoitem', id.toString());
</script>

<section
	class="source-directory-info-item flex flex-col gap-px rounded border border-l-4 border-gray-200 bg-gray-50 px-4 py-3 {item?.identifiedMediaId
		? 'border-l-green-400'
		: 'border-l-red-400'}"
>
	<h2 class="block font-medium break-all">{item.sourceDirectory.entry.name}</h2>
	<p class="text-sm break-all text-gray-500">{formatPathString(item.sourceDirectory.entry.path)}</p>
	<div class="tag-wrapper mt-2">
		<MediaTag type={item.sourceDirectory.type!} />
	</div>
	<div class="input-wrapper mt-2 flex flex-col items-start gap-px">
		<label class="text-sm text-gray-500" for={identifiedNameId}>Identified year and name</label>
		<div
			class="name-and-year flex max-w-lg self-stretch overflow-hidden rounded border border-gray-200"
		>
			<input
				type="number"
				class="w-16 border-0 p-1 outline-0"
				bind:value={item.identifiedMediaYear}
			/>
			<input
				type="text"
				id={identifiedNameId}
				class="grow-1 border-0 border-l-2 border-gray-200 p-1 outline-0"
				bind:value={item.identifiedMediaName}
			/>
		</div>
	</div>
	<div class="media-search-list-container mt-3">
		<MediaSearchList list={item.identifiedMediaInfos || []} bind:value={item.identifiedMediaId!} />
	</div>
</section>
