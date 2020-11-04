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


