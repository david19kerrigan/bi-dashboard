<script lang="ts">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";
	let data = [
		{ ID: "", Name: "", Status: "", Address: "", Custom: "" },
	];
	console.log("here2");
	onMount(async () => {
		fetch("http://localhost:8080/getData")
			.then((response) => response.json())
			.then((pgData) => {
				console.log(pgData);
				data = pgData["data"];
			})
			.catch((error) => {
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
			<td><input type="text" name="id" /></td>
			<td><input type="text" name="name" /></td>
			<td><input type="text" name="status" /></td>
			<td><input type="text" name="address" /></td>
		</tr>
	</table>
</div>

<div>
	<button type="submit">
		Submit Changes
	</button>
</div>
