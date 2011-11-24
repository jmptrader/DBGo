package ra

import (
	"table"
	"st"
	"fmt"
)

// Relational algebra join using nested loops.
func (r *Result) NLJoin(alias string, t2 *table.Table, name string) (*Result, int) {
	t1Column := r.Aliases[alias].ColumnName
	t1 := r.Tables[r.Aliases[alias].TableName]
	t1RowNumbers := make([]int, 0)
	t2RowNumbers := make([]int, 0)
	t2NumberOfRows, status := t2.NumberOfRows()
	if status != st.OK {
		return r, status
	}
	for _, t1RowNumber := range t1.RowNumbers {
		for t2RowNumber := 0; t2RowNumber < t2NumberOfRows; t2RowNumber++ {
			t1Row, status := t1.Table.Read(t1RowNumber)
			if status != st.OK {
				return r, status
			}
			t2Row, status := t2.Read(t2RowNumber)
			if status != st.OK {
				return r, status
			}
			if t1Row["~del"] != "y" && t2Row["~del"] != "y" && t1Row[t1Column] == t2Row[name] {
				t1RowNumbers = append(t1RowNumbers[:], t1RowNumber)
				t2RowNumbers = append(t2RowNumbers[:], t2RowNumber)
			}
		}
	}

	// Correct the order of row numbers of other tables in RA result.
	for _, table := range r.Tables {
		newRowNumbers := make([]int, len(t1RowNumbers))
		for i, keep := range t1RowNumbers {
			newRowNumbers[i] = t1RowNumbers[keep]
		}
		table.RowNumbers = newRowNumbers
		fmt.Println(table.Table.Name, newRowNumbers)
	}
	r.Load(t2)
	t2Table := r.Tables[t2.Name]
	t2Table.RowNumbers = t2RowNumbers
	return r, st.OK
}
