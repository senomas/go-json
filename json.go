package json

import (
	gson "encoding/json"
	"errors"
	"fmt"
)

// Object type
type Object struct {
	raw *gson.RawMessage
	m   map[string]*Object
}

var (
	// ErrNoData const
	ErrNoData = errors.New("no data")
)

// New Object
func New() *Object {
	return &Object{}
}

// Parse func
func (obj *Object) Parse(b []byte) error {
	r := gson.RawMessage(b)
	obj.raw = &r
	obj.m = nil
	return nil
}

// MarshalJSON func
func (obj *Object) MarshalJSON() ([]byte, error) {
	if obj.m != nil {
		b, err := gson.Marshal(obj.m)
		if err != nil {
			return nil, err
		}
		return b, err
	}
	return *obj.raw, nil
}

// UnmarshalJSON func
func (obj *Object) UnmarshalJSON(data []byte) error {
	r := gson.RawMessage(data)
	obj.raw = &r
	mr := make(map[string]*gson.RawMessage)
	err := gson.Unmarshal(data, &mr)
	if err != nil {
		return err
	}
	obj.m = make(map[string]*Object)
	for k, v := range mr {
		obj.m[k] = &Object{raw: v}
	}
	return nil
}

// BytesP func
func (obj *Object) BytesP() []byte {
	b, err := gson.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

// Bytes func
func (obj *Object) Bytes() ([]byte, error) {
	b, err := gson.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return b, err
}

// String func
func (obj *Object) String() (string, error) {
	b, err := gson.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func parseRaw(raw *gson.RawMessage) (map[string]*Object, error) {
	mr := make(map[string]*gson.RawMessage)
	err := gson.Unmarshal([]byte(*raw), &mr)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Object)
	for k, v := range mr {
		m[k] = &Object{raw: v}
	}
	return m, nil
}

// PutP func
func (obj *Object) PutP(value interface{}, p ...string) *Object {
	o, err := obj.Put(value, p...)
	if err != nil {
		panic(err)
	}
	return o
}

// Put func
func (obj *Object) Put(value interface{}, p ...string) (*Object, error) {
	var err error
	li := len(p) - 1
	op := obj
	for k, v := range p {
		if op.m == nil {
			if op.raw != nil {
				op.m, err = parseRaw(op.raw)
				if err != nil {
					return nil, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*op.raw))
				}
			} else {
				op.m = make(map[string]*Object)
			}
		}
		if k == li {
			var bb []byte
			bb, err = gson.Marshal(value)
			if err != nil {
				return op, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*op.raw))
			}
			r := gson.RawMessage(bb)
			op.m[v] = &Object{raw: &r}
			return op, nil
		}
		opp := op.m[v]
		if opp == nil {
			opp = &Object{}
			op.m[v] = opp
		}
		op = opp
	}
	return nil, nil
}

// GetP func
func (obj *Object) GetP(p ...string) *Object {
	o, err := obj.Get(p...)
	if err != nil {
		panic(err)
	}
	return o
}

// Get func
func (obj *Object) Get(p ...string) (*Object, error) {
	var result *Object
	li := len(p) - 1
	var m map[string]*Object
	var err error
	m, err = parseRaw(obj.raw)
	if err != nil {
		return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:1], string(*obj.raw))
	}
	opp := obj
	for k, v := range p {
		op := m[v]
		if op == nil {
			return result, fmt.Errorf("no data, path=%#v, value=\"%s\"", p[:k+1], string(*opp.raw))
		}
		if k == li {
			return op, nil
		}
		m, err = parseRaw(op.raw)
		if err != nil {
			return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
		}
		opp = op
		op = &Object{raw: op.raw, m: m}
	}
	return result, fmt.Errorf("invalid path, path=%#v", p)
}

// GetStringP func
func (obj *Object) GetStringP(p ...string) string {
	o, err := obj.GetString(p...)
	if err != nil {
		panic(err)
	}
	return o
}

// GetString func
func (obj *Object) GetString(p ...string) (string, error) {
	result := ""
	li := len(p) - 1
	var m map[string]*Object
	var err error
	m, err = parseRaw(obj.raw)
	if err != nil {
		return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:1], string(*obj.raw))
	}
	opp := obj
	for k, v := range p {
		op := m[v]
		if op == nil {
			return result, fmt.Errorf("no data, path=%#v, value=\"%s\"", p[:k+1], string(*opp.raw))
		}
		if k == li {
			err = gson.Unmarshal(*op.raw, &result)
			if err != nil {
				return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
			}
			return result, nil
		}
		m, err = parseRaw(op.raw)
		if err != nil {
			return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
		}
		opp = op
		op = &Object{raw: op.raw, m: m}
	}
	return result, fmt.Errorf("invalid path, path=%#v", p)
}

// GetIntP func
func (obj *Object) GetIntP(p ...string) int {
	o, err := obj.GetInt(p...)
	if err != nil {
		panic(err)
	}
	return o
}

// GetInt func
func (obj *Object) GetInt(p ...string) (int, error) {
	result := int(0)
	li := len(p) - 1
	var m map[string]*Object
	var err error
	m, err = parseRaw(obj.raw)
	if err != nil {
		return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:1], string(*obj.raw))
	}
	opp := obj
	for k, v := range p {
		op := m[v]
		if op == nil {
			return result, fmt.Errorf("no data, path=%#v, value=\"%s\"", p[:k+1], string(*opp.raw))
		}
		if k == li {
			err = gson.Unmarshal(*op.raw, &result)
			if err != nil {
				return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
			}
			return result, nil
		}
		m, err = parseRaw(op.raw)
		if err != nil {
			return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
		}
		opp = op
		op = &Object{raw: op.raw, m: m}
	}
	return result, fmt.Errorf("invalid path, path=%#v", p)
}

// GetInt64P func
func (obj *Object) GetInt64P(p ...string) int64 {
	o, err := obj.GetInt64(p...)
	if err != nil {
		panic(err)
	}
	return o
}

// GetInt64 func
func (obj *Object) GetInt64(p ...string) (int64, error) {
	result := int64(0)
	li := len(p) - 1
	var m map[string]*Object
	var err error
	m, err = parseRaw(obj.raw)
	if err != nil {
		return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:1], string(*obj.raw))
	}
	opp := obj
	for k, v := range p {
		op := m[v]
		if op == nil {
			return result, fmt.Errorf("no data, path=%#v, value=\"%s\"", p[:k+1], string(*opp.raw))
		}
		if k == li {
			err = gson.Unmarshal(*op.raw, &result)
			if err != nil {
				return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
			}
			return result, nil
		}
		m, err = parseRaw(op.raw)
		if err != nil {
			return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
		}
		opp = op
		op = &Object{raw: op.raw, m: m}
	}
	return result, fmt.Errorf("invalid path, path=%#v", p)
}

// GetUint64P func
func (obj *Object) GetUint64P(p ...string) uint64 {
	o, err := obj.GetUint64(p...)
	if err != nil {
		panic(err)
	}
	return o
}

// GetUint64 func
func (obj *Object) GetUint64(p ...string) (uint64, error) {
	result := uint64(0)
	li := len(p) - 1
	var m map[string]*Object
	var err error
	m, err = parseRaw(obj.raw)
	if err != nil {
		return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:1], string(*obj.raw))
	}
	opp := obj
	for k, v := range p {
		op := m[v]
		if op == nil {
			return result, fmt.Errorf("no data, path=%#v, value=\"%s\"", p[:k+1], string(*opp.raw))
		}
		if k == li {
			err = gson.Unmarshal(*op.raw, &result)
			if err != nil {
				return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
			}
			return result, nil
		}
		m, err = parseRaw(op.raw)
		if err != nil {
			return result, fmt.Errorf("%v, path=%#v, value=\"%s\"", err, p[:k+1], string(*opp.raw))
		}
		opp = op
		op = &Object{raw: op.raw, m: m}
	}
	return result, fmt.Errorf("invalid path, path=%#v", p)
}
