package main

import "fmt"

func (a *App) getDatasetByID(id int64) (dataset Dataset, err error) {
	rows, err := a.DB.Query("select field_id, name, description, types from datacatalog.public.fields where dataset_id =$1;", id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var field Field
		err = rows.Scan(&field.FieldID, &field.Name, &field.Doc, &field.Type)
		if err != nil {
			return
		}

		dataset.Fields = append(dataset.Fields, field)
	}

	err = rows.Err()
	if err != nil {
		return
	}
	dataset.DatasetID = id
	fmt.Println(dataset)

	return

}
