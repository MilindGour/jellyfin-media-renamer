<script module>
	let _id = 0;

	function getNextId() {
		_id++;
		return joinStrings('srcdirlistitem', _id.toString());
	}
</script>

<script lang="ts">
	import type { MediaType } from '$lib/models/models';
	import { convertToSizeString, joinStrings } from '$lib/stores/util';
	import ToggleButtonGroup from '../toggleButtonGroup/toggleButtonGroup.svelte';

	const id = getNextId();
	const toggleOptions: MediaType[] = ['MOVIE', 'TV'];

	let {
		name,
		value = $bindable()
	}: {
		name: string;
		value: any;
	} = $props();

	let valid = $state(false);

	$effect(() => {
		if (value.selected) {
			valid = !!value.type;
		} else {
			valid = true;
		}
	});
</script>

<div
	class="src-dir-list-item flex gap-2 rounded p-3 {value.selected
		? valid
			? 'bg-green-50 hover:bg-green-100'
			: 'bg-red-50 hover:bg-red-100'
		: 'bg-gray-50 hover:bg-gray-100'}"
>
	<section class="checkbox">
		<input type="checkbox" {id} {name} bind:checked={value.selected} />
	</section>
	<section class="text grow-1">
		<label for={id} class="block cursor-pointer font-medium break-all">
			{value.entry.name}
			<p class="text-sm break-all text-gray-500">{convertToSizeString(value.entry.size)}</p>
		</label>
		<p class="mt-2">
			<ToggleButtonGroup
				id="toggle_{id}"
				bind:value={value.type}
				options={toggleOptions}
				required
			/>
		</p>
	</section>
</div>
