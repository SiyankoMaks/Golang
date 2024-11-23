package main 
 
import ( 
 	"encoding/json" 
 	"fmt" 
 	"html/template" 
 	"net/http" 
 	"strings" 
) 
 
type User struct { 
 	Username string json:"username" 
} 
 
var users = make(map[string]User) 
 
// html-страница
var indexHTML =  
<!DOCTYPE html> 
<html lang="en"> 
<head> 
    <meta charset="UTF-8"> 
    <title>User Management</title> 
</head>

<body> 
    <h1>User Management</h1> 
 
    <!-- Кнопка для просмотра всех пользователей --> 
    <button onclick="viewUsers()">View Users</button> 
    <div id="userList"></div> 
 
    <!-- Форма для добавления пользователя --> 
    <h2>Add User</h2> 
    <input type="text" id="username" placeholder="Enter username"> 
    <button onclick="addUser()">Add User</button> 
 
    <!-- Форма для удаления пользователя --> 
    <h2>Delete User</h2> 
    <input type="text" id="deleteUsername" placeholder="Enter username to delete"> 
    <button onclick="deleteUser()">Delete User</button> 
 
    <script> 
        function viewUsers() { 
            fetch('/users') 
                .then(response => response.json()) 
                .then(data => { 
                    let userList = '<h3>Users:</h3><ul>'; 
                    data.forEach(user => { 
                        userList += '<li>' + user.username + '</li>'; 
                    }); 
                    userList += '</ul>'; 
                    document.getElementById('userList').innerHTML = userList; 
                }) 
                .catch(error => console.error('Error:', error)); 
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
            .catch(error => console.error('Error:', error)); 
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
            .catch(error => console.error('Error:', error)); 
        } 
    </script> 
</body> 
</html> 
 
 
func getUsers(w http.ResponseWriter, _ *http.Request) { 
 	w.Header().Set("Content-Type", "application/json") 
 	userList := make([]User, 0, len(users)) 
 	for _, user := range users { 
  		userList = append(userList, user) 
 	} 
 	json.NewEncoder(w).Encode(userList) 
} 
 
func deleteUser(w http.ResponseWriter, r *http.Request) { 
 	username := strings.TrimPrefix(r.URL.Path, "/users/") 
 	if _, exists := users[username]; exists { 
  	delete(users, username) 
  	fmt.Fprintf(w, "User %s deleted", username) 
 	} else { 
  	http.Error(w, "User not found", http.StatusNotFound) 
 	} 
} 
 
func addUser(w http.ResponseWriter, r *http.Request) { 
 	var newUser User 
 	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil { 
  		http.Error(w, err.Error(), http.StatusBadRequest) 
  		return 
 	} 
 	users[newUser.Username] = newUser 
 	fmt.Fprintf(w, "User %s added", newUser.Username) 
} 
 
func main() { 
 	// Обработчик для главной страницы 
 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
  		tmpl := template.New("index") 
  		tmpl, _ = tmpl.Parse(indexHTML) 
  		tmpl.Execute(w, nil) 
 	}) 
 
 	// Обработчик для пути /users 
 	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) { 
  switch r.Method { 
  	case http.MethodGet:
		getUsers(w, r) 
  	case http.MethodPost: 
   		addUser(w, r) 
  	default: 
   		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) 
  } 
 }) 
 
 	// Обработчик для пути /users/{username} 
 	http.HandleFunc("/users/", deleteUser) 
 
 	fmt.Println("Server running at http://localhost:8080") 
 	if err := http.ListenAndServe(":8080", nil); err != nil { 
		panic(err) 
   	} 
}