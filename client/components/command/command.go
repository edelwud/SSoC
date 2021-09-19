package command

type Command interface {
	Row() []byte
}
