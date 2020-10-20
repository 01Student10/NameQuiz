<script context="module">
    const getAllContentListsRequestBody = {
        service: "nq.QuizContentService", method: "QuizContentService.GetAllContentLists", request: {}
    }
    const getAllGamesRequestBody = {
        service: "nq.GameService", method: "GameService.GetAllGames", request: {}
    }

    export async function preload() {


        const gameRes = await this.fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(getAllGamesRequestBody)
                })

        const allGames = await gameRes.json().then(data => data.games)

        const contentRes = await this.fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(getAllContentListsRequestBody)
                })

        const allContentLists = await contentRes.json().then(data =>

        {
            console.dir(data);
            if (!data.response) {
                //spaw keys and values from the Object..
            return Object.keys(data.nameToId).reduce((ret, key) => {
                    ret[data.nameToId[key]] = key;
                    return ret;
                }, {})
        } else {
            return {}
        }
        })



        return {allGames, allContentLists};
    }
</script>

<script>

    function preloadHomePage() {
       return [allGames, allContentLists]
    }
    import { goto, stores } from '@sapper/app';


    const { session, page } = stores();

    export let allGames
    export let allContentLists

    console.dir("page")
    console.dir(page)

    async function triggerJoin(sessionId) {
        const joinSessionRequestBody =  {
            service: "nq.GameSessionService", method: "GameSessionService.JoinSession",
            request: { sessionId: sessionId, userId: $session.user.id }
        }
        console.dir(joinSessionRequestBody)
        const joinRes = await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(joinSessionRequestBody)
                }).then(r => r.json())
                .then(data => {
                    console.dir("join")
                    console.dir(data)
                    if (data.hasWorked === 1) {
                        alert(data.response)
                    } else {
                        console.dir(data)
                        goto('/home');
                        alert("Error:" + data.response)
                        //navigateAndSave('/createGame',data)
                    }
                });


    }

    async function checkIfRunning(gameId,playlistId) {
        const isSessionRunningRequestBody =  {
            service: "nq.GameSessionService", method: "GameSessionService.IsSessionCreated",
            request: { gameId: gameId }
        }
        const runRes = await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(isSessionRunningRequestBody)
                }).then(r => r.json())
                .then(data => {
                    console.dir(data)
                    if (data.isRunning === 1) {
                        console.dir("check if running true")
                        console.dir(data)
                        join(data.sessionId)
                        goto('/gameSession?chatSocketAddress='+data.chatSocketAddress+
                                '&sessionSocketAddress='+data.sessionSocketAddress + '&playlistId='+playlistId
                         +'&gameSessionId='+data.sessionId)
                    } else {
                        console.dir("check if running false")
                        createAndJoin(gameId, playlistId)

                    }
                });


    }

    function deleteGame(gameId) {
        triggerDelete(gameId)
    }

    function join(sessionId) {
        triggerJoin(sessionId)
    }

    async function triggerDelete(gameId) {
        const deleteGameRequestBody =  {
            service: "nq.GameService", method: "GameService.DeleteGame", request: { id:gameId }
        }

        const _ = await fetch(`http://localhost:8085/rpc`,
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(deleteGameRequestBody)
            }).then(r => r.json())
            .then(data => {
                if(data.response) {
                    alert("Error: " + data.response)
                } else {
                    alert("Game deleted.")
                    goto('/reload?destination=home')
                }


            })
    }


   async function createAndJoin(gameId, playlistId){
        const creatSessionRequestBody =  {
            service: "nq.GameSessionService", method: "GameSessionService.CreateSession", request: { gameId:gameId }
        }



        const _ = await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(creatSessionRequestBody)
                }).then(r => r.json())
       .then(data => {
           console.dir(data.response)
           join(data.sessionId)
           goto('/gameSession?chatSocketAddress='+
                   data.chatSocketAddress+'&sessionSocketAddress='+
                   data.sessionSocketAddress
                   + '&playlistId='+playlistId
                   +'&gameSessionId='+data.sessionId)
       })

    }


</script>

<div>
    {#if allGames}
    <div class="main col-xl-12">
        <div class="card background shadow-sm">
            <div class="row">
            {#each allGames as game}
                    <div class="card games">
                        <div class="card-header game_header text-center">
                            {game.title}
                        </div>
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item">Playlist : {allContentLists[game.playlistId]}</li>
                            <li class="list-group-item">Rounds : {game.amountOfRounds}</li>
                            <li class="list-group-item text-center">
                                <button class="btn btn-dark" on:click|preventDefault={checkIfRunning(game.id,game.playlistId)}>
                                    join
                                </button>
                                <button class="btn btn-dark" on:click|preventDefault={deleteGame(game.id)}>
                                    delete
                                </button>
                            </li>
                        </ul>
                    </div>
        {/each}
            </div>
        </div>
    </div>
    {/if}

</div>