# Todo App

This project is a simple Todo App with a backend running in a Docker container and a PostgreSQL database. Follow the steps below to get the application up and running.

Prerequisites
- Docker installed on your machine
- Docker Hub account (optional for pushing the image)
- .env file configuration for environment variables

## System Design

![image](https://github.com/user-attachments/assets/477cc5fb-cedb-4850-a95e-ad39bc4a1a13)


## Steps to Run the Application

### A) On Local Machine

#### 1. Fork and Clone the Repository

First, clone this repository to your local machine:

```
git clone https://github.com/yourusername/todoapp.git
cd todo-app
```

#### 2. Set up the .env File

Create a .env file in the root of your project and configure the environment variables for your database connection & AWS S3 to store attachments. You can create the file manually or use an example file like so:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_PORT=5432
AWS_S3_BUCKET_NAME=
AWS_REGION=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```

#### 3. Install required dependencies using go mod:

```
go mod tidy
```

GO version used: ```go version go1.22.0 windows/amd64```

#### 4. Run the application:

```
go run main.go
```

The application will start running on http://localhost:8080.


#### 5. Test the Application

Once the application is running, you can test the API endpoints using tools like Postman or cURL. Below are the basic endpoints for the Todo application:

```
- GET /todos - Fetch all todos
- GET /todos/id - Fetch todo by id
- POST /todos - Create a new todo
- PUT /todos/:id - Update an existing todo
- DELETE /todos/:id - Delete a todo
```

Example cURL request to fetch all todos:

```
curl -X GET http://localhost:8080/todos
```

#### 6. Accessing the Database

If you want to connect to the PostgreSQL database, you can do so using a PostgreSQL client like pgAdmin or psql. Use the credentials provided in your .env file.

```
psql -h localhost -U yourusername -d todoapp -p 5432
```


### B) On Docker

#### 1.  Docker Image

You can pull the Docker image from Docker Hub.

```
docker pull mohitkambli8/todoapp:updated
```

Alternatively, if you are building the image locally:

```
docker build -t todoapp .
```

#### 2. Set Up the .env File

Create a .env file in the root of your project and configure the environment variables for your database connection & AWS S3 to store attachments. You can create the file manually or use an example file like so:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_PORT=5432
AWS_S3_BUCKET_NAME=
AWS_REGION=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```

#### 3. Running the Application

If you have already pulled or built the Docker image, you can run the container with the following command:

```
docker run -d -p 8080:8080 --env-file .env mohitkambli8/todoapp:updated
```

- This will start the container in detached mode (-d), and expose port 8080 of the container to port 8080 on your host machine.
- The --env-file .env flag loads the environment variables from your .env file.


#### 4. Test the Application

Once the application is running, you can test the API endpoints using tools like Postman or cURL. Below are the basic endpoints for the Todo application:

```
- GET /todos - Fetch all todos
- GET /todos/id - Fetch todo by id
- POST /todos - Create a new todo
- PUT /todos/:id - Update an existing todo
- DELETE /todos/:id - Delete a todo
```

Example cURL request to fetch all todos:

```
curl -X GET http://localhost:8080/todos
```

#### 5. Accessing the Database

If you want to connect to the PostgreSQL database, you can do so using a PostgreSQL client like pgAdmin or psql. Use the credentials provided in your .env file.

```
psql -h localhost -U yourusername -d todoapp -p 5432
```

#### 6. Stopping the Application

To stop the Docker container, use the following command:

```
docker stop <container_id_or_name>
```

You can find the container ID or name by running:

```
docker ps
```

#### 7. Cleaning Up

If you want to remove the Docker container and the associated image:

1. Remove the container:

```
docker rm <container_id_or_name>
```

2. Remove the image:

```
docker rmi yourusername/todo-app:latest
```


### Testing

Execute the following commands from root directory to test the APIs (Make sure to configure your environment variables in the .env file before testing):

```
  go test -v .\internal\handlers\get_todos_test.go -run TestGetTodos
  go test -v .\internal\handlers\get_todo_by_id_test.go -run TestGetTodoByID
  go test -v .\internal\handlers\create_todo_test.go -run TestCreateTodo
  go test -v .\internal\handlers\update_todo_test.go -run TestUpdateTodo
  go test -v .\internal\handlers\delete_todo_test.go -run TestDeleteTodo
```


### Troubleshooting

- If you face any issues with the database connection, double-check the values in your .env file, especially the host, user, password, and database name.
- Ensure Docker is running properly and check the logs for any error messages.
