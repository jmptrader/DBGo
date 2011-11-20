// Useful constants for the database to run.

package constant

const (
	DefFilePerm               = 0666 // permission for opening .def file of table 
	DataFilePerm              = 0666 // permission for opening .data file of table
	TableDirPerm              = 0755 // permission for creating table directory
	InitFilePerm              = 0666 // permission for opening .init file of database
	MaxColumnNameLength       = 30
	MaxTableNameLength        = 30
	ThePrefix                 = "~" // do not use this prefix to name a database thingy
	MaxTriggerFuncNameLength  = 50
	MaxTriggerParameterLength = 200
	TriggerOperationLength    = 4
	LockTimeout               = 60000000000 // (60 seconds) timeout of table locks in ns
	ExclusiveLockFilePerm     = 0666        // permission for opening .exclusive file of table lock
)

// Returns the extensions names which table files have.
func TableFiles() []string {
	return []string{".data", ".def"}
}

// Returns the column names and lengths which a new table must have. 
func DatabaseColumns() map[string]int {
	return map[string]int{ThePrefix + "del": 1}
}

// Returns the directory suffixes which table directories have. 
func TableDirs() []string {
	return []string{".shared"}
}

// Returns the column names and lengths of a database trigger lookup table.
func TriggerLookupTable() map[string]int {
	return map[string]int{"TABLE": MaxTableNameLength, "COLUMN": MaxColumnNameLength,
		"FUNC": MaxTriggerFuncNameLength, "OP": TriggerOperationLength,
		"PARAM": MaxTriggerParameterLength}
}
