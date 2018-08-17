package dal

type Ifilter interface {
	GetName() string
	Apply(interface{}) (interface{}, error)
	GetType() string
}
