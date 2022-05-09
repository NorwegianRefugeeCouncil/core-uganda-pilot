# Backend Standards

## How to create a new microservice

TODO

## Microservice structure

```
pkg/
  server/
    my-microservice/
      controllers/
        my_resource/
          controller.go
          create.go
          list.go
          get.go
          update.go
          delete.go
      services/
        my_resource/
          service.go
          create.go
          list.go
          get.go
          update.go
          delete.go
      models/
        my_resource/
          model.go
        migrate.go
      migrations/
        1_description.up.go
        1_description.down.go
        2_description.up.go
        2_description.down.go
      types/
      server.go
```

### Controllers

For a RESTful resource create a directory under `controllers`.
`controller.go` sets up the routing for the resource.

For each action create a file, eg. `create.go`.

The controller action is responsible for:
- Unmarshaling the request
- Validating the request payload
- Calling the appropriate methods in the service or model
- Validating the response payload
- Marshaling the response
- Sending the response

### Services

The service layer is an optional layer that should be used for handling more complex business logic and orchestration.
For example, creating a resource may involve creating multiple database resources. This can be managed in a service.
But fetching a single resource from the database may not require a service (in this case the service would just be proxying the call to the model from the controller).

- Start a transaction if required
- Create uuids when creating new resources
- Call the appropriate methods in the model or other services
- Returns an internal struct that can be marshaled to the response by the controller or used in the calling service

```
func (d *fooService) Create(ctx context.Context, foo types.Foo) (*types.Foo, error) {
	tx := d.transactionStore.Begin()

	entityId := uuid.NewV4()

  foo.Id = entityId

  err := d.fooStore.Create(tx, foo)

  if err != nil {
    tx.Rollback()
    return nil, err
  }

  ...other stuff

  tx.Commit()

  return &foo, nil
}
```

### Models

The model layer is responsible for taking an internal struct and writing it to the database.

In `model.go` we will define a generic interface (`FooModel`) that will be implemented by the model for a given database technology (`fooPostgresModel`). 

All model methods should take an optional transaction as an argument.
If passed in this will be used to exectute statements, otherwise fetch the db from the factory.

- Resolve the db connection from either the transaction or factory
- Build a query statement using a query builder (eg. SQLBuilder)
- Execute the statement
- Return the result as an internal struct

```
func (d *fooPostgresStore) InsertFoo(ctx context.Context, db *gorm.DB, foo types.Foo) (*types.Foo, error) {
	if db == nil {
		var err error
		db, err = d.db.Get()

		if err != nil {
			return nil, err
		}
	}

	ddl := d.sqlBuilder.InsertRow(
		"public",
		"foo",
		[]string{
			"id",
			"name",
		},
		[]any{
			entity.ID,
			entity.Name,
		},
	)

	result := db.Exec(ddl.Query, ddl.Args...)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &foo, nil
}
```

### Types

This directory contains all types used within the microservice.
Each type should have it's own file.
The file should also include these methods if required:

- UnmarshalJSON
- MarshalJSON
- Validate

```
type Foo struct {
  ID string `json:"id"`
  Name string `json:"name"`
}

func (f *Foo) UnmarshalJSON(b []byte) error {
  ...
}

func (f *Foo) MarshalJSON() ([]byte, error) {
  ...
}

func (f *Foo) Validate() validation.ErrorList {
	var result ErrorList
  if f.ID == "" {
    errors = append(errors, validation.NewError("id", "ID is required"))
  }
  return result
}
```

### Migrations

Migrations are written in raw SQL.
Each migration should:
 - have an up and a down
 - be idempotent (do we want a migrations table?)
 - follow the naming convention of <migration order number>_<description of migration>.<up|down>.sql

In `models/migrate.go` we will define a function to execute the migrations. (Could this be generalised in the migrate command?)

## Helpers 

### validation.Validate

`validation.Validate` is used by the controller to validate a struct.
It will check each field in the struct for a `Validate(isNew bool)` method.
If it has the method it will run it on the struct.
This will run recursively across all struct fields.

### sqlmanager.SQLBuilder

The `SQLBuilder` is aimed at decoupling the SQL statement generation from the model.
Each method should return an SQL statement (query and args), that can be executed by the model.

## common.TransactionManager

TODO: Move from `common`

Handles starting, committing and rolling back transactions.