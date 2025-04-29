package helpers

func PanicHelper(err error) {
	if err != nil {
		panic(err)
	}
}
