<script context="module">
    export async function preload(page) {
        const playlistId = page.query.playlistId
        const gameSessionId = parseInt(page.query.gameSessionId)
        const chatSocketAddress =  page.query.chatSocketAddress
        const gameSessionSocketAddress =  page.query.sessionSocketAddress
        return { gameSessionId,chatSocketAddress,gameSessionSocketAddress, playlistId };
    }
</script>



<script>
    import Game from '../components/Game.svelte';
    import { goto, stores } from '@sapper/app';

    const { session } = stores();

    export let playlistId = ""
    export let gameSessionId = 0
    export let chatSocketAddress = ""
    export let gameSessionSocketAddress = ""




    const startSessionRequestBody =  {
        service: "nq.GameSessionService", method: "GameSessionService.StartSession",
        request: { sessionId: gameSessionId}
    }

    const leaveSessionRequestBody =  {
        service: "nq.GameSessionService", method: "GameSessionService.LeaveSession",
        request: { sessionId: gameSessionId, userId: $session.user.id}
    }


    async function triggerStart() {

        const startRes = await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(startSessionRequestBody)
                }).then(r => r.json())
                .then(data => {
                    console.dir("start")
                    console.dir(data)
                    if (data.hasWorked === 1) {
                        alert(data.response)
                    } else {
                        console.dir(data)
                        goto('/home');
                        alert("Error:" + data.response)
                    }
                });


    }

    async function triggerLeave() {

        const leaveRes = await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(leaveSessionRequestBody)
                }).then(r => r.json())
                .then(data => {
                    console.dir(data)
                    if (data.hasWorked === 1) {
                        alert(data.response)
                        goto('/home');
                    } else {
                        console.dir(data)
                        alert("Error:" + data.response)
                        goto('/home');

                    }
                });


    }

    function start(){
        triggerStart()

    }

    function leaveGame() {
        triggerLeave()
        goto('/home')
    }

</script>

<div class="main_game">
    <button class="btn btn-dark" on:click|preventDefault|once ={start}>start Game</button>
    <button class="btn btn-outline-dark" on:click|preventDefault={leaveGame}>leave Game</button>
    <br />
    <br />
    <div class="card shadow-sm">
        <div class="game">
                    <Game sessionId={gameSessionId}
                          websocketChatAddress={chatSocketAddress}
                          websocketSessionAddress={gameSessionSocketAddress}
                          userId={$session.user.id}
                          listId={playlistId}
                    />
        </div>
    </div>
</div>
