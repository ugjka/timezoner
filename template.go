package main

var tmpl = `
<!doctype html>

<html lang="en">
<head>
	<meta charset="utf-8">
	<title>Timezoner</title>
	<meta name="description" content="Get timezone information for a given location">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body{
			margin: 40px auto;
			max-width: 650px;
			line-height: 1.6;
			font-size: 18px;
			color: #444;
			padding: 0 10px;
			background-color: #EEEEEE;
		}
		h1,h2,h3{
			line-height: 1.2;
		}
		table, th, td {
			border: 1px solid #444;
		}
		th {
			text-align: left;
			padding-left: 8px;
			padding-right: 8px;
		}
		td {
			padding: 8px;
		}
		input {
			background-color: #EEEDDD;
		}
		a {
			color: #006666;
			text-decoration: none;
		}
	</style>
</head>

<body>
	<h1>Timezoner</h1>
	<p>Get timezone information given by longitude and latitude</p>
	<h2>Try it!</h2>
	<form action="/">
		Longitude:<br>
		<input type="text" name="lon" value="{{if .Data}}{{.Data.Lon}}{{else}}0{{end}}"><br>
		Latitude:<br>
		<input type="text" name="lat" value="{{if .Data}}{{.Data.Lat}}{{else}}0{{end}}"><br>
		<input type="submit" value="Submit">
	</form>
	<br>
	{{if .Err}}<h2>Error: {{.Err}}</h2>{{end}}

	{{if .Data}}
	<table>
		<thead>
			<tr>
				<th>TZID</th>
				<th>Name</th>
				<th>Offset</th>
				<th>Time</th>
			<tr>
		</thead>
		<tbody>
			{{range .Data.Info}}
			<tr>
				<td>{{.TZID}}</td>
				<td>{{.Name}}</td>
				<td>{{.Offset}}</td>
				<td>{{.Time | FormatDate}}</td>
			</tr>
			{{end}}
		</tbody>
	</table>
	{{end}}

	{{if .Data}}<p><a href="/api/json?lon={{.Data.Lon}}&lat={{.Data.Lat}}">JSON API CALL</a></p>{{end}}
	<hr>
	<h3>About</h3>
	<p>I'm Uģis and i mess around with Go in my free time</p>
	<p>My github is <a href="https://github.com/ugjka">here</a>
	<p>Source code for this project <a href="https://github.com/ugjka/timezoner">here</a>
</body>
</html>
`
