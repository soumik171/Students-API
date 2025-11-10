package types

type Student struct {
	Id    int64  // big int
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"required"`
}
