# todo_list
TODO-лист с авторизацией

# Run
```bash
sudo mkdir /data
sudo chmod 777 /data
docker build -t todo-list .
docker run -p 8080:8080 -v ./:/data todo-list

```


# Help

```bash
docker exec -it todo_list_db_1 psql -U user -d db
```