


<svelte:head>
	<title>Login</title>
</svelte:head>


<script>

	import { goto, stores } from '@sapper/app';

	const { session } = stores();

 let uName = ''
 let passwd = ''

 async function getData() {
	 let bodyData = {
		 service: "nq.UserService",
		 method: "UserService.Login",
		 request: { userName: uName, password : passwd}
	 }
	 const _ = await fetch(`http://localhost:8085/rpc`,
			 {
				 method: 'POST',
				 headers: {
					 'Content-Type': 'application/json',
				 },
				 body: JSON.stringify(bodyData)
			 }).then(r => r.json())
			 .then(data => {
				 if(data.hasLoggedIn === 1){
					 $session.user = data.user;
					 goto('/home');

				 } else {
					 alert("login failed, please try again")
				 }
			 });

 }

 function hasUser() {
	getData()
 }

</script>
<div class="main_log col-xl-12">
	<div class="card shadow-sm">
		<div class="card-body">
			<h1>Login</h1>
			<div class="form-label-group">
				<label class="inputfields">
					<small class="form-text text-muted">Username:</small>
					<input type="text" class="form-control" bind:value={uName}>
				</label>
			</div>
			<br />
			<div class="form-label-group">
				<label class="inputfields">
					<small class="form-text text-muted">Password:</small>
					<input type="password" class="form-control" bind:value={passwd}>
				</label>
			</div>
			<br />
			<div class="col-md-12">
					<button class="btn btn-dark" on:click|preventDefault={hasUser}>
						login
					</button>
					<a href="register" id="register_link"> Register </a>
			</div>
		</div>
	</div>
</div>

