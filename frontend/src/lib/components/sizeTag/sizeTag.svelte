<script lang="ts">
	import { convertToSizeString } from '$lib/stores/util';

	const { bytes } = $props();
</script>

<div class="size-tag inline-block rounded bg-gray-400 p-1 text-xs text-white">
	{convertToSizeString(bytes)}
</div>
