<script>
    import Display from "./Display.svelte";
    import Chat from "./Chat.svelte";
    import ScoreEntry from "./ScoreEntry.svelte";
    import { goto } from '@sapper/app';

    export let websocketSessionAddress = ""
    export let websocketChatAddress = ""
    export let sessionId = 0
    export let userId = 0
    export let listId = ""
    let pathToPic = "startPic.jpg"
    let picInfo = ""
    let inRound = "about to start"
    let currentGameTitle = "default Name"
    let entryId = ""
    let userNameToScore = []
    let duration = 30000
    let elapsed = duration;
    let last_time = window.performance.now();
    let frame;
    let isAllowedToChat = false;
    let currentStatus = inRound;

    (function update() {
        frame = requestAnimationFrame(update);

        const time = window.performance.now();
        elapsed += Math.min(
                time - last_time,
                duration - elapsed
        );

        last_time = time;
    }());


    let socket = new WebSocket(websocketSessionAddress);

    // Connection opened
    socket.addEventListener('open', function (event) {
        console.log("It's open");
    });

    // Listen for messages
    socket.addEventListener('message', function (event) {
        isAllowedToChat = true
        console.log("got event message")
        let data = JSON.parse(event.data)

        console.log(data)
        if(data.Status === "Error") {
            alert("Error occured : " + data.ErrorMsg + "\n Try another Game" )
            goto('/home')
        } else {
            if (data.Status !== "Game Over") {
                inRound = data.Status
                pathToPic = data.Path
                picInfo = data.PicInfo
                currentGameTitle = data.GameTitle
                if (Object.keys(data.UserNameToScore).length !== 0) {
                    let res = []
                    for (const [key, value] of Object.entries(data.UserNameToScore)) {
                        console.log(`${key}: ${value}`);
                        let entry = {username: key, score: value}
                        res.push(entry)
                    }
                    userNameToScore = res
                }
                entryId = data.EntryID
                if (inRound !== currentStatus){
                    elapsed = 0
                    currentStatus = inRound
                }

            } else {
                console.dir(data.Winner)
                let resultString = ""
                data.Winner.forEach(element => resultString = resultString.concat(`${element.Name}: ${element.Score} \n`))
                alert("Winner(s): \n" + resultString)
                goto('/home')
            }
        }

    });

</script>
<div class="top row">
    <div class="top_left"></div>
    <div class="top_right">
        <div class="top_right_top row">
            <div class="top_right_top_left"></div>
            <div class="top_right_top_right row">
                <div class="title">
                    <p><b>Title:</b><span class="badge_h badge badge-secondary">{currentGameTitle}</span></p>
                </div>
                <div class="phase">
                    <p><b>Phase:</b><span class="badge_h badge badge-secondary">{inRound}</span></p>
                </div>
            </div>
        </div>
        <div class="top_right_bot">
            <span class="progress_bar">
                Time per round:&nbsp;&nbsp;<progress value="{elapsed / duration}"></progress>
            </span>
        </div>
    </div>
</div>
<div class="bottom row">
    <div class="bottom_left">
        <Display src="images/{pathToPic}" picInfo={picInfo}/>
    </div>
    <div class="bottom_middle">
        <Chat socketAddress={websocketChatAddress} sessionId={sessionId} entryId={entryId} listId={listId} isAllowedToChat={isAllowedToChat}/>
    </div>
    <div class="bottom_right">
        <table class="table table-hover">
            <thead>
                <tr>
                    <th scope="col">Username</th>
                    <th scope="col">Score</th>
                </tr>
            </thead>
            <tbody>
            {#if userNameToScore.length !== 0 }
                {#each userNameToScore as entry}
                    <tr><ScoreEntry username={entry.username} score={entry.score}/></tr>
                {/each}
                {/if}
            </tbody>
        </table>
    </div>
</div>




