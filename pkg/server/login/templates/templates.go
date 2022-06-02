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
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,400;0,500;0,700;1,400&display=swap" rel="stylesheet">
<title>Login to NRC Core</title>
<style>

body.nrc-theme {
	background-color: #FFFFFF;
	font-family: 'Roboto';
}

h1.nrc-theme {
	font-weight: 700;
	font-size: 44px;
	line-height: 52px;
	color: #000000;
	margin-bottom: 10px;
}
h2.nrc-theme {
	font-weight: 500;
	font-size: 24px;
	line-height: 29px;
	color: #000000;
	margin-bottom: 50px;
}
h3.nrc-theme {
	font-weight: 500;
	font-size: 16px;
	line-height: 19px;
	color: #000000;
	margin-bottom: 10px;
}
p.nrc-theme {
	font-weight: 400;
	font-size: 14px;
	line-height: 18px;
	color: #000000;
}

button.nrc-theme {
	background-color: #FF7602;
	border-color: #FF7602;
	color: #FFFFFF;
	height: 50px;
	width: 440px;
	font-size: 20px;
	line-height: 24px;
	font-weight: 500;
}

button.nrc-theme:hover,
button.nrc-theme:active,
button.nrc-theme:focus,
button.nrc-theme:not(:disabled):not(.disabled).active, button.nrc-theme:not(:disabled):not(.disabled):active, .show>button.nrc-theme.dropdown-toggle,
button.nrc-theme:not(:disabled):not(.disabled).active:focus, button.nrc-theme:not(:disabled):not(.disabled):active:focus, .show>button.nrc-theme.dropdown-toggle:focus{
	background-color: #FFCEAE;
	border-color: #FFCEAE;
	color: #FF7602;
    box-shadow: none;
}

button.btn-secondary.nrc-theme {
	font-style: normal;
	font-weight: 400;
	font-size: 18px;
	line-height: 26px;
	letter-spacing: 0.002em;
	text-decoration-line: underline;
	color: #1E4A7D;
	background-color: #FFFFFF;
	border-color: #FFFFFF;
	width: auto;
}

input.nrc-theme {
	height: 50px;
	width: 440px;
	border: 1px solid #666666;
	border-radius: 3px;
	font-weight: 400;
}

label.nrc-theme {
	font-weight: 500;
}

.form-font {
	font-size: 18px;
	line-height: 26px;
	color: #666666;
}

</style>
</head>
{{end}}

{{define "login_subject"}}
  <html>
  {{template "header"}}
    <body class="nrc-theme">
      <div class="d-flex flex-col vh-100 vw-100 align-items-center">

        <div class="container mx-5">
          <div class="row">
            <div class="col">
              <div class="mb-5">
				  <h1 class="nrc-theme">Login</h1>

				  <p class="nrc-theme">Access your personal account</p>

				  <form method="post">
					<div class="form-group mb-3">
					  <label class="form-label mb-1 nrc-theme form-font">Email</label>
					  <input
						class="form-control nrc-theme form-font"
						name="email"
						type="text" placeholder="Email Address" />
					</div>

					<button class="btn btn-primary nrc-theme" type="submit">
					  Next
					</button>

				  </form>

				{{ if .Error }}
				<div class="text-danger my-2 fw-bold">{{.Error}}</p>
				{{ end }}

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
  <body class="nrc-theme">
    <div class="d-flex flex-col vh-100 vw-100 align-items-center">
      <div class="container mx-5">
        <div class="row">
          <div class="col">
            <div class="mb-5">

                <h5>{{.OrganizationName}}</h5>

                {{ range $idp := .IdentityProviders}}
                  <form method="post" action="/login/oidc/{{ $idp.ID }}">
                    <button type="submit" class="btn btn-primary nrc-theme mb-2 w-100">Login with {{ $idp.Name }}</a>
                  </form>
                {{end}}

                {{ if .Error }}
                  <div class="text-danger my-2 fw-bold">{{.Error}}</p>
                {{ end }}

            </div>
          </div>
        </div>
      </div>
    </div>

  </body>
  </html>
{{end}}

{{define "challenge"}}
  <html>
  {{template "header"}}
  <body class="nrc-theme">
    <div class="d-flex flex-col vh-100 vw-100 align-items-center">
      <div class="container mx-5">
        <div class="row">
          <div class="col">
            <div class="mb-5">
			  <h1 class="nrc-theme">Access Consent</h1>
			  <h2 class="nrc-theme">NRC Core wants to access your account</h2>
			  <h3 class="nrc-theme">Make sure you trust NRC Core</h2>
				<p class="nrc-theme">
					You may be sharing sensitive information with this site or app.
					Learn about how NRC Core will handle your data by reviewing
					it’s terms of service and pricacy policies. You can always see
					or remove access in your account.
				</p>


                <form method="post">
					<div>
                  		<button type="submit" class="btn btn-success nrc-theme" formaction="/login/consent/approve">Next</button>
					</div>
					<div>
                  		<button type="submit" class="btn btn-secondary nrc-theme" formaction="/login/consent/decline">Cancel</button>
					</div>
                </form>

                {{ if .Error }}
                  <div class="text-danger my-2 fw-bold">{{.Error}}</p>
                {{ end }}

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
