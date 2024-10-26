## DATABASE

This directory is used for everything with database

### SETUP DATABASE

```
# Run mysql container
$ docker run --name kiber_sql -e MYSQL_ROOT_PASSWORD=pass -d mysql:latest
# Execute sql script
# (need to wait for container to be ready before executing this command)
# (you can modify schema.sql script for it to be )
$ docker exec -i kiber_sql mysql -u root -ppass -h localhost -P 3306 < schema.sql
```

### ACCESS DATABASE

```
$ docker exec -it kiber_sql mysql -u root -ppass -h localhost -P 3306
```
