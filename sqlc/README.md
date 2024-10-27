## DATABASE

This directory is used for everything with database

### SETUP DATABASE

```
# Run mysql container
$ docker run --network host --name kiber_sql -e MYSQL_ROOT_PASSWORD=pass -d mysql:8.0
# Execute sql script
# (need to wait for container to be ready before executing this command)
# (you can modify schema.sql script for it to be )
$ docker exec -i kiber_sql mysql -u root -ppass -h localhost -P 3306 < schema.sql
```

### ACCESS DATABASE

```
$ docker exec -it kiber_sql mysql -u root -ppass -h localhost -P 3306
```

### EXCEPTIONS

Currently mysql container uses 3306 port (the same as host mysql port) so might need to disable mysql on host machine.

### INSERTING INTO `card` TABLE

`uid` column should be without any spaces and in `user_name` column can be anything except empty string for program to work. Eg. " " is valid, "" is not valid.
