<script lang="ts">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";
	let columns: Array<string> = ["ID", "Name", "Status", "Address"];
	export let tableData: Array<object> = [];

	onMount(() => getData());
	async function getData() {
		fetch("http://localhost:8080/getData")
			.then((response) => response.json())
			.then((pgData) => {
				tableData = pgData["data"];
			})
			.catch((error) => {
				return [];
			});
	}
	async function postData(data: object, ACTION_URL: string) {
		fetch(ACTION_URL, {
			method: "POST",
			body: JSON.stringify(data),
		});
		window.location.reload();
		//let test = structuredClone(tableData);
		//test.push(data);
		//$: tableData = test;
	}
	async function deleteRow(rowNum: number, ACTION_URL: string) {
		fetch(ACTION_URL, {
			method: "POST",
			body: JSON.stringify({ rowNum: rowNum }),
		});
		window.location.reload();
		//let test = structuredClone(tableData);
		//test.push(data);
		//$: tableData = test;
	}

	let count: number = -1;
	function getNumber(): string {
		count++;
		let answer = "delete" + String(count);
		return answer;
	}

	let addRows: number = 1;
	function addRowNumber(): number {
		addRows++;
		return addRows;
	}

	const handleSubmit = (e: Event) => {
		if (e["submitter"]["id"] == "subChanges") {
			const ACTION_URL = "http://localhost:8080/postData";
			const formData = new FormData(e.target);
			let data1: object = { rows: [] };
			let temp: object = {};
			let currentId: number = 0;
			for (let field of formData) {
				const [key, value] = field;

				if (key == "id") {
					currentId = Number(value);
					temp[currentId] = {};
				}
				temp[currentId][key] = value;
			}
			for (let row of Object.values(temp)) {
				data1["rows"].push(row);
			}
			postData(data1, ACTION_URL);
		}
		if (e["submitter"]["id"].slice(0, 6) == "delete") {
			const rowNum = e["submitter"]["id"].slice(6);
			const rowId = tableData[rowNum]["id"];
			deleteRow(rowId, "http://localhost:8080/deleteRow");
		}
		if (e["submitter"]["id"] == "addRow") {
			addRowNumber();
		}
	};
</script>

<div>
	<form on:submit|preventDefault={handleSubmit}>
		<table>
			<tr>
				{#each columns as columnNames, index}
					<th>{columnNames}</th>
				{/each}
			</tr>
			{#if tableData}
				{#each tableData as row, index}
					<tr>
						{#each Object.values(row) as column}
							<td
								><input
									id="delete"
									type="text"
									value={column}
								/></td
							>
						{/each}
						<input
							id={getNumber()}
							type="submit"
							value="Delete Row"
						/>
					</tr>
				{/each}
			{/if}
			{#each Array(addRows) as _, i}
				<tr>
					<td><input type="text" name="id" /></td>
					<td><input type="text" name="name" /></td>
					<td><input type="text" name="status" /></td>
					<td><input type="text" name="address" /></td>
				</tr>
			{/each}
		</table>
		<input id="addRow" type="submit" value="Add Row" />
		<br />
		<br />
		<input id="subChanges" type="submit" value="Submit Changes" />
	</form>
</div>
