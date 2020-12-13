Create user with grants
```sql
CREATE USER appuser WITH LOGIN PASSWORD '123456';
\c auth
GRANT ALL ON DATABASE auth TO appuser;
GRANT ALL ON SCHEMA auth TO appuser;
GRANT ALL ON ALL TABLES IN SCHEMA public TO appuser;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO appuser;
```

Show grants for table in database
```sql
\dp
```

List of databases
```env
./krohobor db list
```

Info of database
```env
./krohobor --db=position db read
```


Create dump of all databases
```env
./krohobor dump create
```

Create dump of concrete databases
```env
./krohobor --db=position dump create
```

List of dumps
```env
./krohobor dump list
```

Delete dump
```env
./krohobor --name=all.sql dump delete
```

Restore dump
```env
./krohobor --name=all.sql --target=localhost:5431 dump restore
```
