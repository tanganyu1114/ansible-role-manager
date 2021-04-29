package inventory

import (
	"bytes"
	"errors"
	"strings"
)

type InventoryFileParser interface {
	Parse(data []byte) (Group, error)
	Dump(g Group) ([]byte, error)
}

type inventoryFileParser struct {
}

func NewInventoryFileParser() InventoryFileParser {
	parser := new(inventoryFileParser)
	return InventoryFileParser(parser)
}

func (i inventoryFileParser) Parse(data []byte) (Group, error) {
	if data == nil || len(data) == 0 {
		return nil, errors.New("data is null")
	}

	g := newGroup()
	peek := 0
	isFirstGroup := true

	for idx := 0; idx < len(data); idx++ {
		if data[idx] == '\n' || idx == len(data)-1 {
			line := string(data[peek:idx])
			if idx == len(data)-1 {
				line = string(data[peek:])
			}
			peek = idx + 1
			line = strings.TrimSpace(line)
			if h := ParseHost(line); h != nil {
				_ = g.addHost(h)
			} else if line[0] == '[' && line[len(line)-1] == ']' {
				if !isFirstGroup {
					break
				}
				err := g.setName(line[1 : len(line)-1])
				if err != nil {
					return nil, err
				}
				isFirstGroup = false
			}
		}
	}
	return g, nil
}

func (i inventoryFileParser) Dump(g Group) ([]byte, error) {
	if g == nil || g.GetName() == "" {
		return nil, errors.New("nil Group input")
	}
	buff := bytes.NewBuffer(make([]byte, 0))
	buff.WriteString("[" + g.GetName() + "]\n")
	for _, h := range g.GetHosts() {
		buff.WriteString(h.GetIPString() + "\n")
	}
	return buff.Bytes(), nil
}
