package ctrl

import (
	"encoding/json"
	"fmt"

	"github.com/Hy-Iam-Noval/dacol-2/src/DB"
	"github.com/Hy-Iam-Noval/dacol-2/src/helpers"
	"github.com/Hy-Iam-Noval/dacol-2/src/validation"
)

// first map for the index
// inside contain error store in map where mapping by name form
type Error = map[int]map[string]string

func FindProd(c Ctx) error {
	param := c.Params("id")

	data := new(DB.Product)
	DB.
		Select("product", "*", "id = "+param).
		Single().
		Scan(&data.Id, &data.Id_user, &data.Name, &data.Qty, &data.Price, &data.Img, &data.Descript)
	if (*data == DB.Product{}) {
		return c.Send(nil)
	}

	return SendDatas(c, data)
}

// Get /auth/product
func AllProd(c Ctx) error {
	user := helpers.ParseTokenUser(c.Get(UserKey))
	rows := DB.
		Select("product", "*", fmt.Sprintf("id_user = %d", user.Id)).
		Many()

	datas := []DB.Product{}
	for rows.Next() {
		var data DB.Product
		rows.Scan(&data.Id, &data.Id_user, &data.Name, &data.Qty, &data.Price, &data.Img, &data.Descript)
		if (data != DB.Product{}) {
			datas = append(datas, data)
		}
	}

	return SendDatas(c, datas)
}

func validateProduct(req []DB.Product) map[int]map[string]string {
	var errList map[int]map[string]string

	for _, i := range req {
		/// validate
		err := validation.IsValid(i)
		if err != nil && errList == nil {
			errList = make(map[int]map[string]string)
		}

		if err != nil {
			errList[int(i.Id)] = err
		}

	}

	return errList
}

func AddProd(c Ctx) error {
	req := []DB.Product{}
	user := helpers.ParseTokenUser(c.Get(UserKey))

	// get request
	IsError(c.BodyParser(&req))

	if len(req) == 0 {
		return End(c)
	}

	// validate
	if err := validateProduct(req); err != nil {
		return c.Status(helpers.Invalid).JSON(ErrRes(err))
	}

	// inserting data
	dataInsert := ""
	for _, i := range req {
		jsonImg := DB.JSONImg{}
		json.Unmarshal([]byte(i.Img), &jsonImg)

		data := fmt.Sprintf(
			`(%d,"%s",%d,%d,"%s","%s")`,
			user.Id, i.Name, i.Qty, i.Price, jsonImg.Name, i.Descript)
		if dataInsert == "" {
			dataInsert = data
		} else {
			dataInsert = fmt.Sprintf(`%s,%s`, dataInsert, data)
		}
	}

	res := DB.
		Query(
			fmt.Sprintf(`
				INSERT INTO product(id_user, name, qty, price, img, descript)
				VALUES %s
			`, dataInsert),
		).
		Exec()
	result, _ := res.RowsAffected()
	if result > 0 {
		return End(c)
	}

	return End(c)
}

func DeleteProd(c Ctx) error {
	param := c.Params("id")

	// delete product by id
	DB.
		Delete("product", "id = "+param).
		Exec()

	return End(c)
}

func UpdateProd(c Ctx) error {
	req := []DB.Product{}

	IsError(c.BodyParser(&req))
	if len(req) == 0 {
		return End(c)
	}

	if err := validateProduct(req); err != nil {
		return c.Status(helpers.Invalid).JSON(err)
	}

	// update
	for _, i := range req {
		DB.Query(
			fmt.Sprintf(`
				UPDATE product 
				SET name = "%s", qty = %d, price = %d, img = %s, descript = "%s"
				WHERE id = %d
			`, i.Name, i.Qty, i.Price, i.Img, i.Descript, i.Id,
			)).Exec()
	}

	return End(c)

}
