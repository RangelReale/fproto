package fdep

type DepType struct {
	FileDep *FileDep
	Alias   string
	Name    string
	Item    interface{}
}
