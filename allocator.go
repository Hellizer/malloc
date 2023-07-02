package malloc

//модуль для работы с сохранением объектов в неуправляемую память
// type idx struct {
// 	mapPart int    //контейнер
// 	number  uint64 //номер в контейнере
// }

//Интерфейс должен бфть реализован в сохраняемом объекте
type MemObj interface {
	GetID() uint64       //получение ID объекта
	PutData(ptr uintptr) //сохранение данных из объекта в память по указателю ptr
}

//Структура для Хранилища
type Allocator struct {
	memParts map[int]*memSpace 			//контейнеры выделенной памяти
	index         map[uint64]uintptr 	//
	itemSize      uint64             	//размер сохраняемого элемента
	containerSize uint64             	//количество элементов в контейнере
	lastMemSpace  int                	//последний используемый контейнер
	iterator      uint64             	//итератор перебора по кругу
}

//Создает новое Хранилище с заданными параметрами размер элемента и количество элементов в каждом контейнере
func NewAllocator(itemSize uint64, containerSize uint64) *Allocator {
	alloc := new(Allocator)
	alloc.memParts = make(map[int]*memSpace, 1)
	alloc.itemSize = itemSize
	alloc.containerSize = containerSize
	alloc.memParts[0] = newMemSpace(itemSize * containerSize)
	alloc.index = make(map[uint64]uintptr)
	alloc.iterator = 0
	return alloc
}

//Освобождает выделенную память
func (a *Allocator) Free() {
	for _, v := range a.memParts {
		v.free()
	}
}

//Помещает объект в память и возвращает false если объект с таким ID уже существует
func (a *Allocator) Put(obj MemObj) bool {
	id := obj.GetID()
	_, ok := a.index[id]
	if ok {
		return false
	}

	if a.memParts[a.lastMemSpace].count >= a.containerSize-1 {
		a.lastMemSpace++
		a.memParts[a.lastMemSpace] = newMemSpace(a.itemSize * a.containerSize)
	}

	ptr := a.memParts[a.lastMemSpace].pointer + uintptr((a.memParts[a.lastMemSpace].count * a.itemSize))

	obj.PutData(ptr)
	a.index[id] = ptr
	a.memParts[a.lastMemSpace].count++
	return true
}

//Возвращает uintptr на объект или 0 если объект не найден
func (a *Allocator) Get(id uint64) uintptr {
	v, ok := a.index[id]
	if ok {
		return v
	}
	return 0
}

//выделение места под пустой объект с заданным ID
func (a *Allocator) NewRecord(id uint64) uintptr {
	_, ok := a.index[id]
	if ok {
		return 0
	}
	if a.memParts[a.lastMemSpace].count >= a.containerSize-1 {
		a.lastMemSpace++
		a.memParts[a.lastMemSpace] = newMemSpace(a.itemSize * a.containerSize)
	}
	ptr := a.memParts[a.lastMemSpace].pointer + uintptr((a.memParts[a.lastMemSpace].count * a.itemSize))
	a.index[id] = ptr
	a.memParts[a.lastMemSpace].count++
	return ptr
}

func (a *Allocator) GetNext() uintptr {
	container := (a.iterator) / (a.containerSize)
	number := (a.iterator) % (a.containerSize)
	a.iterator++
	if a.iterator >= uint64(len(a.index)) {
		a.iterator = 0
	}
	return a.memParts[int(container)].pointer + uintptr(a.itemSize*number)
}

func (a *Allocator) Count() uint64 {
	return uint64(len(a.index))
}
