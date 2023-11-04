<script lang="ts">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";

	const HOSTNAME = "http://localhost:8080";
	const URL_GET_DATA = `${HOSTNAME}/getData`;
	const URL_POST_DATA = `${HOSTNAME}/postData`;
	const URL_UPDATE_DATA = `${HOSTNAME}/updateData`;
	const URL_DELETE_DATA = `${HOSTNAME}/deleteRow`;

	interface filter {
		col: number;
		type: number;
		cont: string;
		operation: string;
	}

	interface topLevelInterface<T> {
		[row: string]: T;
	}

	interface postDataInterface {
		rows: Array<dataInterface>
	}

	interface dataInterface {
		[key: string]: string;
	}

	let filterTypes: Array<string> = [
		"equals",
		"contains",
		"less than",
		"greater than",
	];

	export let tableData: topLevelInterface<dataInterface> = {};
	export let filters: Array<filter> = [
		{ col: 0, type: 0, cont: "", operation: "" },
	];
	export let displayedData: topLevelInterface<dataInterface> = {};
	export let columns: Array<string> = ["id", "name", "status", "address"];

	let count: number = -1;
	let addRows: number = 1;

	onMount(() => getData());

	function toggleOperation(id: number): void {
		if (filters[id]["operation"] == "OR") {
			filters[id]["operation"] = "AND";
		} else {
			filters[id]["operation"] = "OR";
		}
		filterData();
	}

	function removeFilter(id: number): void {
		let temp: Array<filter> = [];
		for (var i = 0; i < filters.length; i++) {
			if (i != id) {
				temp.push(filters[i]);
			}
		}
		filters = temp;
		filterData();
	}

	function updateFilterCol(e: Event, id: number): void {
		let text: number = Number((e.target as HTMLButtonElement).value);
		filters[id]["col"] = text;
		filterData();
	}

	function updateFilterType(e: Event, id: number): void {
		let text: number = Number((e.target as HTMLButtonElement).value);
		filters[id]["type"] = text;
		filterData();
	}

	function updateFilterCont(e: Event, id: number): void {
		let text: string = (e.target as HTMLButtonElement).value;
		filters[id]["cont"] = text;
		filterData();
	}

	function filterData(): void {
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
								!(displayedData[key][column] as String).match(
									filter["cont"]
								)
							) {
								delete displayedData[key];
							}
						}
						if (filter["type"] == 0) {
							if (
								!(
									(tableData[key][column] as String) ==
									filter["cont"]
								)
							) {
								delete displayedData[key];
							}
						}
					}
				}
				if (filter["operation"] == "OR" || i == 0) {
					for (const key of Object.keys(tableData)) {
						const column = columns[filter["col"]];
						if (filter["type"] == 1) {
							if (
								(tableData[key][column] as String).match(
									filter["cont"]
								)
							) {
								displayedData[key] = tableData[key];
							}
						}
						if (filter["type"] == 0) {
							if (
								(tableData[key][column] as String) ==
								filter["cont"]
							) {
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

	function addFilter(): void {
		let temp = structuredClone(filters);
		temp.push({ col: 0, type: 0, cont: "", operation: "OR" });
		$: filters = temp;
	}

	function isEqual(one: dataInterface, two: dataInterface): boolean {
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

	async function getData(): Promise<void> {
		fetch(URL_GET_DATA)
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

	async function basicPost(data: object, ACTION_URL: string): Promise<void> {
		fetch(ACTION_URL, {
			method: "POST",
			body: JSON.stringify(data),
		});
		window.location.reload();
		//let test = structuredClone(tableData);
		//test.push(data);
		//$: tableData = test;
	}

	function addRowNumber(): number {
		addRows++;
		return addRows;
	}

	function delRowNumber(): number {
		addRows--;
		return addRows;
	}

	function deleteRow(rowId: number): void {
		const actualId: string = Object.keys(displayedData)[rowId];
		basicPost({ rowNum: actualId }, URL_DELETE_DATA);
	}

	function subChanges(e: Event) {
		const formData = new FormData((e.target as HTMLFormElement).form);
		let data1: postDataInterface = { rows: [] };
		let data2: topLevelInterface<dataInterface> = {};
		let data3: postDataInterface = { rows: [] };
		let temp: topLevelInterface<dataInterface> = {};
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
				temp[currentId][key] = String(value);
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
			basicPost(data1, URL_POST_DATA);
		}

		for (let id of Object.keys(data2)) {
			let reference: object = tableData[id];
			if (!isEqual(reference, data2[id])) {
				data3["rows"].push(data2[id]);
			}
		}
		basicPost(data3, URL_UPDATE_DATA);
	}
</script>

<div>
	<h1>Filters</h1>
	{#each filters as filter, rowI}
		{#if rowI > 0}
			<input
				on:click={() => toggleOperation(rowI)}
				type="submit"
				value={filters[rowI]["operation"]}
			/>
		{/if}
		<br />
		<select name="colName" on:change={(e) => updateFilterCol(e, rowI)}>
			{#each columns as name, colI}
				<option value={colI}>{name}</option>
			{/each}
		</select>
		<select on:change={(e) => updateFilterType(e, rowI)} name="filterType">
			<option value="0">equals</option>
			<option value="1">contains</option>
			<option value="2">less than</option>
			<option value="3">greater than</option>
		</select>
		<input
			on:change={(e) => updateFilterCont(e, rowI)}
			type="text"
			value=""
		/>
		<input
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
	<form>
		<table>
			<tr>
				{#each columns as columnNames, index}
					<th>
						<input type="submit" value={columnNames + " ⬆"} />
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
							type="submit"
							value="Delete Row"
							on:click={() => deleteRow(rowI)}
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
						type="submit"
						value="Delete Row"
						on:click={delRowNumber}
					/>
				</tr>
			{/each}
		</table>
		<input type="submit" value="Add Row" on:click={addRowNumber} />
		<br />
		<br />
		<input
			type="submit"
			value="Submit Changes"
			on:click={(e) => subChanges(e)}
		/>
	</form>
</div>
