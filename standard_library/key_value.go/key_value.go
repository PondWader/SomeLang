package keyvalue

import (
	"database/sql"
	"fmt"
	"main/interop"
	"main/interpreter"

	_ "github.com/mattn/go-sqlite3"
)

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
			"get": 0,
			"set": 1,
		},
		PropertyDefs: []interpreter.TypeDef{
			GetDef,
			SetDef,
		},
	},
}

type KeyValueDb struct {
	db *sql.DB
}

const create string = `
CREATE TABLE IF NOT EXISTS key_value (
  key VARCHAR(255) NOT NULL PRIMARY KEY,
  value VARCHAR(65536) NOT NULL
);`

func Open(file string) []any {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(create); err != nil {
		panic(err)
	}

	fmt.Println(db, err)
	return interop.CreateRuntimeStruct(&KeyValueDb{db})
}

var GetDef = interpreter.FuncDef{}

func (kv *KeyValueDb) Get(self []any, key string) {
	result, err := self[0].(*sql.DB).Query("SELECT * FROM key_value WHERE key = ?", key)
	fmt.Println(result, err)
}

var SetDef = interpreter.FuncDef{}

func (kv *KeyValueDb) Set(self []any, key string, value string) {

}
