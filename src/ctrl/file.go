package ctrl

import "fmt"

func FileUpload(c Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return FiberBadReq(err)
	}

	for _, i := range form.File["img"] {
		if err := c.SaveFile(i, fmt.Sprintf("./public/%s", i.Filename)); err != nil {
			panic(err)
		}
	}
	return End(c)
}
