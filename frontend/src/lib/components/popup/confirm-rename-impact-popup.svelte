<script lang="ts">
	import { PopupComponent } from '.';
	import { Button } from '$lib/components';
	import type { ConfirmMediaRequestItem, DestConfig } from '$lib/models';
	import { Log } from '$lib/services/logger';
	import { convertToSizeString } from '$lib/stores/util';

	const { data, onclose } = $props();
	const log = new Log('Confirm Rename Impact Popup');

	const selectionImpacts = $derived(
		computeRenameImpact(data?.selections || [], data?.allDestinations || [])
	);

	function computeRenameImpact(
		selections: ConfirmMediaRequestItem[],
		allDestinations: DestConfig[]
	) {
		const impactObject = selections
			.map((sel) => getSingleItemImpactOnDisk(sel, allDestinations))
			.reduce(impactListReducer, {});

		return impactObjectToImpactInfo(impactObject, allDestinations);
	}

	function impactObjectToImpactInfo(
		reducedImpactList: { [k: string]: number },
		allDestinations: DestConfig[]
	) {
		let output = [];
		for (const [mountPoint, impactBytes] of Object.entries(reducedImpactList)) {
			const mountPointObject = allDestinations.find((d) => d.mount_point === mountPoint);
			output.push({
				mount_point: mountPoint,
				total_size_bytes: mountPointObject!.total_size_kb * 1000,
				used_size_bytes_before: mountPointObject!.used_size_kb * 1000,
				free_size_bytes_before: mountPointObject!.free_size_kb * 1000,
				used_size_bytes_after: mountPointObject!.used_size_kb * 1000 + impactBytes,
				free_size_bytes_after: mountPointObject!.free_size_kb * 1000 - impactBytes,
				difference: -1 * impactBytes
			});
		}

		return output;
	}

	function impactListReducer(
		obj: { [k: string]: number },
		impact: { [k: string]: number }
	): { [k: string]: number } {
		for (const [k, v] of Object.entries(impact)) {
			if (k in obj) {
				obj[k] += v;
			} else {
				obj[k] = v;
			}
		}
		return obj;
	}

	function getSingleItemImpactOnDisk(
		input: ConfirmMediaRequestItem,
		allDestinations: DestConfig[]
	) {
		const selectedTotalSize = input.selected.reduce(
			(sum, sel) => sum + (sel.media.size + (sel?.subtitle?.size || 0)),
			0
		);

		const ignoredTotalSize = input.ignored.reduce((sum, item) => sum + item.size, 0);

		const fromMountPointStr =
			allDestinations.find((dest) => input.selected[0].media.path.startsWith(dest.mount_point))
				?.mount_point || null;
		const toMountPointStr = input.destination.mount_point;

		if (!fromMountPointStr) {
			console.error('Unknown mount_point computed from the given selections. Cannot continue!');
			throw Error('Invalid mount_point computed');
		}

		const output: { [k: string]: number } = {};
		output[fromMountPointStr] = -(selectedTotalSize + ignoredTotalSize);

		if (fromMountPointStr === toMountPointStr) {
			output[toMountPointStr] += selectedTotalSize;
		} else {
			output[toMountPointStr] = selectedTotalSize;
		}

		log.info('impacts:', output);
		return output;
	}
</script>

<PopupComponent title={data.title || 'Confirm impact of size on mount points'}>
	{#snippet body()}
		{#each selectionImpacts as impactInfo}
			<table class="not-first:mt-6">
				<tbody>
					<tr><td>Mount Point:</td><td class="pl-16 font-semibold">{impactInfo.mount_point}</td></tr
					>
					<tr
						><td>Total Size:</td><td class="pl-16 font-semibold"
							>{convertToSizeString(impactInfo.total_size_bytes)}</td
						></tr
					>
					<tr
						><td>Used Size (Before):</td><td class="pl-16 font-semibold"
							>{convertToSizeString(impactInfo.used_size_bytes_before)}</td
						></tr
					>
					<tr
						><td>Used Size (After):</td><td class="pl-16 font-semibold"
							>{convertToSizeString(impactInfo.used_size_bytes_after)}</td
						></tr
					>
					<tr
						><td>Impact on Size:</td><td
							class="pl-16 font-semibold {impactInfo.difference >= 0
								? 'text-green-600'
								: 'text-red-600'}">{convertToSizeString(impactInfo.difference)}</td
						></tr
					>
				</tbody>
			</table>
		{/each}
	{/snippet}
	{#snippet footer()}
		<Button onclick={() => onclose(true)}>Confirm</Button>
		<Button onclick={() => onclose(false)}>Cancel</Button>
	{/snippet}
</PopupComponent>
