package structure
// Todo: apply to all other collections
type CollectionInterface interface {
	GetId()	string
	Append(item interface{}) *Collection
	GetItems() []interface{}
	GetLength() uint64
}

type Collection struct {
	id     string
	items  []interface{}
	length uint64
	CollectionInterface
}

func NewCollection(id string) *Collection {
	return &Collection{
		id: id,
	}
}

func (self *Collection) GetId() string {
	return self.id
}

func (self *Collection) Append(item interface{}) *Collection {
	self.items = append(self.items, item)
	self.length = uint64(len(self.items))
	return self
}

func (self *Collection) AppendItems(items []interface{}) *Collection {
	self.items = append(self.items, items)
	self.length = uint64(len(self.items))
	return self
}

func (self *Collection) GetItems() []interface{} {
	return self.items
}

func (self *Collection) GetLength() uint64 {
	return self.length
}
