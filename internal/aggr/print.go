package aggr

type Printable interface {
	String() string
	IsEmpty() bool
}
