package stack

// Stack é uma estrutura de dados de pilha
type Stack struct {
	items []interface{}
}

// Push adiciona um elemento ao topo da pilha
func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

// Pop remove e retorna o elemento no topo da pilha
func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

// Peek retorna o elemento no topo da pilha sem removê-lo
func (s *Stack) Peek() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}

// IsEmpty verifica se a pilha está vazia
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size retorna o número de elementos na pilha
func (s *Stack) Size() int {
	return len(s.items)
}
