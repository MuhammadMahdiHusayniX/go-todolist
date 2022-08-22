# go-todolist
TodoList server with gmail authentication

Run:
docker compose up

Open 127.0.0.1:3000 in browser

List of API route below (Using bearer token)

GET ALL TODO => METHOD: GET, ROUTE: http://127.0.0.1:3000/todo

CREATE TODO => METHOD: POST, ROUTE: http://127.0.0.1:3000/todo, JSON BODY: { "Task": "Optimize code", "CreatedBy": "Mahdi" }

MARK TODO AS COMPLETE => METHOD: GET, ROUTE: http://127.0.0.1:3000/todo/complete/:id

DELETE TODO => METHOD: DELETE, ROUTE: http://127.0.0.1:3000/todo/:id