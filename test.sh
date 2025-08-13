curl -i -k -X POST http://localhost:8080/users/signup \
-H "Content-Type: application/json" \
-d '{"username": "testuser", "password": "testpassword"}'

token=$(curl -X POST http://localhost:8080/users/login \
-H "Content-Type: application/json" \
-d '{"username": "testuser", "password": "testpassword"}' | jq -r .token)

curl -i -k -X GET http://localhost:8080/users \
-H "Authorization: Bearer $token"

curl -i -k -X PUT http://localhost:8080/users \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $token" \
-d '{"username": "newusername", "password": "newpassword"}'

curl -i -k -X GET http://localhost:8080/users \
-H "Authorization: Bearer $token"

curl -i -k -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $token" \
-d '{"title": "New Task", "description": "Task description"}'

curl -i -k -X GET http://localhost:8080/tasks \
-H "Authorization: Bearer $token"

curl -i -k -X GET http://localhost:8080/tasks/5 \
-H "Authorization: Bearer $token"

curl -i -k -X PUT http://localhost:8080/tasks/5 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $token" \
-d '{"title": "Updated Task", "description": "Updated task description"}'

curl -i -k -X DELETE http://localhost:8080/tasks/5 \
-H "Authorization: Bearer $token"
