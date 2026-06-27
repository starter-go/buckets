package buckets

////////////////////////////////////////////////////////////////////////////////

type MetaName string

func (name MetaName) String() string {
	return string(name)
}

////////////////////////////////////////////////////////////////////////////////

const (
	MetaNamePrefix MetaName = "object."

	MetaObjectLength MetaName = MetaNamePrefix + "length"
	MetaObjectSum    MetaName = MetaNamePrefix + "sum"
	MetaObjectType   MetaName = MetaNamePrefix + "type"

	MetaObjectCreatedAt     MetaName = MetaNamePrefix + "created-at"      // as int like '123456789'
	MetaObjectCreatedAtTime MetaName = MetaNamePrefix + "created-at-time" // as text like '1970-01-01 00:00:00'
	MetaObjectUpdatedAt     MetaName = MetaNamePrefix + "updated-at"      // as int like '123456789'
	MetaObjectUpdatedAtTime MetaName = MetaNamePrefix + "updated-at-time" // as text like '1970-01-01 00:00:00'

	MetaObjectX MetaName = MetaNamePrefix + "object.x"
)

////////////////////////////////////////////////////////////////////////////////

type MetaMap map[MetaName]string

func (m MetaMap) Init() MetaMap {
	if m == nil {
		m = make(MetaMap)
	}
	return m
}

func (m MetaMap) Put(name MetaName, value string) {
	if m == nil {
		return
	}
	if (name == "") || (value == "") {
		return
	}
	m[name] = value
}

func (m MetaMap) Get(name MetaName) string {
	if m == nil {
		return ""
	}
	return m[name]
}

func (m MetaMap) Remove(name MetaName) {
	if m == nil {
		return
	}
	m[name] = ""
}

func (m MetaMap) Trim() MetaMap {

	src := m.Init()
	countEmpty := 0

	for k, v := range src {
		if k == "" || v == "" {
			countEmpty++
			break
		}
	}
	if countEmpty == 0 {
		return src
	}

	dst := make(MetaMap)
	for k, v := range src {
		dst.Put(k, v)
	}
	return dst
}

func (m MetaMap) Export(dst map[string]string) map[string]string {

	const empty = ""
	src := m

	if dst == nil {
		dst = make(map[string]string)
	}

	for k, v := range src {
		name := k.String()
		if name == empty || v == empty {
			continue
		}
		dst[name] = v
	}

	return dst
}

func (m MetaMap) Import(src map[string]string) {

	if m == nil {
		return
	}

	for k, v := range src {
		m.Put(MetaName(k), v)
	}
}

////////////////////////////////////////////////////////////////////////////////
// EOF
