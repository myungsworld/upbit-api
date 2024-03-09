package models

type Response503Error struct {
	Error struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	} `json:"error"`
}

type Response429Error struct {
	Name string `json:"name"`
}

type Response400Error struct {
	Error struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	} `json:"error"`
}
