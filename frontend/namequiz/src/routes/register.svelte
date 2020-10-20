

<style>
    h1, figure, p, div {
        text-align: center;
        margin: 0 auto;
    }
    h2 {
        text-align: center;

    }
    label {
        margin-top: 50px;
        margin-bottom: 10px;
    }

    h1 {
        font-size: 2.8em;
        text-transform: uppercase;
        font-weight: 700;
        margin: 0 0 0.5em 0;
    }

    figure {
        margin: 0 0 1em 0;
    }

    img {
        width: 100%;
        max-width: 400px;
        margin: 0 0 1em 0;
    }

    p {
        margin: 1em auto;
    }

    @media (min-width: 480px) {
        h1 {
            font-size: 4em;
        }
    }
</style>

<svelte:head>
    <title>Register</title>
</svelte:head>


<script>

    import { goto, stores } from '@sapper/app';

    const { session } = stores();

    let uName = 'Bill'
    let passwd = '42'
    let email = 'example@test.com'


    async function getData() {
        let bodyData = {
            service: "nq.UserService",
            method: "UserService.CreateUser",
            request: { username: uName, password : passwd, mail : email}
        }
        const res = await fetch(`http://localhost:8085/rpc`,
                {
                    method: 'POST',
                    headers: {

                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(bodyData)
                }).then(r => r.json())
                .then(data => {
                    if(data.hasCreatedUser === 1){
                        $session.user = data.user;
                        goto('/home');

                    } else {
                        alert(data.Response)

                    }
                });

    }

    function createUser() {
        //console.dir("about to get data:")
        getData()
    }

</script>

<!--<h1>Register</h1>-->
<!--<br/>-->
<!--<br/>-->
<!--<div>-->
<!--    <label>-->
<!--        <h2>Username:</h2>-->
<!--        <input  bind:value={uName}>-->
<!--    </label>-->
<!--    <br/>-->
<!--    <br/>-->
<!--    <label>-->
<!--        <h2>Email:</h2>-->
<!--        <input  bind:value={email}>-->
<!--    </label>-->
<!--    <br/>-->
<!--    <br/>-->
<!--    <label>-->
<!--        <h2>Password:</h2>-->
<!--        <input  bind:value={passwd}>-->
<!--    </label>-->
<!--    <br/>-->
<!--    <br/>-->
<!--    <button on:click|preventDefault={createUser}>-->
<!--        register-->
<!--    </button>-->
<!--    <a href="."> already have an account? </a>-->
<!--</div>-->
<div class="main_reg">
    <div class="card shadow-sm">
        <div class="card-body">
            <h1>Register</h1>
            <div class="form-label-group">
                <label class="inputfields">
                    <small class="form-text text-muted">Username:</small>
                    <input type="text" class="form-control" bind:value={uName}>
                </label>
            </div>
            <div class="form-label-group">
                <label class="inputfields">
                    <small class="form-text text-muted">Email:</small>
                    <input type="email" class="form-control" bind:value={email}>
                </label>
            </div>
            <div class="form-label-group">
                <label class="inputfields">
                    <small class="form-text text-muted">Password:</small>
                    <input type="password" class="form-control" bind:value={passwd}>
                </label>
            </div>
            <br />
            <div class="col-md-12">
                <button class="btn btn-dark" on:click|preventDefault={createUser}>
                    Register
                </button>
                <a href="." id="register_link"> already have an account? </a>
            </div>
        </div>
    </div>
</div>

