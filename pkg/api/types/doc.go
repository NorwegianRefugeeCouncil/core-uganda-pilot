/**

The Form feature allows user to design forms for data collection.
The fields that can be attached to the forms are of different types.

It is also possible to add "Links" between records, such that a
record in a form contains a link to a record in another form.
These are "Reference" fields.

For example

Countries
==========
NAME   CODE
Uganda UG
Kenya  KY

Projects
========
NAME  COUNTRY
PR1   UG (references Countries[Uganda])
PR2   KY (references Countries[Kenya])

Certain fields can also be marked as Key fields. When this happens,
there will be no two records with the same values for the key fields.
The combination of values in the Key fields are unique within all the
records of a form.

For example
YEAR (key)    MONTH (key)    COMMENT
2021          Jan            It's cold
2021          Feb            Still cold  <- allowed because 2021/Feb != 2021/Jan
2021          Jan            Colder      <- !!not allowed, there already exist a field with 2021/Jan.

There exists a special case when a form has 1 key field which happens to be also a reference field.
This is in many regards an "Inheritance Link"

For Example

Individuals
===========
ID   Name
1    John
2    Sandy

Staff
======================
Individual   StaffCode
1            Staff-1
2            Staff-2

Here, the "Individual" field is both a Key and a Reference field to the Individual form.
That means that there can be no 2 records in the Staff form for the same Individual.
In that sense, the "Staff" form "Extends" the "Individual" form.

We can take advantage of this special scenario to implement the business needs of NRC.

The main complexity for case management revolves around the required flexibility of operating
in different geographic areas, where there are different data collection needs. At the same time,
there are global requirements that all data collection operations must abide by.

For example, NRC Global requires that we ask ALL beneficiaries the same set of questions.
At the same time, certain countries have very local-specific questions that they need to ask.

We can solve the problem by having multiple profiles that "Inherit" from each other using the
special "Inheritance" link.

For example, the "Individual" form might represent any person that NRC enters in contact with.

Then, the "Individual Beneficiary" form might represent a person that is also a Beneficiary. Since we need
to ask certain questions to all beneficiaries regardless of their country, we can add the fields in this
form that correspond to those questions we need to ask. We can also configure this form so that it
"Extends" the "Individual" form. This could be summarized as "A Beneficiary is an Individual, but an
Individual is not necessarily a Beneficiary"

Furthermore, every country/area that collects information about beneficiaries could create their own
area-specific form "Kenya Beneficiary", that extends the "Global Beneficiary". In this form,
they would add the area-specific questions. This could be summarized as "A Kenya Beneficiary is
always a Global Beneficiary, but a Global Beneficiary is not always a Kenya Beneficiary".

This gives a hierarchical structure like so

- Household
    |- Kenya Household
    |- Uganda Household
    |- Colombia Household

- Temporary Tented Settlement
    |- Kenya TTS
    |- Uganda TTS
    |- Colombia TTS

- Individual
  |- Global Beneficiary
  |   |- Kenya Beneficiary
  |   |- Uganda Beneficiary
  |   |- Colombia Beneficiary
  |- Other forms about individuals...

It would be nice to not pollute the API of the Forms in order to enable the Case Management on top of it.
The ideal scenario would be to add a "Case" object which "surrounds" the Form with additional logic.

Proposal:

This new database table would store which recipients can be attached to which cases.
When the user would create a new Case Record, the information in this table would allow
us to restrict what kind of entity the Case can be attached to for the recipient.
Eg.

Uganda ICLA Case can only be attached to Uganda Individual
Kenya Shelter Case can only be attached to Kenya Household
etc.

CMS_ALLOWED_RECIPIENTS: Stores which recipients can be attached to cases
=============================================================================
CASE_ID  CASE_DATABASE_ID  RECIPIENT_FORM_ID  RECIPIENT_FORM_DATABASE_ID
=============================================================================
1        my-database       case-type-1        recipient-form-1   my-database
1        my-database       case-type-1        recipient-form-1   my-database


Proposal:
This new database table would store which Forms are "Recipients".
When the user would create a new Case Type, the user would need to specify
what kind of Recipients this case is allowed to be attached to. (To populate
the above table). This table would basically restrict which forms the user
can select.

Otherwise, the user could Select some random form that is not at all a
Recipient form, such as "Project Codes". Here, if we add "Uganda Individual",
"Kenya Individual", "Kenya Household", the case type creator could only
be presented with these 3 options.

# CMS_RECIPIENT_FORMS: restricts which forms can be used as recipients
# ========================
# FORM_ID     DATABASE_ID   IS_ABS.
# ========================
# kenya-b     my-database   false
# uganda-b    my-database   false
# global-b    my-database   true

Proposal:
We need to allow referrals between cases. This new database table would store
the referrals.

CMS_REFERRALS: stores case referrals
============================================================================
CASE_ID  REFERRED_FROM_CASE_ID  REFERRED_FROM_CASE_DATABASE_ID  REFERRED_BY
============================================================================
1        my-case-form           my-database                     johndoe@email.com
2        my-case-form           my-database                     stacey@email.com

sample payload for creating a new type of case
POST cms.nrc.no/v1/casedefinition
{
	"name": "my-case-definition",
    "databaseId": "abc123",
    "fields": [{"name":"my-field","fieldType":{"text":{}}}],
	"recipients": [
		{"formRef": {"formId":"my-recipient-form","databaseId":"abc123"}},
		{"form":    {"name":"anonymous", "fields":[{"name":"my-field"}]}} // potential in the future?
	]
}


sample payload for creating a new record of a case
POST core.nrc.no/v1/record
{
  "formId": "my-case-definition-id",
  "databaseId": "abc123",
  "values": {
    "my-field-1": "my-value",
	"my-other-field": "my-other-value",
    "status": "open"
  }
}

*/

package types
