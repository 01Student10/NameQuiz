<script>

    import { stores } from '@sapper/app';
    const { session } = stores();
    import Message from "./Message.svelte";

    let message;
    let messages = [];
    export let isAllowedToChat;
    export let socketAddress = ""
    export let sessionId = 0
    export let listId = ""
    export let entryId = ""
    const socket = new WebSocket(socketAddress);

    // Connection opened
    socket.addEventListener('open', function (event) {
        console.log("It's open");
    });

    // Listen for messages
    socket.addEventListener('message', function (event) {
        console.log("got event message")
        let data = JSON.parse(event.data)
        console.log(data)
        let messageObj = {guess: data.guess, wasRight: data.wasRight === 1 ? true : false, userName: data.userName}
        messages = messages.concat(messageObj);
        updateScroll();
    });

    async function fetchMessage() {

        const chatMessageBody = {
            userId: $session.user.id,
            userName: "",
            guess: message,
            sessionId: sessionId,
            listId : listId,
            entryId : entryId,
            wasRight : 2
        }
        console.dir("chatMessageBody")
        console.dir(chatMessageBody)

        const sendMessage = await fetch(`http://localhost:8086/Chat/UserMessage/`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(chatMessageBody)
                });


    }


    function onSendMessage() {
        if (message.length > 0) {
            fetchMessage()
            message = "";
        }
    }

    function updateScroll() {
        const chatWindow = document.getElementById('chatWindow');
        setTimeout(() => {
            chatWindow.scrollTop = chatWindow.scrollHeight;
        }, 0);
    }
</script>

    <div class="card">
        <div class="chatwindow overflow-auto" id="chatWindow">
            <ul id="messages">
                {#if messages}
                {#each messages as message}
                    <Message guess={message.guess} userName={message.userName} wasRight={message.wasRight}/>
                {/each}
                {/if}
            </ul>
        </div>
        <div class="chatwindow_input">
            <form action="">
                {#if isAllowedToChat}
                    <div class='input-group mb-3'>
                        <div class='input-group-prepend'>
                            <span class='input-group-text' id='basic-addon3'>Guess</span>
                        </div>
                        <input type='text' class='form-control' autocomplete="off" bind:value={message} aria-describedby='basic-addon3'>
                        <div class='input-group-append'>
                            <button class='btn btn-outline-dark' on:click|preventDefault={onSendMessage}>Send</button>
                        </div>
                    </div>
                {/if}
            </form>
        </div>
    </div>
