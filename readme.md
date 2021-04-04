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
./krohobor --database=postgres-source db list
./krohobor --database=postgres-target db list
```

Info of database
```env
./krohobor --database=postgres-source --dbname=test1 db read
./krohobor --database=postgres-target --dbname=test1 db read
```

Create a database
```env
./krohobor --database=postgres-source --dbname=test_new db create
./krohobor --database=postgres-target --dbname=test_new db create
```

Delete a database
```env
./krohobor --database=postgres-source --dbname=test_new db delete
./krohobor --database=postgres-target --dbname=test_new db delete
```

Create dump of all databases
```env
./krohobor --database=postgres-source dump create
./krohobor --database=postgres-source --storage=local dump create
```

Create dump of concrete databases
```env
./krohobor --database=postgres-source --dbname=test1 dump create
./krohobor --database=postgres-source --storage=local --dbname=test1 dump create
```

List of dumps
```env
./krohobor dump list
./krohobor --storage=local dump list
```

Delete dump
```env
./krohobor --name=test1.sql dump delete
```

Restore dump
```env
./krohobor --name=all.sql --database=postgres-target dump restore
```

Restore dump to database
```env
./krohobor --name=test1.sql --database=postgres-target --dbname=new dump restore
```
