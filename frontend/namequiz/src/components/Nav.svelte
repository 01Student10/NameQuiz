<script>
	import { goto, stores } from '@sapper/app';


	const { session } = stores();

	export let segment;

	function logout() {
		goto('/')
	}


</script>

<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
	<div class="collapse navbar-collapse" id="navbarSupportedContent">
		<ul class="navbar-nav mr-auto">
			<li class="nav-item active">
				<a class="nav-link" aria-current="{segment === 'home' ? 'page' : undefined}" href="home">Home <span class="sr-only">(current)</span></a>
			</li>
			<li class="nav-item">
				<a class="nav-link" aria-current="{segment === 'createGame' ? 'page' : undefined}" href="createGame">Create Game</a>
			</li>
		</ul>
		<form class="form-inline my-2 my-lg-0">
			{#if {$session}}
				<span class="nav_greeting">Welcome,</span> &nbsp;&nbsp;
				<span class="nav_greeting" id="nav_username">
				{$session.user.username}</span>
			{/if}
			<button class="btn btn-light logoutButton" on:click|preventDefault={logout}>logout</button>
		</form>
	</div>
</nav>
