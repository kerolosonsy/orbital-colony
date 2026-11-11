package main

import (
	"fmt"
	"orbital/colony"
	"strconv"
)

type FlightDrive interface{ Step(colony.Direction) bool }
type FlightExpr interface {
	Eval(map[string]int64) (int64, error)
}
type flightAtom struct{ token CommandToken }

func (a flightAtom) Eval(env map[string]int64) (int64, error) {
	if a.token.Kind == "number" {
		return strconv.ParseInt(a.token.Text, 10, 64)
	}
	n, ok := env[a.token.Text]
	if !ok {
		return 0, &CommandError{a.token.Offset, "undefined variable"}
	}
	return n, nil
}

type flightBinary struct {
	op          CommandToken
	left, right FlightExpr
}

func (b flightBinary) Eval(env map[string]int64) (int64, error) {
	left, err := b.left.Eval(env)
	if err != nil {
		return 0, err
	}
	right, err := b.right.Eval(env)
	if err != nil {
		return 0, err
	}
	var n int64
	switch b.op.Kind {
	case "+":
		n = left + right
	case "-":
		n = left - right
	case "*":
		n = left * right
	case "/":
		if right == 0 {
			return 0, &CommandError{b.op.Offset, "division by zero"}
		}
		n = left / right
	}
	if n < -1000000 || n > 1000000 {
		return 0, &CommandError{b.op.Offset, "expression out of range"}
	}
	return n, nil
}

type flightStatement struct {
	command, target CommandToken
	expr            FlightExpr
}
type flightParser struct {
	tokens []CommandToken
	next   int
}

func (p *flightParser) peek() CommandToken { return p.tokens[p.next] }
func (p *flightParser) take(kind string) (CommandToken, error) {
	t := p.peek()
	if t.Kind != kind {
		return t, &CommandError{t.Offset, "expected " + kind}
	}
	p.next++
	return t, nil
}
func flightPrimary(p *flightParser) (FlightExpr, error) {
	t := p.peek()
	switch t.Kind {
	case "number", "name":
		p.next++
		return flightAtom{t}, nil
	case "(":
		p.next++
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.take(")")
		return expr, err
	case "+", "-":
		p.next++
		expr, err := flightPrimary(p)
		if err != nil {
			return nil, err
		}
		return flightBinary{t, flightAtom{CommandToken{"number", "0", t.Offset}}, expr}, nil
	}
	return nil, &CommandError{t.Offset, "expected expression"}
}
func (p *flightParser) term() (FlightExpr, error) {
	left, err := flightPrimary(p)
	if err != nil {
		return nil, err
	}
	for p.peek().Kind == "*" || p.peek().Kind == "/" {
		op := p.peek()
		p.next++
		right, err := flightPrimary(p)
		if err != nil {
			return nil, err
		}
		left = flightBinary{op, left, right}
	}
	return left, nil
}
func (p *flightParser) expression() (FlightExpr, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.peek().Kind == "+" || p.peek().Kind == "-" {
		op := p.peek()
		p.next++
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		left = flightBinary{op, left, right}
	}
	return left, nil
}
func (p *flightParser) statement() (flightStatement, error) {
	command, err := p.take("name")
	if err != nil {
		return flightStatement{}, err
	}
	if command.Text != "let" && command.Text != "thrust" {
		return flightStatement{}, &CommandError{command.Offset, "unknown command"}
	}
	target, err := p.take("name")
	if err != nil {
		return flightStatement{}, err
	}
	if command.Text == "let" {
		if _, err = p.take("="); err != nil {
			return flightStatement{}, err
		}
	} else if _, ok := flightDirections[target.Text]; !ok {
		return flightStatement{}, &CommandError{target.Offset, "unknown direction"}
	}
	expr, err := p.expression()
	if err != nil {
		return flightStatement{}, err
	}
	if _, err = p.take(";"); err != nil {
		return flightStatement{}, err
	}
	return flightStatement{command, target, expr}, nil
}

var flightDirections = map[string]colony.Direction{"east": colony.East, "north": colony.North, "west": colony.West, "south": colony.South}

func (s flightStatement) run(env map[string]int64, drive FlightDrive) error {
	n, err := s.expr.Eval(env)
	if err != nil {
		return err
	}
	if s.command.Text == "let" {
		env[s.target.Text] = n
		return nil
	}
	if n < 0 || n > 32 {
		return &CommandError{s.command.Offset, "thrust outside 0..32"}
	}
	for i := int64(0); i < n; i++ {
		if !drive.Step(flightDirections[s.target.Text]) {
			return &CommandError{s.command.Offset, "blocked thrust"}
		}
	}
	return nil
}
func ExecuteFlight(source string, drive FlightDrive) error {
	if drive == nil {
		return fmt.Errorf("drive required")
	}
	tokens, err := LexFlight(source)
	if err != nil {
		return err
	}
	parser := flightParser{tokens: tokens}
	var program []flightStatement
	for parser.peek().Kind != "eof" {
		statement, err := parser.statement()
		if err != nil {
			return err
		}
		program = append(program, statement)
	}
	env := map[string]int64{}
	for _, statement := range program {
		if err := statement.run(env, drive); err != nil {
			return err
		}
	}
	return nil
}
