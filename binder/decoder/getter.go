package decoder

type Getter interface {
	Get(string) string
	Values(string) []string
}
