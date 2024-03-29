package keyvalue

import (
	"database/sql"
	"main/interpreter"
	"main/interpreter/interop"

	_ "github.com/mattn/go-sqlite3"
)

// Definition of open function that returns a database struct
var OpenDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{
			Type: interpreter.TypeString,
		},
	},
	ReturnType: interpreter.StructDef{
		GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeStruct},
		Properties: map[string]int{
			"get":    0,
			"set":    1,
			"delete": 2,
			"close":  3,
		},
		PropertyDefs: []interpreter.TypeDef{
			GetDef,
			SetDef,
			DeleteDef,
			CloseDef,
		},
		Name: "KeyValueDb",
	},
}

type KeyValueDb struct {
	db *sql.DB
}

// Opens the database
func Open(file string) []any {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS key_value (
	  key VARCHAR(255) NOT NULL PRIMARY KEY,
	  value VARCHAR(65535) NOT NULL
	);`); err != nil {
		panic(err)
	}

	return interop.CreateRuntimeStruct(&KeyValueDb{db}, []string{"Get", "Set", "Delete", "Close"})
}

var GetDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{Type: interpreter.TypeString},
	},
	ReturnType: interpreter.GenericTypeDef{Type: interpreter.TypeString},
}

// Gets a value by it's key
func (kv *KeyValueDb) Get(key string) string {
	row := kv.db.QueryRow("SELECT value FROM key_value WHERE key = ?;", key)
	var result string
	row.Scan(&result)
	return result
}

var SetDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{Type: interpreter.TypeString},
		interpreter.GenericTypeDef{Type: interpreter.TypeString},
	},
}

// Sets a value by it's key
func (kv *KeyValueDb) Set(key string, value string) {
	kv.db.Exec("INSERT INTO key_value VALUES (?, ?) ON CONFLICT DO UPDATE SET value = ?;", key, value, value)
}

var DeleteDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{Type: interpreter.TypeString},
	},
}

// Deletes a value by it's key
func (kv *KeyValueDb) Delete(key string) {
	kv.db.Exec("DELETE FROM key_value WHERE key = ?", key)
}

var CloseDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
}

// Closes open database
func (kv *KeyValueDb) Close() {
	kv.db.Close()
}
