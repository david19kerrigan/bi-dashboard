<script lang="js">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";
	export let data = [{ ID: "", Name: "", Status: "", Address: "", Custom: "" }];
	onMount(async () => {
	  fetch("localhost:1080/getData")
	  .then(response => response.json())
	  .then(pgData => {
			console.log(data);
			data += pgData;
	  }).catch(error => {
		console.log(error);
		return [];
	  });
	});
</script>

<div class="counter">
	<table>
		<tr>
			{#each Object.keys(data[0]) as column}
				<th>{column}</th>
			{/each}
		</tr>
		{#each Object.values(data) as row}
			<tr>
				{#each Object.values(row) as column}
					<td>{column}</td>
				{/each}
			</tr>
		{/each}
		<tr>
			{#each Object.values(data[0]) as column}
				<td><input type="text" name="fname" /></td>
			{/each}
		</tr>
	</table>
</div>
