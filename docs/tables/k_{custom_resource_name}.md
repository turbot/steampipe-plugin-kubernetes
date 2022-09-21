# Table: {custom_resource_name}

Query data from the custom resource called `{custom_resource_name.group_name}`, e.g., `certificates.cert-manager.io`, `salesforce_custom_app__c`. A table is automatically created to represent each object in the `objects` argument.

## Examples

### Inspect the table structure

List all tables:

```sql
.inspect kubernetes
+---------------------------------+---------------------------------------------------------+
| table                           | description                                             |
+---------------------------------+---------------------------------------------------------+
| salesforce_account_contact_role | Represents the role that a Contact plays on an Account. |
| salesforce_campaign             | Represents Saleforce object Campaign.                   |
| salesforce_case                 | Represents Salesforce object Case.                      |
| salesforce_custom_app__c        | Represents Salesforce object CustomApp__c.              |
+---------------------------------+---------------------------------------------------------+
```

To get details of a specific custom resource table, inspect it by name:

```sql
.inspect "certificates.cert-manager.io"
+---------------------+--------------------------+-------------------------+
| column              | type                     | description             |
+---------------------+--------------------------+-------------------------+
| created_by_id       | text                     | ID of app creator.      |
| created_date        | timestamp with time zone | Created date.           |
| id                  | text                     | App record ID.          |
| is_deleted          | boolean                  | True if app is deleted. |
| last_modified_by_id | text                     | ID of last modifier.    |
| last_modified_date  | timestamp with time zone | Last modified date.     |
| name                | text                     | App name.               |
| owner_id            | text                     | Owner ID.               |
| system_modstamp     | timestamp with time zone | System Modstamp.        |
+---------------------+--------------------------+-------------------------+
```

### Get all values from salesforce_custom_app\_\_c

```sql
select
  *
from
  salesforce_custom_app__c;
```

### List custom apps added in the last 24 hours

```sql
select
  id,
  name,
  owner_id
from
  salesforce_custom_app__c
where
  created_date = now() - interval '24 hrs';
```

### Get details for a custom app by ID

```sql
select
  *
from
  salesforce_custom_app__c
where
  id = '7015j0000019GVgAAM';
```
