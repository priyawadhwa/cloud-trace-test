package tmc

type TMC struct {
	m map[string]string
}

func New() *TMC {
	return &TMC{
		m: map[string]string{},
	}
}

// Get returns the value associated with the passed key.
func (t *TMC) Get(key string) string {
	v, _ := t.m[key]
	return v
}

// Set stores the key-value pair.
func (t *TMC) Set(key string, value string) {
	t.m[key] = value
}
