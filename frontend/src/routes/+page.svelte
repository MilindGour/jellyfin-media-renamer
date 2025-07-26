<script lang="ts">
	import { Dropdown, DropdownService } from '$lib/dropdown';

	const items = [
		{ label: 'Media TV', path: '/mnt/media/tv' },
		{ label: 'Media Movies', path: '/mnt/media/movies' },
		{ label: 'Transmission Downloads', path: '/var/lib/transmission-daemon/downloads' },
		{ label: 'Data TV', path: '/mnt/data/tv' },
		{ label: 'Data Movie', path: '/mnt/data/movie' }
	];

	function formatPath(path: string): string {
		if (path.startsWith('/')) {
			path = path.substring(1);
		}
		return path.split('/').join(' > ');
	}

	function getSelectedItemFromDropdown() {
		// This method will get the selected item from the dropdown
		const value = DropdownService.getValueOf('dropdown1');
		console.log('Dropdown1 value =', $state.snapshot(value));
	}
</script>

<Dropdown id="dropdown1" options={items} itemTemplate={dropdownTemplate} initialValue={items[2]} />

<button onclick={getSelectedItemFromDropdown} class="cursor-pointer">Get Selected Item</button>

{#snippet dropdownTemplate(item: any)}
	<div class="item-instance">
		<div class="font-semibold">{item.label}</div>
		<div class="text-sm text-gray-500">{formatPath(item.path)}</div>
	</div>
{/snippet}
