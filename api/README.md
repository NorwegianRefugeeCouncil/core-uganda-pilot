# nrc::core::api

### Case management API for Humanitarian Organizations

Part of the [NRC](http://nrc.no/) Core application. This project is the API layer for Core.

Humanitarian organizations often need to collect information about beneficiaries. Traditionnally, these organizations
use [kobo](https://github.com/kobotoolbox) or [ODK Central](https://getodk.org/) to do this. Though, these tools are not
well suited for beneficiary case management, but are rather designed for one-off data collection.

This project aims to solve some of the problems we identified with traditional beneficiary intake tools

### 1. Security & data compartmentalization

For large humanitarian organizations spanning multiple countries, it's essential that the field officers are only able
to access the data that they need. This project aims at providing the organization with customizable permission rules

For example, if a field officer needs to pull the beneficiaries information for offline use, he should only be able to
access that data.

### 2. Define and reuse data models

Organizations need to define the structure of the data they will collect. They often need to also reuse the same forms
throughout multiple campaigns. This project aims at providing the users a robust way of **defining** and **versioning**
data models throughout their lifecycle.

## Concepts

### 1. `FormElement`

A `FormElement` represents a single **Form Control** or **Form Group**. There are different types of `FormElements`,
such as `text`, `longtext`, `checkbox`, `radio`, etc. Each `FormElement` also contains the validation rules for that
element, such as `required`, `minLength` etc.

If the `FormElement` is a `FormGroup`, then the `FormElement` contains `children`. These `children` are recursively
rendered inside the form.

Example

```yaml
type: text
required: true
key: firstName
```

### 2. `FormSchema`

A `FormSchema` lives inside a `FormDefinition`. A `FormSchema` represents both the **visual layout** and **logical
layout** of the form. The `FormSchema` contains the different visual elements that the end-user will see when filling
out a form. Each `version` of a `FormDefinition` will have it's own `FormSchema`

Example:

```yaml

name: v1
schema:
  root:
    type: section
    children:
      - type: text
        key: firstName
      - type: checkbox
        key: active
```

### 1. `FormDefinition`

A `FormDefinition` represents a user-defined form definition, along with all the different versions of that form. For
example, a `FormDefinition` will contain the `v1` `v2` and `v3` versions along with their `FormSchema`.

Example:

```yaml

kind: FormDefinition
apiGroup: core.nrc.no/v1
metadata:
  name: GeneralIntake        // unique name for the definition
spec:
  group: my.custom.group
  names:
    plural: generalintakes   // will be available at /apis/my.custom.group/{v1|v2}/generalintakes
    singular: generalintake
    kind: GeneralIntake
  versions:
    - name: v1
      schema:
        root:
          type: section
          children:
            - type: text
              key: firstName
    - name: v2
      schema:
        root:
          type: section
          children:
            - type: text
              key: firstName
            - type: text
              key: lastName
```

### Implementation details

This API leverages the work by the [kubernetes](https://github.com/kubernetes/kubernetes) team, especially around
the `apimachinery` package. It also uses a [kine](https://github.com/k3s-io/kine) etcd3-compatible storage.

It essentially is a very light version of the kubernetes `apiserver`.

1. The user defines a new form by submitting a `FormDefinition` to the API server
2. The `FormDefinitionController` detects a new `FormDefinition`, and will create a matching
   `CustomResourceDefinition`
3. The `CustomResourceDefinition` controller will expose an API endpoint at the address
   `apis/{formDefinition.spec.group}/{formDefinition.spec.version[].name}/{formDefinition.spec.names.plural}`
   For example, the api endpoint could be `apis/acme.org/v1/superforms`. This endpoint allows the usual `GET`, `POST`
   , `PATCH`, `PUT`

### How to start this thing

Requirements:

```
# requires go 1.16

# We depend on code that lives in the kubernetes/kubernetes
# github repo. It lives under the /staging/src/k8s.io, which
# is not importable by go modules. We then download it locally,
# and have the go.mod redirects to point to the correct directory.
# This will download kubernetes/kubernetes under /third-party 

./hack/update-deps.sh

# Download go modules
go mod download

# Set up go vendor directory
go mod vendor

# Starts a kine instance with a local sqlite store 
go run ./cmd/kine 

# Starts the actual core api server
go run ./cmd/server --etcd-servers 127.0.0.1:2379

# Try the api a little bit
curl http://localhost:8001/apis/core.nrc.no/v1
curl http://localhost:8001/apis/core.nrc.no/v1/formdefinitions

# that will **soon** work as well
kubectl config set-cluster core --server http://localhost:8001
kubectl config set-context core --cluster core
kubectl config use-context core
kubectl get formdefinitions
...

```

### Example Form Definition

<details>
<summary>Example FormDefinition</summary>
<p>

# example.json

```

{
  "kind": "FormDefinition",
  "apiVersion": "core.nrc.no/v1",
  "metadata": {
    "name": "formtests.test.com"
  },
  "spec": {
    "group": "test.com",
    "names": {
      "plural": "formtests",
      "singular": "formtest",
      "kind": "FormTest"
    },
    "versions": [
      {
        "name": "v1",
        "storage": true,
        "served": true,
        "schema": {
          "formSchema": {
            "root": {
              "type": "section",
              "children": [
                {
                  "key": "firstName",
                  "label": [
                    {
                      "locale": "en",
                      "value": "First Name"
                    },
                    {
                      "locale": "fr",
                      "value": "Prenom"
                    }
                  ],
                  "description": [
                    {
                      "locale": "en",
                      "value": "Enter the first name of the beneficiary"
                    },
                    {
                      "locale": "fr",
                      "value": "Entrez le prénom du bénéficiaire"
                    }
                  ],
                  "type": "shortText",
                  "required": true
                },
                {
                  "key": "lastName",
                  "label": [
                    {
                      "locale": "en",
                      "value": "Last Name"
                    },
                    {
                      "locale": "fr",
                      "value": "Nom de famille"
                    }
                  ],
                  "description": [
                    {
                      "locale": "en",
                      "value": "Enter the first name of the beneficiary"
                    },
                    {
                      "locale": "fr",
                      "value": "Entrez le nom de famille du bénéficiaire"
                    }
                  ],
                  "type": "shortText",
                  "required": true
                }
              ]
            }
          }
        }
      }
    ]
  }
}

```

</p>
</details>  

#### Post the example form definition

```
curl --data-binary @example.json http://localhost:8001/apis/test.com/v1/formtests -H "Content-Type: application/json"
```

<details>
<summary>Example form submission</summary>

<p>

```
# submission.json
{
  "apiVersion": "test.com/v1",
  "kind": "FormTest",
  "metadata": {
    "name": "example-form-submission"
  },
  "spec": {
    "firstName": "Ludovic",
    "lastName": "Cleroux"
  }
}
```

</p>

</details>

#### Post the form submission

```
curl --data-binary @submission.json http://localhost:8001/apis/test.com/v1/formtests -H "Content-Type: application/json"
```

#### Get the form submission

```
# Get
curl http://localhost:8001/apis/test.com/v1/formtests/example-form-submission | jq
# List
curl http://localhost:8001/apis/test.com/v1/formtests/ | jq
```