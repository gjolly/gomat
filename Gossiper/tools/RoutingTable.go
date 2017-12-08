package tools

type RoutingTable struct {
	table map[string]string
}

func newRoutingTable() *RoutingTable {
	return &RoutingTable{make(map[string]string, 0)}
}

func (r *RoutingTable) add(key string, value string) {
	r.table[key] = value
}

func (r RoutingTable) FindNextHop(dest string) string {
	return r.table[dest]
}

func (r RoutingTable) GetTable() map[string]string {
	return r.table
}