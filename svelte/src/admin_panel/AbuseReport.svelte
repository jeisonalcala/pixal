<script>
import { formatDate } from "../util/Formatting.svelte";
import Expandable from "../util/Expandable.svelte";
import { createEventDispatcher } from "svelte";
let dispatch = createEventDispatcher()

export let report
export let ip_report_count
let preview = false

$: can_grant = report.status !== "granted"
$: can_reject = report.status !== "rejected"

let set_status = async (action, report_type) => {
	dispatch("resolve_report", {action: action, report_type: report_type})
}
</script>

<Expandable expanded={report.status === "pending"} click_expand>
	<div slot="header" class="header">
		<div class="icon_cell">
			<img class="file_icon" src={"/api/file/"+report.file.id+"/thumbnail"} alt="File thumbnail"/>
		</div>

		<div class="title">{report.file.name}</div>
			<div class="stats">Type<br/>
				{report.file.abuse_type === "" ? report.type : report.file.abuse_type}
			</div>
		{#if report.status !== "pending"}
			<div class="stats">Status<br/>{report.status}</div>
		{/if}
		<div class="stats">R<br/>{report.reports.length}</div>
		<div class="stats">V<br/>{report.file.views}</div>
		<div class="stats">DL<br/>{Math.round(report.file.bandwidth_used / report.file.size)}</div>
	</div>
	<div class="details">
		<div class="toolbar">
			<div style="flex: 1 1 auto">
				<a class="button" target="_blank" href={"/u/"+report.file.id} rel="noreferrer">
					<i class="icon">open_in_new</i> Open file
				</a>
				<button class:button_highlight={preview} on:click={() => {preview = !preview}}>
					<i class="icon">visibility</i> Preview
				</button>
				{#if can_grant}
					<button class="button_highlight" on:click={() => {set_status("grant", report.type)}}>
						<i class="icon">done</i> Block ({report.type})
					</button>
				{/if}
				{#if can_reject}
					<button class="button_red" on:click={() => {set_status("reject", "")}}>
						<i class="icon">delete</i> Ignore
					</button>
				{/if}
			</div>
			<div style="flex: 0 1 auto">
				<button on:click={() => {set_status("grant", "terrorism")}}>terrorism</button>
				<button on:click={() => {set_status("grant", "gore")}}>gore</button>
				<button on:click={() => {set_status("grant", "child_abuse")}}>child_abuse</button>
				<button on:click={() => {set_status("grant", "malware")}}>malware</button>
				<button on:click={() => {set_status("grant", "doxing")}}>doxing</button>
				<button on:click={() => {set_status("grant", "revenge_porn")}}>revenge_porn</button>
				<button on:click={() => {set_status("grant", "copyright")}}>copyright</button>
			</div>
		</div>
		<div style="text-align: center;">
			{#if preview}
			<br/>
			<iframe
				title="File preview"
				src="/u/{report.file.id}?embed"
				style="border: none; width: 100%; height: 600px; border-radius: 16px;"
			></iframe>
			{/if}
		</div>
		<div class="table_scroll">
			<table>
				<tr>
					<td>Time</td>
					<td>IP</td>
					<td>Type</td>
					<td>Status</td>
					<td colspan="2">Reports from this IP</td>
				</tr>
				{#each report.reports as user_report}
					<tr>
						<td>{formatDate(user_report.time, true, true, false)}</td>
						<td>{user_report.ip_address}</td>
						<td>{user_report.type}</td>
						<td>{user_report.status}</td>
						<td>{ip_report_count[user_report.ip_address]}</td>
						<td>
							{#if can_grant}
								<button on:click={() => dispatch("resolve_by_ip", {ip: user_report.ip_address, action: "grant"})}>
									Accept all
								</button>
							{/if}
							{#if can_reject}
								<button on:click={() => dispatch("resolve_by_ip", {ip: user_report.ip_address, action: "reject"})}>
									Ignore all
								</button>
							{/if}
						</td>
					</tr>
					{#if user_report.description !== ""}
						<tr>
							<td>Description</td>
							<td colspan="5" style="white-space: pre-wrap; word-wrap: break-word;">
								{user_report.description}
							</td>
						</tr>
					{/if}
				{/each}
			</table>
		</div>
	</div>
</Expandable>

<style>
.header {
	display: flex;
	flex-direction: row;
	line-height: 1.2em;
}
.icon_cell {
	flex: 0 0 auto;
	line-height: 0;
	margin-right: 6px;
}
.title {
	flex: 1 1 auto;
	align-self: center;
	word-break: break-all;
}
.stats {
	flex: 0 0 auto;
	padding: 0 4px;
	border-left: 1px solid var(--separator);
	text-align: center;
}
.details {
	flex: 1 1 auto;
}
.toolbar {
	display: flex;
	flex-direction: row;
}
.file_icon {
	width:48px;
	height:48px;
}
</style>
