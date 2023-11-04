<script lang="ts">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";

	interface filter {
		col: number;
		type: number;
		cont: string;
		operation: string;
	}

	interface topLevelInterface<row> {
		row?: dataInterface;
	}

	interface dataInterface {
		id: number;
		name: string;
		status: string;
		address: string;
	}

	let filterTypes: Array<string> = [
		"equals",
		"contains",
		"less than",
		"greater than",
	];

	export let tableData: topLevelInterface<number> = {};
	export let filters: Array<filter> = [
		{ col: 0, type: 0, cont: "", operation: "" },
	];
	export let displayedData: topLevelInterface<number> = {};
	export let columns: Array<string> = ["id", "name", "status", "address"];

	function toggleOperation(id: number) {
		if (filters[id]["operation"] == "OR") {
			filters[id]["operation"] = "AND";
		} else {
			filters[id]["operation"] = "OR";
		}
		filterData();
	}

	function removeFilter(id: number) {
		let temp: Array<filter> = [];
		for (var i = 0; i < filters.length; i++) {
			if (i != id) {
				temp.push(filters[i]);
			}
		}
		filters = temp;
		filterData();
	}

	function updateFilterCol(e: Event, id: number) {
		let text: number = Number((e.target as HTMLButtonElement).value);
		filters[id]["col"] = text;
		filterData();
	}

	function updateFilterType(e: Event, id: number) {
		let text: number = Number((e.target as HTMLButtonElement).value);
		filters[id]["type"] = text;
		filterData();
	}

	function updateFilterCont(e: Event, id: number) {
		let text: string = (e.target as HTMLButtonElement).value;
		filters[id]["cont"] = text;
		filterData();
	}

	function filterData() {
		displayedData = {};
		let filtersUsed: boolean = false;
		for (var i = 0; i < filters.length; i++) {
			let filter: filter = filters[i];
			if (filter["cont"] == "") {
				continue;
			} else {
				filtersUsed = true;
				if (filter["operation"] == "AND") {
					for (const key of Object.keys(displayedData)) {
						const column = columns[filter["col"]];
						if (filter["type"] == 1) {
							if (
								!displayedData[Number(key)][column].match(
									filter["cont"]
								)
							) {
								delete displayedData[key];
							}
						}
						if (filter["type"] == 0) {
							if (!tableData[key][column] == filter["cont"]) {
								delete displayedData[key];
							}
						}
					}
				}
				if (filter["operation"] == "OR" || i == 0) {
					for (const key of Object.keys(tableData)) {
						const column = columns[filter["col"]];
						if (filter["type"] == 1) {
							if (tableData[key][column].match(filter["cont"])) {
								displayedData[key] = tableData[key];
							}
						}
						if (filter["type"] == 0) {
							if (tableData[key][column] == filter["cont"]) {
								displayedData[key] = tableData[key];
							}
						}
					}
				}
			}
		}
		if (!filtersUsed) {
			displayedData = structuredClone(tableData);
		}
	}

	function addFilter() {
		let temp = structuredClone(filters);
		temp.push({ col: 0, type: 0, cont: "", operation: "OR" });
		$: filters = temp;
	}

	function isEqual(one: object, two: object) {
		let answer: boolean = true;
		let key1: Array<string> = Object.keys(one);
		let key2: Array<string> = Object.keys(two);
		if (key1.length != key2.length) {
			return false;
		}
		for (let index = 0; index < key1.length; index++) {
			let key: string = key1[index];
			if (key != key2[index]) {
				return false;
			}
			if (one[key] != two[key]) {
				return false;
			}
		}
		return true;
	}

	onMount(() => getData());
	async function getData() {
		fetch("http://localhost:8080/getData")
			.then((response) => response.json())
			.then((pgData) => {
				for (const row of pgData["data"]) {
					tableData[row["id"]] = row;
					displayedData[row["id"]] = row;
				}
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
	function delRowNumber(): number {
		addRows--;
		return addRows;
	}

	const handleSubmit = (e: Event) => {
		if (e["submitter"]["id"] == "subChanges") {
			const ACTION_URL = "http://localhost:8080/postData";
			const formData = new FormData(e.target);
			let data1: object = { rows: [] };
			let data2: object = {};
			let data3: object = { rows: [] };
			let temp: object = {};
			let currentId: number = -1;
			for (let field of formData) {
				const [key, value] = field;
				if (key == "id") {
					if (value == "") {
						currentId = -1;
					} else {
						currentId = Number(value);
					}
					if (currentId != -1) {
						temp[currentId] = {};
					}
				}
				if (currentId != -1) {
					temp[currentId][key] = value;
				}
			}
			for (let id of Object.keys(temp)) {
				if (!Object.keys(tableData).includes(id)) {
					data1["rows"].push(temp[id]);
				} else {
					data2[id] = temp[id];
				}
			}
			if (Object.keys(data1).length > 0) {
				postData(data1, ACTION_URL);
			}

			for (let id of Object.keys(data2)) {
				let reference: object = tableData[id];
				if (!isEqual(reference, data2[id])) {
					data3["rows"].push(data2[id]);
				}
			}
			postData(data3, "http://localhost:8080/updateData");
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
	<h1>Filters</h1>
	{#each filters as filter, rowI}
		{#if rowI > 0}
			<input
				id={String(rowI)}
				on:click={() => toggleOperation(rowI)}
				type="submit"
				value={filters[rowI]["operation"]}
			/>
		{/if}
		<br />
		<select
			id={String(rowI)}
			name="colName"
			on:change={(e) => updateFilterCol(e, rowI)}
		>
			{#each columns as name, colI}
				<option value={colI}>{name}</option>
			{/each}
		</select>
		<select
			on:change={(e) => updateFilterType(e, rowI)}
			id={String(rowI)}
			name="filterType"
		>
			<option value="0">equals</option>
			<option value="1">contains</option>
			<option value="2">less than</option>
			<option value="3">greater than</option>
		</select>
		<input
			on:change={(e) => updateFilterCont(e, rowI)}
			id={String(rowI)}
			type="text"
			value=""
		/>
		<input
			id={String(rowI)}
			type="submit"
			value="Remove Filter"
			on:click={() => removeFilter(rowI)}
		/>
		<br />
	{/each}

	<input type="submit" value="Add Filter" on:click={addFilter} />
</div>
<br />
<br />

<div>
	<form on:submit|preventDefault={handleSubmit}>
		<table>
			<tr>
				{#each columns as columnNames, index}
					<th>
						<input
							id="subChanges"
							type="submit"
							value={columnNames + " ⬆"}
						/>
					</th>
				{/each}
			</tr>
			{#if tableData}
				{#each Object.values(displayedData) as row, rowI}
					<tr>
						{#each Object.values(row) as column, colI}
							<td
								><input
									name={columns[colI]}
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
					<input
						id={String(i)}
						type="submit"
						value="Delete Row"
						on:click={delRowNumber}
					/>
				</tr>
			{/each}
		</table>
		<input id="addRow" type="submit" value="Add Row" />
		<br />
		<br />
		<input id="subChanges" type="submit" value="Submit Changes" />
	</form>
</div>
