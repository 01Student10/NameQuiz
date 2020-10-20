<script context="module">
    const getAllContentListsRequestBody = {
        service: "nq.QuizContentService", method: "QuizContentService.GetAllContentLists", request: {}
    }
    export async function preload() {
        const res1 = await this.fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(getAllContentListsRequestBody)
                })

        const res2 = await this.fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(getAllContentListsRequestBody)
                })

        const namesToIds = await res1.json(data => {
            console.dir(data.nameToId)
            data.nameToId
        }
            )
        const names = await res2.json().then(data => {
            var result = []
            for (let entry in data.nameToId) {
                result = result.concat(entry)
            }
            return result

        })
        return {namesToIds, names};
    }
</script>


<script>
    import {goto, stores} from '@sapper/app';

    const {session} = stores();

    export let namesToIds = {}
    export let names = []
    let selectedPlaylistName = names[0]
    let playlistName = ""
    let amountOfRounds = 1

    async function getData() {
        let bodyData = {
            service: "nq.GameService",
            method: "GameService.CreateGame",
            request: {
                title: playlistName,
                amountOfSongs : amountOfRounds,
                playlistId: namesToIds.nameToId[selectedPlaylistName],
                owner: $session.user.id
        }}

         await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(bodyData)
                }).then(r => r.json())
                .then(data => {
                    if (data.response !== "") {
                        alert("successfully created game!")
                        goto('/home');
                    } else {
                        alert(data.response)
                    }
                });

     }

    function createGame() {
        getData()
    }
</script>

<div class="main_create_game">
    <div class="card shadow-sm">
        <h1>create a game here</h1>
        <div class="card-body">
            <small class="form-text text-muted">Choose a playlist:</small>
            {#if names.length !== 0 }
                <label for="gameselect"></label><select class="form-control" id="gameselect" bind:value={selectedPlaylistName}>
                        { #each names as name}
                            <option value={name}>{name}</option>
                        {/each}
                    </select>
            {/if}

            <br/>
            <br/>
            <div class="form-label-group">
                <label class="inputfields">
                    <small class="form-text text-muted">Playlist Title:</small>
                    <input type="text" class="form-control" bind:value={playlistName}>
                </label>
            </div>
            <br/>
            <br/>
            <div class="form-label-group">
                <label class="inputfields">
                    <small class="form-text text-muted"># Rounds to play:</small>
                    <input type="text" class="form-control" bind:value={amountOfRounds}>
                </label>
            </div>
            <br/>
            <br/>
            <div class="col-md-12">
                <button class="btn btn-dark" on:click|preventDefault={createGame}>
                    create Game
                </button>
            </div>
        </div>
    </div>
</div>

