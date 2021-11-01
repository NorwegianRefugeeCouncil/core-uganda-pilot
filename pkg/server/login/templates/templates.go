package templates

import (
	"html/template"
)

const tplText = `

{{define "header"}}
<head>
<!-- Required meta tags -->
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
<!-- Bootstrap CSS -->
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
<title>Hello, world!</title>
</head>
{{end}}

{{define "login_subject"}}
	<html>
	{{template "header"}}
	<body class="bg-light">
	<div class="d-flex flex-col vh-100 vw-100 align-items-center">

	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card shadow mb-5">
	<div class="card-body">
	<h5 class="card-title">Login</h5>

	<form method="post">

	<input type="hidden" name="login_challenge" value="{{.LoginChallenge}}"/>

	<div class="form-group mb-3">
	<label class="form-label">Enter your email address</label>
	<input class="form-control" name="email" type="text" placeholder="Email Address"/>
	</div>

	<button class="btn btn-primary" type="submit">Login</button>

	</form>

	{{ if .Error }}
    <div class="text-danger my-2 fw-bold">{{.Error}}</p>
	{{ end }}

	</div>
	</div>
	</div>
	</div>
	</div>
	</div>

	</body>
	</html>
{{end}}

{{define "login_idp"}}
	<html>
	{{template "header"}}
	<body class="bg-light">
	<div class="d-flex flex-col vh-100 vw-100 align-items-center">
	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card shadow mb-5">
	<div class="card-body">

	<h5 class="card-title">{{.OrganizationName}}</h5>

	{{ range $idp := .IdentityProviders}}
		<a class="btn btn-primary mb-2 w-100" href="#">Login with {{ $idp.Name }}</a>
	{{end}}

	{{ if .Error }}
    <div class="text-danger my-2 fw-bold">{{.Error}}</p>
	{{ end }}

	</div>
	</div>
	</div>
	</div>
	</div>
	</div>

	</body>
	</html>
{{end}}
`

var Template *template.Template = nil

func init() {
	var err error
	Template, err = template.New("").Parse(tplText)
	if err != nil {
		panic(err)
	}
}
