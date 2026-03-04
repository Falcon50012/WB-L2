package cut

// Config содержит параметры запуска утилиты.
type Config struct {
	Fields        string // -f
	Delimiter     string // -d
	OnlyDelimited bool   // -s
}
