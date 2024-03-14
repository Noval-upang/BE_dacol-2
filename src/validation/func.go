package validation

import (
	"fmt"

	"github.com/Hy-Iam-Noval/dacol-2/src/DB"
)

func unique(f Field) bool {
	user := DB.User{}
	DB.
		Select("user", "*", fmt.Sprintf(`email = "%s"`, f.Field().String())).
		Single().
		Scan(&user.Id, &user.Email, &user.Password)

	return user == DB.User{}
}

func uniqueProduct(f Field) bool {
	val, ok := f.Parent().Interface().(DB.Product)
	if !ok {
		panic(fmt.Sprintf("Type %T", val))
	}

	prod := new(DB.Product)

	DB.
		Select("product", "*", "id_suser = "+string(rune(val.Id)), fmt.Sprintf(`name = "%s"`, f.Field().String())).
		Single().
		Scan(&prod.Id, &prod.Id_user, &prod.Name, &prod.Qty, &prod.Price, &prod.Img, &prod.Descript)

	return *prod == DB.Product{}

}
