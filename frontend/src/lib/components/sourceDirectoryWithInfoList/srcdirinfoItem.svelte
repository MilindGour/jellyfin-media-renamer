<script module>
	let _id = 0;

	function getNextId() {
		return ++_id;
	}
</script>

<script lang="ts">
	import { formatPathString, joinStrings } from '$lib/stores/util';
	import MediaSearchList from '../mediasearchlist/media-search-list.svelte';

	import MediaTag from '../mediaTag/media-tag.svelte';

	const id = getNextId();
	const identifiedNameId = joinStrings('srcdirinfoitem', id.toString());

	let mediaId: string;
</script>

<section class="source-directory-info-item flex flex-col gap-px rounded bg-gray-50 px-4 py-3">
	<h2 class="block font-medium break-all">Media dirty name from directory</h2>
	<p class="text-sm text-gray-500">{formatPathString('/Users/milindgour/Desktop')}</p>
	<div class="tag-wrapper mt-2">
		<MediaTag type="Movie" />
	</div>
	<div class="input-wrapper mt-2 flex flex-col items-start gap-px">
		<label class="text-sm text-gray-500" for={identifiedNameId}>Identified Name</label>
		<input type="text" id={identifiedNameId} class="min-w-md rounded border-gray-200 p-1" />
	</div>
	<div class="media-search-list-container mt-3">
		<MediaSearchList bind:value={mediaId} />
	</div>
</section>
