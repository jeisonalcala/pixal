## Filesystem Methods

<details class="api_doc_details request_post">
	<summary><span class="method">POST</span>/filesystem/{path}</summary>
	<div>
		<h3>Description</h3>
		<p>
			Creates a new directory or uploads a file to an existing directory.
		</p>

		<h3>Parameters</h3>
		<p>
			The form parameters <b>must</b> be sent in the order displayed below
			for the realtime error checking to work. If 'name' comes after
			'file' it will be ignored.
		</p>
		<table>
			<tr>
				<td>Param</td>
				<td>Location</td>
				<td>Description</td>
			</tr>
			<tr>
				<td>type</td>
				<td>Form Values</td>
				<td>The type of node to create, can either be 'directory', or 'file'</td>
			</tr>
			<tr>
				<td>name</td>
				<td>Form Values</td>
				<td>
					Name of the directory to create, or of file to create. Not
					required if 'type' is 'file'
				</td>
			</tr>
			<tr>
				<td>file</td>
				<td>Form Values</td>
				<td>
					Multipart file to upload to the directory. Will be ignored
					if 'type' is 'directory'
				</td>
			</tr>
		</table>

		<h3>Returns</h3>
<pre>HTTP 200: OK
{
	"success": true,
	"id": "abc123" // ID of the newly uploaded file
}</pre>
		todo
	</div>
</details>

<details class="api_doc_details request_get">
	<summary><span class="method">GET</span>/filesystem/{path}</summary>
	<div>
		<h3>Description</h3>
		<p>
			Returns information about the requested path.
		</p>
		<h3>Parameters</h3>
		<table>
			<tr>
				<td>Param</td>
				<td>Required</td>
				<td>Location</td>
				<td>Description</td>
			</tr>
			<tr>
				<td>path</td>
				<td>true</td>
				<td>URL</td>
				<td>Path to the directory or file to request</td>
			</tr>
			<tr>
				<td>download</td>
				<td>false</td>
				<td>URL</td>
				<td>
					If the URL paramater '?download' is passed the requested
					file will be downloaded (if it is a file)
				</td>
			</tr>
		</table>
		<h3>Returns</h3>
		<h4>When the requested entity is a directory:</h4>
		<pre>HTTP 200: OK
{
	"success": true,
	"name": "some dir",
	"path": "/some dir",
	"type": "directory",
	"child_directories": [
		{
			"name": "some other directory",
			"type": "directory",
			"path": "/some dir/some other directory"
		}
	],
	"child_files": [
		{
			"name": "11. Lenny Kravitz - Fly away.ogg",
			"type": "file",
			"path": "/some dir/11. Lenny Kravitz - Fly away.ogg"
		}
	]
}</pre>
		<h4>When the requested entity is a file:</h4>
		<pre>HTTP 200: OK
{
	"success": true,
	"name": "11. Lenny Kravitz - Fly away.ogg",
	"path": "/some dir/11. Lenny Kravitz - Fly away.ogg",
	"type": "file",
	"file_info": {
		"success": true,
		"id": "Jf_u5TI9",
		"name": "11. Lenny Kravitz - Fly away.ogg",
		"date_upload": "2018-07-04T22:24:48Z",
		"date_last_view": "2018-07-04T22:24:48Z",
		"size": 9757269,
		"views": 0,
		"mime_type": "application/ogg",
		"thumbnail_href": "/file/Jf_u5TI9/thumbnail"
	}
}</pre>
	</div>
</details>

<details class="api_doc_details request_delete">
	<summary><span class="method">DELETE</span>/filesystem/{path}</summary>
	<div>
		<h3>Description</h3>
		<p>
			Deletes a filesystem node.
		</p>
		<h3>Parameters</h3>
		<table>
			<tr>
				<td>Param</td>
				<td>Required</td>
				<td>Location</td>
				<td>Description</td>
			</tr>
			<tr>
				<td>path</td>
				<td>true</td>
				<td>URL</td>
				<td>Path of the entity to delete</td>
			</tr>
		</table>
		<h3>Returns</h3>
<pre>HTTP 200: OK
{
	"success": true
}</pre>
	</div>
</details>
