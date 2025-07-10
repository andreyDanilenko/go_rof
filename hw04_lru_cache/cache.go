package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// cacheItem хранит ключ и значение элемента кэша.
type cacheItem struct {
	key   Key         // Ключ нужен для удаления из map при выталкивании
	value interface{} // Само значение элемента
}

type LruCache struct {
	Capacity int               // Максимальная ёмкость кэша
	Queue    List              // Очередь на основе двусвязного списка
	Items    map[Key]*ListItem // Словарь для быстрого доступа к элементам
}

// lruCache реализует LRU-кэш.
type lruCache struct {
	capacity int               // Максимальная ёмкость кэша
	queue    List              // Очередь на основе двусвязного списка
	items    map[Key]*ListItem // Словарь для быстрого доступа к элементам
}

// NewCache создаёт новый LRU-кэш заданной ёмкости.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(), // Используем раннюю реализацию работы со списком
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	// Проверяем, есть ли уже такой ключ в кэше
	if item, exists := c.items[key]; exists {
		// Обновляем значение
		item.Value.(*cacheItem).value = value
		// Перемещаем элемент в начало списка (как недавно использованный)
		c.queue.MoveToFront(item) // Используем метод MoveToFront
		return true
	}

	// Создаём новый элемент кэша
	newCacheItem := &cacheItem{
		key:   key,
		value: value,
	}

	// Добавляем в начало списка
	listItem := c.queue.PushFront(newCacheItem) // Используем метод PushFront
	c.items[key] = listItem

	// Если превысили ёмкость, удаляем последний элемент
	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back() // Используем метод Back
		if lastItem != nil {
			// Удаляем из словаря
			delete(c.items, lastItem.Value.(*cacheItem).key)
			// Удаляем из списка
			c.queue.Remove(lastItem) // Используем метод Remove
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		// Перемещаем элемент в начало списка (как недавно использованный)
		c.queue.MoveToFront(item) // Используем метод MoveToFront
		// Возвращаем значение
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	// Очищаем словарь, создавая новый с той же capacity
	c.items = make(map[Key]*ListItem, c.capacity)
	// Очищаем список, создавая новый
	c.queue = NewList() // Используем функцию NewList
}
