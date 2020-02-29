package dt

// Value is the value interface.
type Value interface {
	Int() int
	Float() float64
	Bool() bool
	String() string
}
