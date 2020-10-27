package tmc

type TMC struct {
	M map[string]string
}

func New() *TMC {
	return &TMC{
		M: map[string]string{},
	}
}

// Get returns the value associated with the passed key.
func (t *TMC) Get(key string) string {
	v, _ := t.M[key]
	return v
}

// Set stores the key-value pair.
func (t *TMC) Set(key string, value string) {
	t.M[key] = value
}
