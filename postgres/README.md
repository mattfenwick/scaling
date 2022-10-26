# PG

## queries

insert:

```sql
insert into documents 
    (parsed, parse_error)
values 
    (format('{"abc": [123,"def"]}}')::json, 'def');
```

or:

```sql
insert into documents 
    (parsed, parse_error)
values 
    ('{ "customer": "John Doe", "items": {"product": "Beer","qty": 6}}', 'def');

select parsed -> 'items' -> 'product' as "hi" from documents ;
```