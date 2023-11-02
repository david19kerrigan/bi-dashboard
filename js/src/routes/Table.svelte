<script lang="ts">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";
	let data = [{ ID: "", Name: "", Status: "", Address: "", Custom: "" }];
	onMount(async () => {
		fetch("http://localhost:8080/getData")
			.then((response) => response.json())
			.then((pgData) => {
				data = pgData["data"];
			})
			.catch((error) => {
				console.log(error);
				return [];
			});
	});
	const handleSubmit = (e) => {
		// getting the action url
		const ACTION_URL = "http://localhost:8080/postData";

		// get the form fields data and convert it to URLSearchParams
		const formData = new FormData(e.target);
		const data = new URLSearchParams();
		let list = {};
		for (let field of formData) {
			const [key, value] = field;
			data.append(key, value);
			console.log(key, value);
		}
		data.append("test", "value");
		console.log(data);

		fetch(ACTION_URL, {
			method: "POST",
			body: JSON.stringify({id: "test", name: "test", status: "test", address: "test"}),
		});
	};
</script>

<div class="counter">
	<form on:submit={handleSubmit}>
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
				<td>
					<input type="text" name="id" />
				</td>
				<td><input type="text" name="name" /></td>
				<td><input type="text" name="status" /></td>
				<td><input type="text" name="address" /></td>
			</tr>
		</table>
		<input type="submit" value="Submit Changes" />
	</form>
</div>

<div>
	<button type="submit"> Submit Changes </button>
</div>
