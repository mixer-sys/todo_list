# todo_list
TODO-лист с авторизацией

# Run
```bash
sudo mkdir /data
sudo chmod 777 /data
docker build -t todo-list .
docker run -p 8080:8080 -v ./:/data todo-list

```