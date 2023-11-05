<script lang="ts">
	import { writable } from "svelte/store";
	import { onMount } from "svelte";

	const HOSTNAME = "http://localhost:8080";
	const URL_GET_DATA = `${HOSTNAME}/getData`;
	const URL_NEW_DATA = `${HOSTNAME}/newData`;
	const URL_UPDATE_DATA = `${HOSTNAME}/updateData`;
	const URL_DELETE_DATA = `${HOSTNAME}/deleteRow`;

	interface filter {
		col: number;
		type: string;
		cont: string;
		operation: string;
	}

	interface topLevelInterface {
		[row: string]: dataInterface;
	}

	interface postDataInterface {
		rows: Array<dataInterface>;
	}

	interface dataInterface {
		[key: string]: string | number;
	}

	enum StatusValues {
		active = "Active",
		churned = "Churned",
		inquiry = "Inquiry",
		onboarding = "Onboarding",
	}

	enum Json {
		rows = "rows",
		data = "data",
	}

	enum Validation {
		name = "[A-Za-z\\s]+",
		address = "[A-Za-z0-9'\\.\\-\\s,]+",
	}

	enum DataColumns {
		id = "id",
		name = "name",
		status = "status",
		address = "address",
	}

	enum FilterColumns {
		col = "col",
		type = "type",
		cont = "cont",
		operation = "operation",
	}

	enum FilterState {
		OR = "OR",
		AND = "AND",
	}

	enum FilterTypes {
		equals = "equals",
		contains = "contains",
		less_than = "less than",
		greater_than = "greater than",
	}

	const filterTypes: Array<string> = [
		FilterTypes.equals,
		FilterTypes.contains,
		FilterTypes.less_than,
		FilterTypes.greater_than,
	];

	const columns: Array<string> = [
		DataColumns.id,
		DataColumns.name,
		DataColumns.status,
		DataColumns.address,
	];

	const filterColumns: Array<string> = [
		FilterColumns.col,
		FilterColumns.type,
		FilterColumns.cont,
		FilterColumns.operation,
	];

	const statusValues: Array<string> = [
		StatusValues.active,
		StatusValues.churned,
		StatusValues.inquiry,
		StatusValues.onboarding,
	];

	const emptyFilter = {
		col: 0,
		type: FilterTypes.equals,
		cont: "",
		operation: FilterState.OR,
	};

	export let tableData: topLevelInterface = {};
	export let filters: Array<filter> = [structuredClone(emptyFilter)];
	export let displayedData: topLevelInterface = {};

	let count: number = -1;
	let addRows: number = 0;

	onMount(() => getData());

	function toggleOperation(id: number): void {
		if (filters[id][FilterColumns.operation] === FilterState.OR) {
			filters[id][FilterColumns.operation] = FilterState.AND;
		} else {
			filters[id][FilterColumns.operation] = FilterState.OR;
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
		filters[id][FilterColumns.col] = text;
		filterData();
	}

	function updateFilterType(e: Event, id: number): void {
		let text: string = String((e.target as HTMLButtonElement).value);
		filters[id][FilterColumns.type] = text;
		filterData();
	}

	function updateFilterCont(e: Event, id: number): void {
		let text: string = (e.target as HTMLButtonElement).value;
		filters[id][FilterColumns.cont] = text;
		filterData();
	}

	function filterData(): void {
		displayedData = {}; // Start with no data
		let filtersUsed: boolean = false;

		for (var i = 0; i < filters.length; i++) {
			let filter: filter = filters[i];

			if (filter[FilterColumns.cont] === "") {
				// if there's no input ignore the filter
				continue;
			} else {
				filtersUsed = true;

				if (filter[FilterColumns.operation] === FilterState.AND) {
					// for AND filtering -> filter displayedData
					for (const key of Object.keys(displayedData)) {
						const column = columns[filter[FilterColumns.col]];
						if (
							filter[FilterColumns.type] ===
							FilterTypes.greater_than
						) {
							if (
								displayedData[key][column] >
								filter[FilterColumns.cont]
							) {
								delete displayedData[key];
							}
						}
						if (
							filter[FilterColumns.type] === FilterTypes.less_than
						) {
							if (
								displayedData[key][column] <
								filter[FilterColumns.cont]
							) {
								delete displayedData[key];
							}
						}
						if (
							filter[FilterColumns.type] === FilterTypes.contains
						) {
							if (
								!(displayedData[key][column] as String).match(
									filter[FilterColumns.cont]
								)
							) {
								delete displayedData[key];
							}
						}
						if (filter[FilterColumns.type] === FilterTypes.equals) {
							if (
								!(
									(tableData[key][column] as String) ==
									filter[FilterColumns.cont]
								)
							) {
								delete displayedData[key];
							}
						}
					}
				}
				if (
					filter[FilterColumns.operation] === FilterState.OR ||
					i === 0
				) {
					// for OR filtering -> filter tableData and combine the result with displayedData
					for (const key of Object.keys(tableData)) {
						const column = columns[filter[FilterColumns.col]];
						if (
							filter[FilterColumns.type] ===
							FilterTypes.greater_than
						) {
							if (
								tableData[key][column] >
								filter[FilterColumns.cont]
							) {
								displayedData[key] = tableData[key];
							}
						}
						if (
							filter[FilterColumns.type] === FilterTypes.less_than
						) {
							if (
								tableData[key][column] <
								filter[FilterColumns.cont]
							) {
								displayedData[key] = tableData[key];
							}
						}
						if (
							filter[FilterColumns.type] === FilterTypes.contains
						) {
							if (
								(tableData[key][column] as String).match(
									filter[FilterColumns.cont]
								)
							) {
								displayedData[key] = tableData[key];
							}
						}
						if (filter[FilterColumns.type] === FilterTypes.equals) {
							if (
								(tableData[key][column] as String) ==
								filter[FilterColumns.cont]
							) {
								displayedData[key] = tableData[key];
							}
						}
					}
				}
			}
		}
		if (!filtersUsed) {
			// If no filters were applied use tableData
			displayedData = structuredClone(tableData);
		}
	}

	function addFilter(): void {
		filters.push(emptyFilter);
		$: filters = filters;
	}

	// For checking if two rows of data are the same
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
				tableData = {};
				displayedData = {};
				for (const row of pgData[Json.data]) {
					tableData[row[DataColumns.id]] = row;
					displayedData[row[DataColumns.id]] = row;
				}
			})
			.catch((error) => {
				return [];
			});
	}

	async function basicPost(data: object, ACTION_URL: string): Promise<void> {
		await fetch(ACTION_URL, {
			method: "POST",
			body: JSON.stringify(data),
		})
			.then(() => {
				getData();
			})
			.then(() => {
				$: tableData = tableData;
				$: displayedData = tableData;
			});
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

	function validateData(type: string, value: string): boolean {
		if (type === DataColumns.name) {
			return /^([A-Za-z\s]+)$/.test(value);
		}
		if (type === DataColumns.address) {
			const regex = new RegExp(Validation.name);
			return /^([A-Za-z0-9'.-\s,]+)$/.test(value);
		}
		return true;
	}

	function subChanges(e: Event): void {
		const formData = new FormData((e.target as HTMLFormElement).form);

		let newDataPost: postDataInterface = { rows: [] };
		let newDataArray: Array<dataInterface> = [];

		let updatedDataArray: Array<dataInterface> = [];
		let updateDataPost: postDataInterface = { rows: [] };

		let unpackedData: Array<dataInterface> = [];
		let currentRowData: dataInterface = {};

		for (let field of formData) {
			// iterate thru all input fields of the table
			const [key, value] = field;
			if (!validateData(key, String(value))) {
				return;
			}
			if (Object.keys(currentRowData).includes(key)) {
				unpackedData.push(currentRowData);
				currentRowData = {};
			}
			currentRowData[key] = String(value);
		}
		unpackedData.push(currentRowData);

		newDataArray = unpackedData.slice(
			Object.keys(tableData).length,
			unpackedData.length
		);
		updatedDataArray = unpackedData.slice(0, Object.keys(tableData).length);

		if (newDataArray.length > 0) {
			// Evalute input fields that have new data
			newDataPost[Json.rows] = newDataArray;
			basicPost(newDataPost, URL_NEW_DATA);
		}

		for (var i = 0; i < updatedDataArray.length; i++) {
			const id = Object.keys(tableData)[i];
			// Evalute input fields that have updated data
			let reference: dataInterface = tableData[id];
			let currentRow = updatedDataArray[i];
			delete reference[DataColumns.id];
			if (!isEqual(reference, currentRow)) {
				currentRow[DataColumns.id] = id;
				updateDataPost[Json.rows].push(currentRow);
			}
		}
		basicPost(updateDataPost, URL_UPDATE_DATA);

		addRows = 0;
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
		<select on:change={(e) => updateFilterCol(e, rowI)}>
			{#each columns as name, colI}
				<option value={colI}>{name}</option>
			{/each}
		</select>
		<select on:change={(e) => updateFilterType(e, rowI)} name="filterType">
			{#each filterTypes as type, typeI}
				<option value={type}>{type}</option>
			{/each}
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
						<input type="submit" value={columnNames} />
					</th>
				{/each}
			</tr>
			{#if tableData}
				{#each Object.values(displayedData) as row, rowI}
					<tr>
						{#each Object.values(row) as column, colI}
							<td>
								{#if columns[colI] === DataColumns.id}
									{column}
								{:else if columns[colI] === DataColumns.address}
									<input
										name={columns[colI]}
										type="text"
										value={column}
										class="field"
										pattern={Validation.address}
									/>
								{:else if columns[colI] === DataColumns.name}
									<input
										name={columns[colI]}
										type="text"
										value={column}
										class="field"
										pattern={Validation.name}
									/>
								{:else if columns[colI] === DataColumns.status}
									<select name={columns[colI]}>
										{#each statusValues as status, statusI}
											<option value={status}
												>{status}</option
											>
										{/each}
									</select>
								{:else}
									<input
										name={columns[colI]}
										type="text"
										value={column}
									/>
								{/if}
							</td>
						{/each}
						<input
							type="submit"
							value="Delete Row"
							on:click={() => deleteRow(rowI)}
						/>
					</tr>{/each}
			{/if}
			{#each Array(addRows) as _, i}
				<tr>
					{#each columns as col, j}
						<td>
							{#if columns[j] === DataColumns.id}{:else if columns[j] === DataColumns.status}
								<select name={columns[j]}>
									{#each statusValues as status, statusI}
										<option value={status}>{status}</option>
									{/each}
								</select>
							{:else if columns[j] === DataColumns.address}
								<input
									type="text"
									value=""
									name={columns[j]}
									class="field"
									pattern={Validation.address}
									required
								/>
							{:else if columns[j] === DataColumns.name}
								<input
									type="text"
									name={columns[j]}
									value=""
									class="field"
									pattern={Validation.name}
									required
								/>
							{:else}
								<input type="text" name={col} />
							{/if}
						</td>
					{/each}
					<input
						type="submit"
						value="Delete Row"
						on:click={delRowNumber}
					/>
				</tr>
			{/each}
		</table>
		<input type="submit" value="Add Row" on:click={addRowNumber} />
		<input type="submit" value="Add Column" on:click={addRowNumber} />
		<br />
		<br />
		<input
			type="submit"
			value="Submit Changes"
			on:click={(e) => subChanges(e)}
		/>
	</form>
</div>
