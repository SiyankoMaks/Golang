package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

type User struct {
	Username string `json:"username"`
}

var (
	users    = make(map[string]User)
	port     = ":8080"
	indexHTML = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Управление Пользователями</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #000000;
            color: #DFF6FF;
            margin: 0;
            padding: 20px;
        }
        h1, h2 {
            color: #47B5FF;
        }
        button {
            background-color: #47B5FF;
            color: #000000;
            border: none;
            padding: 10px 20px;
            cursor: pointer;
            margin: 5px 0;
        }
        button:hover {
            background-color: #256D85;
        }
        input {
            padding: 10px;
            margin: 5px 0;
            border: 1px solid #47B5FF;
            border-radius: 5px;
        }
        #userList {
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <h1>Управление Пользователями</h1>

    <button onclick="viewUsers()">Просмотреть пользователей</button>
    <div id="userList"></div>

    <h2>Добавить пользователя</h2>
    <input type="text" id="username" placeholder="Введите имя пользователя">
    <button onclick="addUser()">Добавить</button>

    <h2>Удалить пользователя</h2>
    <input type="text" id="deleteUsername" placeholder="Введите имя пользователя для удаления">
    <button onclick="deleteUser()">Удалить</button>

    <script>
        function viewUsers() {
            fetch('/users')
                .then(response => response.json())
                .then(data => {
                    let userList = '<h3>Список пользователей:</h3><ul>';
                    data.forEach(user => {
                        userList += '<li>' + user.username + '</li>';
                    });
                    userList += '</ul>';
                    document.getElementById('userList').innerHTML = userList;
                })
                .catch(error => console.error('Ошибка:', error));
        }

        function addUser() {
            const username = document.getElementById('username').value;
            fetch('/users', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: username })
            })
            .then(response => response.text())
            .then(data => {
                alert(data);
                viewUsers();
            })
            .catch(error => console.error('Ошибка:', error));
        }

        function deleteUser() {
            const username = document.getElementById('deleteUsername').value;
            fetch('/users/' + username, {
                method: 'DELETE'
            })
            .then(response => response.text())
            .then(data => {
                alert(data);
                viewUsers();
            })
            .catch(error => console.error('Ошибка:', error));
        }
    </script>
</body>
</html>
`
)

func main() {
	r := chi.NewRouter()

	// Главная страница
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").Parse(indexHTML))
		tmpl.Execute(w, nil)
	})

	// Маршрут для пользователей
	r.Route("/users", func(r chi.Router) {
		r.Get("/", getUsers)          // Просмотр всех пользователей
		r.Post("/", addUser)          // Добавление пользователя
		r.Delete("/{username}", deleteUser) // Удаление пользователя
	})

	// Запуск сервера
	serverAddress := "http://localhost" + port
	println("Сервер запущен по адресу", serverAddress)
	http.ListenAndServe(port, r)
}

// Просмотр всех пользователей
func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}

// Удаление пользователя
func deleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if _, exists := users[username]; exists {
		delete(users, username)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователь успешно удалён"))
	} else {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
	}
}

// Добавление пользователя
func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := users[newUser.Username]; exists {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}
	users[newUser.Username] = newUser
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно добавлен"))
}
