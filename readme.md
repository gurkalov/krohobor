Create user with grants
```sql
CREATE USER mytest WITH LOGIN PASSWORD '123456';
\c auth
GRANT ALL ON ALL TABLES IN SCHEMA public TO mytest;
```

Show grants for table in database
```sql
\dp
```
