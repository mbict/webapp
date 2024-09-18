package binder

type (
	getter interface {
		Get(string) string
		Values(string) []string
	}

	mapGetter map[string][]string
)

func (m mapGetter) Get(key string) string {
	if m == nil {
		return ""
	}

	if vs, ok := m[key]; ok && len(vs) >= 1 {
		return vs[0]
	}

	return ""
}

func (m mapGetter) Values(key string) []string {
	if m == nil {
		return nil
	}

	return m[key]
}
