<script>
	import Nav from '../components/Nav.svelte';

	export let segment;
</script>

{#if segment !== undefined && segment !== "register"}
	<Nav {segment}/>
{/if}

<main>
	<slot></slot>
</main>
