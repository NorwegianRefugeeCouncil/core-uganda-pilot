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

{{define "navbar"}}
<nav class="navbar navbar-expand-sm navbar-light bg-light">
  <a class="navbar-brand" href="#">Core Admin</a>
  <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
    <span class="navbar-toggler-icon"></span>
  </button>
  <div class="collapse navbar-collapse" id="navbarSupportedContent">
    <ul class="navbar-nav mr-auto">
      <li class="nav-item">
        <a class="nav-link" href="/admin/clients">OAuth2 Clients</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="/admin/organizations">Organizations</a>
      </li>
    </ul>
  </div>
</nav>
{{end}}

{{define "client_list"}}
	<html>
	{{template "header"}}
	<body class="bg-dark">
	{{template "navbar"}}
	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card my-4">
	<div class="card-body">
	<h5 class="card-title">Oauth2 Clients</h5>
	{{if not .Clients}}
		No clients
	{{else}}
		<div class="list-group mb-3">
		{{range $client := .Clients}}
			<a class="list-group-item list-group-item-action" href="/admin/clients/{{$client.ClientID}}">{{$client.ClientName}}</a>
		{{end}}
		</div>
	{{end}}

	<p><a class="btn btn-primary" href="/admin/clients/add">Add Client<a/></p>
	</div>
	</div>
	</div>
	</div>
	</div>
	</body>
	</html>
{{end}}

{{define "client_add"}}
	<html>
	{{template "header"}}
	<body class="bg-dark">
	{{template "navbar"}}
	<form method="post" {{ if not .IsNew}}action="/admin/clients/{{.ClientID}}"{{end}}>
	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card my-4">
	<div class="card-body">
	<p><a class="card-title" href="/admin/clients">Back to list</a></p>

	{{if .ClientSecret}}
		<div class="p-4 my-3 border shadow bg-warning">
		<p class="text-danger">Save this secret! You will not be able to retrieve it later!</p>
		<div class="d-flex flex-row">
		<script>
		function toggle(){
			const input = document.getElementById('client_secret')
			const btn = document.getElementById('client_secret_toggle')
			btn.textContent = btn.textContent === 'Show' ? 'Hide' : 'Show'
			input.type = input.type === 'text' ? 'password' : 'text'
		}
		function copy(){
			const input = document.getElementById('client_secret')
			const oldType = input.type
			input.type = "text"
			input.select()
			document.execCommand("copy");
			input.type = oldType
		}
		</script>
		<input
			id="client_secret"
			class="form-control text-monospace me-2"
			type="password"
			value="{{.ClientSecret}}"/>
		<button
			id="client_secret_toggle"
			type="button"
			onclick="toggle()"
			class="btn btn-outline-primary me-2">Show</button>
		<button
			id="client_secret_copy"
			type="button"
			onclick="copy()"
			class="btn btn-outline-primary">Copy</button>
		</div>
		</div>
	{{end}}

	<div class="form-group mb-4">
		<label class="form-label">Client ID</label>
		<input class="form-control"
			   type="text"
			   value="{{.ClientID}}"
               name="client_id"
               placeholder="client id"
               {{if not .IsNew}}disabled="disabled"{{end}}/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Client Name</label>
		<input class="form-control" type="text" value="{{.ClientName}}" name="client_name" placeholder="client name"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Client URI</label>
		<input class="form-control" type="text" value="{{.ClientURI}}" name="client_uri" placeholder="client uri"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Grant Types</label>
		<input class="form-control" type="text" value="{{.GrantTypes}}" name="grant_types" placeholder="grant types"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Scope</label>
		<input class="form-control" type="text" value="{{.Scope}}" name="scope" placeholder="scope"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Response Types</label>
		<input class="form-control" type="text" value="{{.ResponseTypes}}" name="response_types" placeholder="response types"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Token Endpoint Auth Method</label>
		<input class="form-control" type="text" value="{{.TokenEndpointAuthMethod}}" name="token_endpoint_auth_method" placeholder="Token Endpoint Auth Method"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Redirect URIs</label>
		<input class="form-control" type="text" value="{{.RedirectURIs}}" name="redirect_uris" placeholder="redirect uris"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Allowed CORS Origins</label>
		<input class="form-control" type="text" value="{{.AllowedCORSOrigins}}" name="allowed_cors_origins" placeholder="allowed cors origins"/>
	</div>

	<button class="btn btn-primary me-2" type="submit">{{ if .IsNew }}Create Client{{ else }}Save Client{{ end }}</button>
	{{ if not .IsNew}}<button type="submit" class="btn btn-danger" type="submit" formaction="/admin/clients/{{.ClientID}}/delete">Delete Client</button>{{end}}

	{{if .Error}}
		<div class="text-danger p-2 fw-bold">{{.Error}}</div>
	{{end}}
	</div>
	</div>
	</div>
	</div>
	</div>
	</form>
	</body>
	</html>
{{end}}

{{define "organization_list"}}
	<html>
	{{template "header"}}
	<body class="bg-dark">
	{{template "navbar"}}
	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card my-4">
	<div class="card-body">
	<h5 class="card-title">Organizations</h5>
	{{if not .Organizations}}
		No organizations
	{{else}}
		<div class="list-group mb-3">
		{{range $organization := .Organizations}}
			<a class="list-group-item list-group-item-action" href="/admin/organizations/{{$organization.ID}}">{{$organization.Name}}</a>
		{{end}}
		</div>
	{{end}}

	{{if .Error}}
		<div class="text-danger p-2 fw-bold">{{.Error}}</div>
	{{end}}


	<p><a class="btn btn-primary" href="/admin/organizations/add">Add Organization<a/></p>
	</div>
	</div>
	</div>
	</div>
	</div>
	</body>
	</html>
{{end}}

{{define "organization_add"}}
	<html>
	{{template "header"}}
	<body class="bg-dark">
	{{template "navbar"}}
	<form method="post" {{ if not .IsNew}}action="/admin/organizations/{{.ID}}"{{end}}>
	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card my-4">
	<div class="card-body">
	<p><a class="card-title" href="/admin/organizations">Back to list</a></p>

	<div class="form-group mb-4">
		<label class="form-label">Organization ID</label>
		<input class="form-control"
			   type="text"
			   value="{{.ID}}"
               name="id"
               placeholder="Organization ID"
               disabled="disabled"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Organization Name</label>
		<input class="form-control" type="text" value="{{.Name}}" name="organization_name" placeholder="Organization Name"/>
	</div>

	<div class="card mb-3">
	<div class="card-body">
	<h5 class="card-title">Identity Providers</h5>
	{{if not .IdentityProviders}}
		<div class="mb-3">
			No identity providers
		</div>
	{{else}}
		<div class="list-group mb-3">
		{{range $identityProvider := .IdentityProviders}}
			<a class="list-group-item list-group-item-action" href="/admin/organizations/{{.ID}}/identity-providers/{{$identityProvider.ID}}">{{$identityProvider.Name}}</a>
		{{end}}
		</div>
	{{end}}
	<p><a class="btn btn-primary" href="/admin/organizations/{{.ID}}/identity-providers/add">Add Identity Provider<a/></p>
	</div>
	</div>

	<button class="btn btn-primary me-2" type="submit">{{ if .IsNew }}Create Organization{{ else }}Save Organization{{ end }}</button>
	{{ if not .IsNew}}<button type="submit" class="btn btn-danger" type="submit" formaction="/admin/organizations/{{.ID}}/delete">Delete Organization</button>{{end}}

	{{if .Error}}
		<div class="text-danger p-2 fw-bold">{{.Error}}</div>
	{{end}}
	</div>
	</div>
	</div>
	</div>
	</div>
	</form>
	</body>
	</html>
{{end}}

{{define "idp_add"}}
	<html>
	{{template "header"}}
	<body class="bg-dark">
	{{template "navbar"}}
	<form method="post" {{ if not .IsNew}}action="/admin/organizations/{{.OrganizationID}}/identity-providers/{{.ID}}"{{end}}>
	<div class="container">
	<div class="row">
	<div class="col">
	<div class="card my-4">
	<div class="card-body">
	<p><a class="card-title" href="/admin/organizations/{{.OrganizationID}}">Back to Organization</a></p>

	<div class="form-group mb-4">
		<label class="form-label">Client ID</label>
		<input class="form-control"
			   type="text"
			   value="{{.ClientID}}"
               name="client_id"
               placeholder="client id"
               {{if not .IsNew}}disabled="disabled"{{end}}/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Client Secret</label>
		<input class="form-control"
			   type="text"
			   value="{{.ClientSecret}}"
               name="client_secret"
               placeholder="Client Secret"/>
	</div>

	<div class="form-group mb-4">
		<label class="form-label">Issuer</label>
		<input class="form-control" type="text" value="{{.Issuer}}" name="issuer" placeholder="Issuer"/>
	</div>

	<button class="btn btn-primary me-2" type="submit">{{ if .IsNew }}Create Identity Provider{{ else }}Save Identity Provider{{ end }}</button>
	{{ if not .IsNew}}<button type="submit" class="btn btn-danger" type="submit" formaction="/admin/organizations/{{.OrganizationID}}/identity-providers/{{.ID}}/delete">Delete Identity Provider</button>{{end}}

	{{if .Error}}
		<div class="text-danger p-2 fw-bold">{{.Error}}</div>
	{{end}}
	</div>
	</div>
	</div>
	</div>
	</div>
	</form>
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
