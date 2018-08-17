package dal

type IDestination interface {
	Exec(payload interface{})
	GetType() string
	GetName() string
}
