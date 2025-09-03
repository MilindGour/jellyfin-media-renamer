<script module>
	let _id = 0;

	function getNextId() {
		_id++;
		return joinStrings('srcdirlistitem', _id.toString());
	}
</script>

<script lang="ts">
	import type { DirEntry } from '$lib/models/models';
	import { convertToSizeString, joinStrings } from '$lib/stores/util';
	import ToggleButtonGroup from '../toggleButtonGroup/toggleButtonGroup.svelte';

	const id = getNextId();
	const toggleOptions = ['Movie', 'Tv'];

	let {
		entry,
		name,
		value = $bindable(null),
		valid = $bindable(true)
	}: { entry: DirEntry; name: string; value: any; valid: boolean } = $props();

	let toggleValid = $state(false);
	let toggleValue = $state(null);
	let isSelected = $state(false);

	$effect(() => {
		if (isSelected) {
			valid = toggleValid;
			value = valid ? { entry, type: toggleValue } : null;
		} else {
			valid = true;
			value = null;
		}
	});
</script>

<div
	class="src-dir-list-item flex gap-2 rounded p-3 {isSelected
		? valid
			? 'bg-green-50 hover:bg-green-100'
			: 'bg-red-50 hover:bg-red-100'
		: 'bg-gray-50 hover:bg-gray-100'}"
>
	<section class="checkbox">
		<input type="checkbox" {id} {name} bind:checked={isSelected} />
	</section>
	<section class="text grow-1">
		<label for={id} class="block cursor-pointer font-medium break-all">
			{entry.name}
			<p class="text-sm text-gray-500">{convertToSizeString(entry.size)}</p>
		</label>
		<p class="mt-2">
			<ToggleButtonGroup
				id="toggle_{id}"
				bind:valid={toggleValid}
				bind:value={toggleValue}
				options={toggleOptions}
				required
			/>
		</p>
	</section>
</div>
