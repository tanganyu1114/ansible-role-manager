package inventory

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
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
	groupNameReg := regexp.MustCompile(`^\[(\S+)\]`)
	ipReg := regexp.MustCompile(`^\s*(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`)
	peek := 0
	isFirstGroup := true

	for idx := 0; idx < len(data); idx++ {
		if data[idx] == '\n' || idx == len(data)-1 {
			line := string(data[peek:idx])
			if idx == len(data)-1 {
				line = string(data[peek:])
			}
			peek = idx + 1
			if match := groupNameReg.FindStringSubmatch(line); len(match) == 2 {
				if !isFirstGroup {
					break
				}
				err := g.setName(match[1])
				if err != nil {
					return nil, err
				}
				isFirstGroup = false
			} else if match := ipReg.FindStringSubmatch(line); len(match) == 5 {
				ip := [4]byte{}
				for j := 1; j < 5; j++ {
					b, err := strconv.Atoi(match[j])
					if err == nil && b > 255 {
						err = errors.New("no matching to ip")
					}
					if err != nil {
						return nil, err
					}
					ip[j-1] = byte(b)
				}
				_ = g.addHost(NewIPv4Host(ip))
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
		buff.WriteString(h.GetIp().IP.String() + "\n")
	}
	return buff.Bytes(), nil
}
