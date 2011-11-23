// Relational algebra results.

package ra

import (
	"strings"
	"constant"
	"table"
	"logg"
	"fmt"
	"st"
)

// Table and the selected rows in the table.
type TableResult struct {
	Table      *table.Table
	RowNumbers []int
}

// Mapping of table name to column name.
type TableColumn struct {
	TableName, ColumnName string
}

// Relational algebras result. For convenience, the thingy is called RA result.
type Result struct {
	Tables  map[string]*TableResult
	Aliases map[string]*TableColumn
}

// Initializes a new Result.
func New() (r *Result) {
	r = new(Result)
	r.Tables = make(map[string]*TableResult)
	r.Aliases = make(map[string]*TableColumn)
	return
}

// Returns a copy of the Result. 
func (r *Result) Copy() *Result {
	aCopy := New()
	for str, tableResult := range r.Tables {
		trCopy := new(TableResult)
		trCopy.RowNumbers = make([]int, len(tableResult.RowNumbers))
		for i, r := range tableResult.RowNumbers {
			trCopy.RowNumbers[i] = r
		}
		trCopy.Table = tableResult.Table
		aCopy.Tables[str] = trCopy

	}
	for str, tableColumn := range r.Aliases {
		tcCopy := new(TableColumn)
		tcCopy.ColumnName = tableColumn.ColumnName
		tcCopy.TableName = tableColumn.TableName
		aCopy.Aliases[str] = tcCopy
	}
	return aCopy
}

// Load all rows of a table into RA result.
func (r *Result) Load(t *table.Table) (*Result, int) {
	_, exists := r.Tables[t.Name]
	if exists {
		return r, st.TableAlreadyExists
	}
	// rowNumbers = list(range(t.NumberOfRows()))
	rowNumbers := make([]int, 0)
	numberOfRows, status := t.NumberOfRows()
	if status != st.OK {
		return r, status
	}
	for i := 0; i < numberOfRows; i++ {
		rowNumbers = append(rowNumbers[:], i)
	}
	r.Tables[t.Name] = &TableResult{t, rowNumbers}
	// Load columns of the table.
	for columnName, _ := range t.Columns {
		if !strings.HasPrefix(columnName, constant.ThePrefix) {
			r.Aliases[columnName] = &TableColumn{t.Name, columnName}
		}
	}
	return r, st.OK
}

// For debugging purpose, prints the RA result.
func (r *Result) Report() {
	var content string
	for name, t := range r.Tables {
		content += "Table: " + name + "\t" + fmt.Sprint(t.RowNumbers) + "\n"
	}
	for alias, c := range r.Aliases {
		content += "Alias " + alias + "\tis " + c.TableName + "." + c.ColumnName + "\n"
	}
	logg.Debug("ra", "Report", content)
}

// Reads a row and return a map representation (name1:value1, name2:value2...)
func (r *Result) Read(rowNumber int) (map[string]string, int) {
	row := make(map[string]string)
	for _, column := range r.Aliases {
		columnName := column.ColumnName
		table := r.Tables[column.TableName]
		tableRow, status := table.Table.Read(table.RowNumbers[rowNumber])
		if status != st.OK {
			break
		}
		row[columnName] = tableRow[columnName]
	}
	return row, st.OK
}

// Returns the number of rows in RA result.
func (r *Result) NumberOfRows() int {
	// All tables have equal number of rows.
	for _, table := range r.Tables {
		return len(table.RowNumbers)
	}
	return 0
}

// Returns TableResult of the specified table name.
func (r *Result) Table(name string) (*TableResult, int) {
	t, exists := r.Tables[name]
	if !exists {
		return nil, st.TableNotFound
	}
	return t, st.OK

}
