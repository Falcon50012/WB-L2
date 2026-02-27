package grep

type Config struct {
	Pattern     string
	After       int
	Before      int
	Context     int
	IgnoreCase  bool
	InvertMatch bool
	LineNumber  bool
	Fixed       bool
	Count       bool
}
