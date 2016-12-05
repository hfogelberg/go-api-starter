package main

var routes = Routes{
	Route{
		"Signup",
		"POST",
		"/api/signup",
		Signup,
	},
	Route{
		"Login",
		"POST",
		"/api/login",
		Login,
	},
	Route{
		"Notes",
		"GET",
		"/api/notes",
		GetNotes,
	},
	Route{
		"Greeting",
		"GET",
		"/api/greeting",
		GetGreeting,
	},
	Route{
		"Create",
		"POST",
		"/api/notes",
		CreateNote,
	},
	Route{
		"Note",
		"GET",
		"/api/notes/{id}",
		GetNoteById,
	},
}
