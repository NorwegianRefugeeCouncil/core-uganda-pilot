# Table of Contents

- [database](#database)
- [folder](#folder)
- [form](#form)
- [record](#record)
- [fields](#field)
	- [text field](#text-field)
	- [multiline text field](#multiline-text-field)
	- [single choice field](#single-choice-field)
	- [multiple choice field](#multiple-choice-field)
	- [date field](#date-field)
	- [month field](#month-field)
	- [week field](#week-field)
	- [reference field](#reference-field)
	- [sub form field](#sub-form-field)
- [field properties](#field-properties)
	- [required fields](#required-fields)
	- [key fields](#key)

# Concepts

#### Database

Databases are containers for data. Access can be granted on a per-database basis.

#### Folder

Folders are logical groups for forms. Access can be granted on a per-folder basis.

#### Form

Forms define the structure of the data captured.

#### Record

Record represent an entry in a form. Whenever a user fills a form and submits it, a new record is created.

## Fields

[forms](#form) are composed of multiple [fields](#fields). When designing a [form](#form), the user can add as
many [fields](#fields) as needed to collect the appropriate data.

### Text Field

Records text as single-line input

### Multiline Text Field

Records multi-line text

### Single Choice Field

Allows the user to choose between a predefined set of choices

### Multiple Choice Field

Allows the user to choose multiple options within a predefined set of choices

### Date Field

Allows the user to select a date `2021-12-31`

### Month Field

Allows the user to select a month `2021-12`

### Week Field

Allows the user to select a week `2021-W52`

### Reference Field

Allows the user to select a record from another [forms](#form)

### Sub Form Field

**Sub Forms** are regular [forms](#form), but they are attached to a "parent" [record](#record). For example, one might
define a [form](#form) called `Case`, and add a `Sub Form` field `Follow-Ups`. The user could add `Follow-Ups` to a
single `Case` [record](#record).

## Field Properties

### Required Fields

Fields marked as `required` will require the user to enter a value when creating a record. If the user submits a record
without that field filled in, the record will be denied.

### Key

When one or multiple [fields](#fields) are marked as `key` fields, that means that the value of those
`key` fields **MUST** be unique within all records of that form. Also, the record values for that field
**MUST** be supplied by the user (the `key` fields are automatically `required` as well)

#### Notes

- `Reference` fields **can** be marked as `key` fields.
- `Sub Form` fields **cannot** be marked as `key` or `required` fields
- `Multiple Selection` fields **cannot** be marked as `key` fields
- `Header` fields **cannot** be marked as `key` fields
- `Key` fields are always `required`

#### Example

For example, a user might define a form called `Monthly Reports`. The user would add a
`Reporting Month` field (of type `Month`).

| Month (key)  | Amount
---------------|-------
| 2020-01      | 100$

If the `Reporting Month` field is marked as a `key` field, then the users will not be able to add a record twice with
the same month.

| Month (key)  | Amount |   |
---------------|--------|---|
| 2020-01      | 100$   |
| 2020-01      | 120$   | **This record would not be allowed because January appears twice!**
| 2020-02      | 150$   | **This record would be allowed!**

The `key` property can be enabled on multiple fields in a form. In such case, any combination of the fields that are
marked as `key` cannot appear twice.

| Month (key)  | Project (key) | Amount |   |
---------------|---------------|--------|---|
| 2020-01      | Project A     | 100$   |
| 2020-01      | Project A     | 120$   | **This record would not be allowed**
| 2020-01      | Project B     | 110$   | **This record would be allowed**
| 2020-02      | Project A     | 180$   | **This record would be allowed**


