package compiler

import (
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/op"
)

type loop struct {
	continuePos []int
	breakPos    []int
}

type Code struct {
	id           string
	name         string
	isNamed      bool
	parent       *Code
	children     []*Code
	symbols      *SymbolTable
	instructions []op.Code
	constants    []any
	names        []string
	source       string
	functionID   string

	// Used during compilation only
	loops      []*loop
	pipeActive bool
}

func (c *Code) ID() string {
	return c.id
}

func (c *Code) CodeName() string {
	return c.name
}

func (c *Code) addName(name string) uint16 {
	c.names = append(c.names, name)
	return uint16(len(c.names) - 1)
}

func (c *Code) IsNamed() bool {
	return c.isNamed
}

func (c *Code) FunctionID() string {
	return c.functionID
}

func (c *Code) Parent() *Code {
	return c.parent
}

func (c *Code) newChild(name, source, funcID string) *Code {
	child := &Code{
		id:         fmt.Sprintf("%s.%d", c.id, len(c.children)),
		name:       name,
		isNamed:    name != "",
		parent:     c,
		symbols:    c.symbols.NewChild(),
		source:     source,
		functionID: funcID,
	}
	c.children = append(c.children, child)
	return child
}

func (c *Code) InstructionCount() int {
	return len(c.instructions)
}

func (c *Code) Instruction(index int) op.Code {
	return c.instructions[index]
}

func (c *Code) ConstantsCount() int {
	return len(c.constants)
}

func (c *Code) Constant(index int) any {
	return c.constants[index]
}

func (c *Code) NameCount() int {
	return len(c.names)
}

func (c *Code) Name(index int) string {
	return c.names[index]
}

func (c *Code) Source() string {
	return c.source
}

func (c *Code) LocalsCount() int {
	return int(c.symbols.Count())
}

func (c *Code) Local(index int) *Symbol {
	return c.symbols.Symbol(uint16(index))
}

func (c *Code) GlobalsCount() int {
	return int(c.symbols.Root().Count())
}

func (c *Code) Global(index int) *Symbol {
	return c.symbols.Root().Symbol(uint16(index))
}

func (c *Code) Root() *Code {
	curr := c
	for curr.parent != nil {
		curr = curr.parent
	}
	return curr
}

func (c *Code) MarshalJSON() ([]byte, error) {
	state, err := stateFromCode(c)
	if err != nil {
		return nil, err
	}
	return json.Marshal(state)
}

func (c *Code) Flatten() []*Code {
	var codes []*Code
	codes = append(codes, c)
	for _, child := range c.children {
		codes = append(codes, child.Flatten()...)
	}
	return codes
}
