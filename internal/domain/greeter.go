package domain

import "fmt"

type Greeter interface {
	SayHello(name string) string
	GenerateHellos(name string, count int) []string
}

type SimpleGreeter struct{}

func NewSimpleGreeter() *SimpleGreeter {
	return &SimpleGreeter{}
}
func (s *SimpleGreeter) SayHello(name string) string {
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
func (s *SimpleGreeter) GenerateHellos(name string, count int) []string {
	hellos := make([]string, 0, count)
	for i := 0; i < count; i++ {
		hellos = append(hellos, fmt.Sprintf("Hello #%d from %s", i+1, name))
	}
	return hellos
}
