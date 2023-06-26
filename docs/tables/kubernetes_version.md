# Table: kubernetes_version

The `kubernetes_version` table can be used to query client and server version information for the current context.

> Note: Should only return a single row of data.

## Examples

### Get version information

```sql
select
  *
from
  kubernetes_version;
```
