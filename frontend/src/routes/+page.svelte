<script lang="ts">
	import { Dropdown, DropdownService } from '$lib/components/dropdown';
	import { Button } from '$lib/components/button';

	function formatPath(path: string): string {
		if (path.startsWith('/')) {
			path = path.substring(1);
		}
		return path.split('/').join(' > ');
	}

	function handleScanDirClick() {
		// This method will get the selected item from the dropdown
		const value = DropdownService.getValueOf('dropdown1');
		console.log('Dropdown1 value =', $state.snapshot(value));
	}
</script>

<section class="page flex flex-col gap-8">
	<section
		class="form-section flex flex-col flex-wrap items-stretch gap-2 sm:flex-row sm:items-start"
	>
		<label class="basis-full" for="dropdown1">Please select media source directory</label>
		<Dropdown id="dropdown1" options={[]} itemTemplate={dropdownTemplate} />
		<Button type="primary" onclick={handleScanDirClick}>Scan Directory</Button>
	</section>
	<section class="list-section">This is list section.</section>
</section>

{#snippet dropdownTemplate(item: any)}
	<div class="item-instance">
		<div class="font-semibold">{item.label}</div>
		<div class="text-sm text-gray-500">{formatPath(item.path)}</div>
	</div>
{/snippet}
