# mysql-waiter

This is supposed to act as a dependency shim for mysql containers.
The need for this shim arises when you want to have a mysql container that will be initialized with data,
if you use some kind of wait-for-it script to listen for the mysql port or docker healthchecks (at a point in time where they are actually honored for depends_on), you will soon notice that the mysql service is available even though the import might not be done.

This is where you will use this container, it takes environment variables to configure what to wait for.
```
environment:
  - MYSQL_HOST=mysql-container-name
  - MYSQL_USER=hopefully-not-root
  - MYSQL_PASSWORD=change_me
  - MYSQL_DB=database-that-is-imported
  - MYSQL_TABLE=a=table-within-that-db
```

It will retry connecting to the server, selecting the database and looking for that table, if successfull it will listen on port 3000.
That way you can use this shim with any of your healthchecks or wait-for-it scripts and just look for port 3000 to know when your mysql container is ready

