package unit

var types = map[string]Schema{}

func RegisterUnitType(name string, sch Schema) {
	types[name] = sch
}

func GetSchema(typ string) (Schema, bool) {
	v, ok := types[typ]
	return v, ok
}
