package dt

// Value is the value interface.
type Value interface {
	Float() float64
	Bool() bool
	String() string
}
