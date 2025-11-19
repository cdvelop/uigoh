package formPlatform

type entity struct {
	Name string //ej: users.user, module.product
	// Legend        string // e.g.: Person, User, Product
	IsTable   bool
	TableName string //table name db ej: user, product
	// ParentStruct any
	Fields []field
	// StructHandler *structHandler

	HtmlForm string //html form

	lastHtmlID int // counter for multiple id values
}
