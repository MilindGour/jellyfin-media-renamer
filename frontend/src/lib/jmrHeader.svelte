<script lang="ts">
	function logoClickHandler(e: Event) {
		e.preventDefault();
		window.location.assign('/');
	}
</script>

<header class="flex justify-center py-8 sm:justify-start">
	<a href="/" onclick={logoClickHandler}>
		<img src="/jmr.svg" alt="logo" />
	</a>
</header>
