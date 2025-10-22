<script lang="ts">
	import type { MediaInfo } from '$lib/models';
	import type { ChangeEventHandler } from 'svelte/elements';

	const {
		item,
		name,
		onchange
	}: { item: MediaInfo; name: string; onchange: ChangeEventHandler<HTMLInputElement> } = $props();
</script>

<label class="media-searchlist-item flex items-center gap-2 rounded bg-gray-100 p-2">
	<input class="" type="radio" value={item.mediaId} {name} {onchange} />
	<section class="info flex gap-2">
		<img
			class="aspect-[2/3] max-h-16 rounded"
			src={item.thumbnailUrl || 'noimg.svg'}
			alt="poster"
		/>
		<section class="text">
			<h3 class="text-md font-medium">{item.name} ({item.yearOfRelease})</h3>
			<p class="line-clamp-2 text-sm text-gray-500">{item.description || 'No Description'}</p>
		</section>
	</section>
</label>
