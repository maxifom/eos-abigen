package abitypes

type ABI struct {
	Version          string        `json:"version"`
	Types            []Type        `json:"types"`
	Structs          []Struct      `json:"structs"`
	Actions          []Action      `json:"actions"`
	Tables           []Table       `json:"tables"`
	RicardianClauses []interface{} `json:"ricardian_clauses"`
	ErrorMessages    []interface{} `json:"error_messages"`
	ABIExtensions    []interface{} `json:"abi_extensions"`
	Variants         []interface{} `json:"variants"`
	ActionResults    []interface{} `json:"action_results"`
	KvTables         KvTables      `json:"kv_tables"`
}

type Action struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	RicardianContract string `json:"ricardian_contract"`
}

type KvTables struct {
}

type Struct struct {
	Name   string  `json:"name"`
	Base   string  `json:"base"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Table struct {
	Name      string        `json:"name"`
	IndexType string        `json:"index_type"`
	KeyNames  []interface{} `json:"key_names"`
	KeyTypes  []interface{} `json:"key_types"`
	Type      string        `json:"type"`
}

type Type struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}
