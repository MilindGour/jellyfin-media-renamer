<script lang="ts">
	import type { MediaInfo } from '$lib/models/models';
	import type { ChangeEventHandler } from 'svelte/elements';

	const {
		item,
		name,
		onchange
	}: { item: MediaInfo; name: string; onchange: ChangeEventHandler<HTMLInputElement> } = $props();
</script>

<label class="media-searchlist-item flex items-center gap-3 rounded bg-gray-100 p-2">
	<input type="radio" value={item.mediaId} {name} {onchange} />
	<section class="info flex gap-2">
		<img class="max-h-16 rounded" src={item.thumbnailUrl} alt="{item.name} poster" />
		<section class="text">
			<h3 class="text-md font-medium">{item.name} ({item.yearOfRelease})</h3>
			<p class="text-sm text-gray-500">{item.description}</p>
		</section>
	</section>
</label>
