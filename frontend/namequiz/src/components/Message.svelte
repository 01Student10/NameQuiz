<script>
    export let guess;
    export let wasRight;
    export let userName;

</script>

{#if wasRight}
    <p id="right_guess"><b>{userName}</b> : <i>{guess}</i></p>
{:else}
    <p id="wrong_guess"><b>{userName}</b> : <i>{guess}</i></p>
{/if}