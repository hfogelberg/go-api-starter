package main

var routes = Routes{
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
		"/api/create",
		CreateNote,
	},
	Route{
		"Note",
		"GET",
		"/api/notes/{id}",
		GetNoteById,
	},
}
