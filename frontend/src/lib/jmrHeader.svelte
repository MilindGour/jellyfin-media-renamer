<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from './components';
	import { PopupService } from './components/popup';

	const ps = new PopupService();

	function logoClickHandler(e: Event) {
		e.preventDefault();
		window.location.assign('/');
	}
</script>

<header class="flex items-center justify-center py-8 sm:justify-start">
	<a href="/" onclick={logoClickHandler}>
		<img src="/jmr.svg" alt="logo" />
	</a>
	<div class="right-container ml-auto">
		<Button onclick={() => goto('/new-media')}>New Media</Button>
		<Button onclick={() => ps.showFileTransferStatusPopup()}>FTP</Button>
	</div>
</header>
