<script context="module">
    export async function preload(page) {
        const dest = page.query.destination
        return { dest};
    }
</script>


<script>
    import { goto } from '@sapper/app';

    export let dest = ''
    goto('/'+dest)
</script>

<h1>reloading...</h1>

