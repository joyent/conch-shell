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

		<h3>Full Report</h3>
		<ul>
			<li><a href="/full">Text</a></li>
			<li><a href="/full.csv">CSV</a></li>
		</ul>

		<h3>Remediation Times</h3>
		<ul>
		{{ range .AzNames }}
			<li><a href="/reports/times/{{.}}">{{.}}</a></li>
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

var MboGraphsReportsIndex = `
<html>
	<body>
		<h1>Conch: Hardware Failures for {{ .Name }}</h1>

		<ul>
		{{ range $type, $data := .Data.TimesByType }}
			<li><a href="/reports/times/{{ $.Name }}/{{ $type }}">{{ $type }}</a><ul>
				<li>Mean: {{ $data.Mean }}</li>
				<li>Median: {{ $data.Median }}</li>
			</ul></li>
		{{ end }}
		</ul>

	</body>
</html>
`

var MboGraphsReportsBySubtype = `
<html>
	<body>
		<h1>Conch: Hardware Failures for {{.Az}}, Type {{.Name}} </h1>

		<ul>
		{{ range $type, $data := .Data }}
			<li>{{ $type }}<ul>
				<li>Mean: {{ $data.Mean }}</li>
				<li>Median: {{ $data.Median }}</li>
			</ul></li>
		{{ end }}
		</ul>

	</body>
</html>
`
