<script lang="ts">
	import Icon from 'svelte-awesome';
	import chevronDown from 'svelte-awesome/icons/chevronDown';
	import { slideDown } from '$lib/animations';

	let { id, options, labelProp = '', itemTemplate = null, value = $bindable(null) } = $props();

	let open = $state(false);
	let el: HTMLElement;
	let label = $state('');

	$effect(() => {
		if (open) {
			document.body.removeEventListener('click', handleOutsideClick);
			document.body.addEventListener('click', handleOutsideClick);
		} else {
			document.body.removeEventListener('click', handleOutsideClick);
		}
	});

	$effect(() => {
		if (value === null) {
			label = 'Please Select';
		} else {
			label = labelProp?.length > 0 ? value[labelProp] : value;
		}
	});

	function closeDropdown() {
		open = false;
		document.body.removeEventListener('click', handleOutsideClick);
	}
	function handleDropdownHeadClick() {
		open = !open;
	}
	function handleOutsideClick(e: MouseEvent) {
		const ddBounds = el?.getBoundingClientRect();
		if (!ddBounds) {
			return;
		}
		if (
			e.clientX < ddBounds.left ||
			e.clientX > ddBounds.right ||
			e.clientY < ddBounds.top ||
			e.clientY > ddBounds.bottom
		) {
			closeDropdown();
		}
	}

	function handleItemSelect(item: any) {
		value = item;
	}
</script>

<div
	bind:this={el}
	data-id={id}
	class="jmr-dropdown relative inline-block min-w-full sm:min-w-sm"
	class:open
>
	<button
		type="button"
		class="dropdown-head inline-flex w-full cursor-pointer items-center gap-3 rounded-md border border-gray-400 p-2 text-start"
		onclick={handleDropdownHeadClick}
	>
		<span class="text grow-1">{label}</span>
		<span class="flex items-center justify-center" class:rotate-180={open}
			><Icon data={chevronDown} /></span
		>
	</button>
	{#if open}
		<div
			transition:slideDown
			class="dropdown-list absolute z-10 mt-px flex w-full flex-col overflow-hidden rounded-md border border-gray-400 bg-white"
		>
			{#each options as option (option)}
				<button
					onclick={() => handleItemSelect(option)}
					class="item block cursor-pointer p-2 text-left"
					class:bg-gray-200={JSON.stringify(option) === JSON.stringify(value)}
					class:hover:bg-gray-50={JSON.stringify(option) !== JSON.stringify(value)}
				>
					{#if itemTemplate}
						{@render itemTemplate(option)}
					{:else}
						{label}
					{/if}
				</button>
			{/each}
		</div>
	{/if}
</div>
