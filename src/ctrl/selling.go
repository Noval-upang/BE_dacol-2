package ctrl

import (
	"encoding/json"
	"fmt"

	"github.com/Hy-Iam-Noval/dacol-2/src/DB"
	"github.com/Hy-Iam-Noval/dacol-2/src/helpers"
)

func AddSelling(c Ctx) error {
	req := new(DB.Selling)

	// get req
	IsError(c.BodyParser(req))

	// add selling
	DB.Create("selling", *req).Exec()

	// update
	sellingTarget := new([]struct {
		Id  int `json:"id"`
		Qty int `json:"qty"`
	})
	json.Unmarshal([]byte(req.Product_list), sellingTarget)

	for _, i := range *sellingTarget {
		DB.Query(
			fmt.
				Sprintf(`
					UPDATE product  
						SET qty = qty - %d 
						WHERE id = %d`,
					i.Qty, i.Id),
		).Exec()
	}

	return End(c)
}

func AllByIDSelling(c Ctx) error {
	param := c.Params("id")

	datas := []DB.Selling{}
	rows := DB.
		Select("selling", "*", fmt.Sprintf(`product_list LIKE '%%"id":%s%%'`, param)).
		Many()

	for rows.Next() {
		data := DB.Selling{}

		if err := rows.Scan(&data.Id, &data.Product_list, &data.Provit, &data.Date); err != nil {
			return FiberBadReq(err)
		}
		if (data != DB.Selling{}) {
			datas = append(datas, data)
		}
	}

	return SendDatas(c, datas)
}

func AllSelling(c Ctx) error {
	user := helpers.ParseTokenUser(c.Get(helpers.UserKey))
	datas := []DB.Selling{}

	for rows := DB.
		Query(
			fmt.Sprintf(`
				SELECT  
					selling.id, selling.product_list,selling.provit, selling.date
				FROM selling,  
					json_tree(product_list) 
				WHERE EXISTS 
					(SELECT * FROM product WHERE id = json_tree.value AND id_user = %d)
			`, user.Id),
		).
		Many(); rows.Next(); {
		data := DB.Selling{}
		err := rows.Scan(&data.Id, &data.Product_list, &data.Provit, &data.Date)
		IsError(err)
		if (data != DB.Selling{}) {
			datas = append(datas, data)
		}
	}

	return SendDatas(c, datas)
}
