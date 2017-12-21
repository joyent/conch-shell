// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package templates

var MboGraphsIndex = `
<html>
	<body>
		<h1>Conch : MBO Hardware Failures</h1>
		<h2>Text Reports</h2>
		<ul>
		{{ range .AzNames }}
			<li><a href="/reports/{{.}}">{{.}}</a></li>
		{{ end }}
		</ul>


		<h2>Graphs</h2>
		<h3>By Type</h3>
		<ul>
		{{ range .AzNames }}
			<li><a href="/graphics/{{.}}/by_type.png">{{.}}</a></li>
		{{ end }}
		</ul>

		<h3>By Vendor</h3>
		<ul>
		{{ range .AzNames }}
			<li><a href="/graphics/{{.}}/by_vendor.png">{{.}}</a></li>
		{{ end }}
		</ul>
	</body>
</html>
`
