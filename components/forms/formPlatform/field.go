package formPlatform

import "reflect"

type sourceData interface {
	DataSource() any
}

type dataSource struct {
	data sourceData
}

type className string

type fieldset struct {
	cssClasses []className
}

type field struct {
	Index  uint32 // index of the field
	Legend string // e.g.: ID, Name, Phone

	DbType     dbFieldType // e.g.: INT, VARCHAR(255)
	Unique     bool        // unique and unalterable field in db
	NotNull    bool
	PrimaryKey bool    // primary key of the table
	ForeignKey *entity // foreign key of the table

	Parent             *entity              // pointer to the entity parent
	reflectType        reflect.Type         // type of the structure
	reflectStructField *reflect.StructField // type of the structure

	// Input configuration
	// HtmlID           string  // unique ID for HTML rendering
	dataSource
	fieldset

	Id                 string              `ctx:"ui"`
	Name               string              `ctx:"ui,db"` //eg address,phone
	Type               string              `ctx:"ui"`    // input type  eg : text, password, email
	htmlName           string              `ctx:"ui"`    //eg input,select,textarea
	customName         string              `ctx:"ui"`    //eg only_text, only_number...
	allowSkipCompleted bool                `ctx:"ui"`    //permite que el campo no sea completado
	hideTyping         bool                `ctx:"ui"`    //oculta el valor mientras se escribe
	PlaceHolder        string              `ctx:"ui"`
	Title              string              `ctx:"ui"` //info
	Min                string              `ctx:"ui"` //valor mínimo
	Max                string              `ctx:"ui"` //valor máximo
	Maxlength          string              `ctx:"ui"` //ej: maxlength="12"
	Autocomplete       string              `ctx:"ui"`
	Rows               string              `ctx:"ui"` //filas ej 4,5,6
	Cols               string              `ctx:"ui"` //columnas ej 50,80
	Step               string              `ctx:"ui"`
	Oninput            string              `ctx:"ui"` // ej: "miRealtimeFunction()" = oninput="miRealtimeFunction()"
	Onkeyup            string              `ctx:"ui"` // ej: "miNormalFuncion()" = onkeyup="miNormalFuncion()"
	Onchange           string              `ctx:"ui"` // ej: "miNormalFuncion()" = onchange="myFunction()"
	Accept             string              `ctx:"ui"`
	Multiple           string              `ctx:"ui"` // multiple
	Value              string              `ctx:"ui"`
	Class              []className         `ctx:"ui"` // clase css ej: class="age"
	DataSet            []map[string]string `ctx:"ui"` // dataset ej: data-id="123" = map[string]string{"id": "123"}
	options            []map[string]string `ctx:"ui"` // ej: [{"m": "male"}, { "f": "female"}]

	permitted
}
