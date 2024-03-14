package DB

type (
	User struct {
		Id        uint   `json:"id"`
		Email     string `json:"email" validate:"required,email,unique"`
		Password  string `json:"password" validate:"required,min=8,eqfield=Password2"`
		Password2 string `json:"password2"`
	}

	JSONImg struct {
		Name string `json:"name"`
		Size int    `json:"size"`
	}

	/// img will be JSONImg
	Product struct {
		Id       uint   `json:"id"` // in form this field will become index
		Id_user  uint   `json:"id_user"`
		Name     string `json:"name" validate:"required,min=5,uniqueProduct"`
		Qty      int    `json:"qty" validate:"required"`
		Price    int    `json:"price" validate:"required"`
		Img      string `json:"img" validate:"imageValidate"`
		Descript string `json:"descript"`
	}

	Selling struct {
		Id           uint   `json:"id"`
		Product_list string `json:"product_list"`
		Provit       uint   `json:"provit"`
		Date         string `json:"date"`
	}
)
