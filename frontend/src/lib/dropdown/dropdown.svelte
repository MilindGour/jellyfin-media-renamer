<script lang="ts">
	import Icon from 'svelte-awesome';
	import chevronDown from 'svelte-awesome/icons/chevronDown';
	import { slideDown } from '$lib/animations';
	import { DropdownModel } from './dropdown-model.svelte.ts';
	import { onDestroy, onMount } from 'svelte';
	import { DropdownService } from './dropdown-service.svelte.ts';

	const { id, options, labelProp = '', itemTemplate = null, initialValue = null } = $props();

	let dd: DropdownModel = new DropdownModel({
		options,
		labelProperty: labelProp,
		value: initialValue
	});

	$effect(() => {
		if (dd.open) {
			document.body.removeEventListener('click', handleOutsideClick);
			document.body.addEventListener('click', handleOutsideClick);
		} else {
			document.body.removeEventListener('click', handleOutsideClick);
		}
	});

	onMount(() => {
		console.log('[dbg] mounting dropdown with initialValue:', dd.value);
		DropdownService.register(id, dd);
	});
	onDestroy(() => {
		DropdownService.unregister(id);
	});

	export function getValue() {
		return dd.value;
	}

	function closeDropdown() {
		dd.open = false;
		document.body.removeEventListener('click', handleOutsideClick);
	}
	function handleDropdownHeadClick() {
		dd.open = !dd.open;
	}
	function handleOutsideClick(e: MouseEvent) {
		const ddBounds = dd.el!.getBoundingClientRect();
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
		dd.value = item;
	}
</script>

<div bind:this={dd.el} class="jmr-dropdown relative inline-block min-w-xs" class:open={dd.open}>
	<button
		type="button"
		class="dropdown-head inline-flex w-full cursor-pointer items-center gap-3 rounded-md border border-gray-400 p-2 text-start"
		onclick={handleDropdownHeadClick}
	>
		<span class="text grow-1">{dd.label || 'Please select'}</span>
		<span class="flex items-center justify-center" class:rotate-180={dd.open}
			><Icon data={chevronDown} /></span
		>
	</button>
	{#if dd.open}
		<div
			transition:slideDown
			class="dropdown-list absolute mt-px flex w-full flex-col overflow-hidden rounded-md border border-gray-400 bg-white"
		>
			{#each dd.options as option}
				<button
					onclick={() => handleItemSelect(option)}
					class="item block cursor-pointer p-2 text-left"
					class:bg-gray-200={JSON.stringify(option) === JSON.stringify(dd.value)}
					class:hover:bg-gray-50={JSON.stringify(option) !== JSON.stringify(dd.value)}
				>
					{#if itemTemplate}
						{@render itemTemplate(option)}
					{:else}
						{dd.label}
					{/if}
				</button>
			{/each}
		</div>
	{/if}
</div>
