package repository

import (
	"strconv"
	"testTask/internal/model"
	"time"
)

func (r *Repository) AddRate(data *model.Rates) {
	a_date, err := time.Parse("02.01.2006", data.Date)
	if err != nil {
		r.Log.ErrLog("[ ERR ] convert date" + err.Error())
		return
	}
	for i, item := range data.Item {

		value, err := strconv.ParseFloat(item.Description, 32)
		if err != nil {
			r.Log.ErrLog("[ ERR ] convert to float" + err.Error())
			return
		} else {
			println("[ INF ]", "write item", i, "to database")
			_, err := r.Db.Exec("insert into R_CURRENCY(title, code, value, a_date) values (?, ?, ?, ?)", item.Fullname, item.Title, value, a_date)
			if err != nil {
				r.Log.ErrLog("[ ERR ] insert in database :" + " data item.Title, value: " + item.Title + " " + err.Error())
				return
			}
		}
	}
}

func (r *Repository) GetRate(date time.Time, code string) (*[]model.RatesOut, error) {

	query := "select title, code, value, a_date from R_CURRENCY where a_date = ?"
	params := []interface{}{date}
	if len(code) != 0 {
		query += " and code = ?"
		params = append(params, code)
	}
	rows, err := r.Db.Query(query, params...)
	if err != nil {
		r.Log.ErrLog("[ ERR ] data query" + err.Error())
		return nil, err
	}

	defer rows.Close()
	var RatesOut []model.RatesOut

	for rows.Next() {
		var Rates model.RatesOut
		rows.Scan(&Rates.Title, &Rates.Code, &Rates.Value, &Rates.Date)
		if err != nil {
			r.Log.ErrLog("[ ERR ] read result from DS" + err.Error())
			return nil, err
		}
		RatesOut = append(RatesOut, Rates)
	}

	return &RatesOut, nil
}
