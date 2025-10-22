<script lang="ts">
	import type { MediaType } from '$lib/models';

	const { type }: { type: MediaType } = $props();
</script>

<div
	class="media-tag w-fit rounded px-1 py-0.5 text-xs uppercase {type === 'MOVIE'
		? 'bg-green-200'
		: 'bg-purple-200'}"
>
	{type}
</div>
