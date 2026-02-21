package sorter

// Options содержит параметры сортировки
type Options struct {
	Delimiter string // -d
	Column    int    // -k
	Numeric   bool   // -n
	Reverse   bool   // -r
	Unique    bool   // -u
}
