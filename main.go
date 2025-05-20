package main

import (
	"fmt"
	"gop_shlyop/cache"
)

func main() {
	stringIntCache := cache.New[string, int]()

	stringIntCache.Set("apple", 10)
	stringIntCache.Set("banana", 20)
	fmt.Println("Добавлены 'apple': 10, 'banana': 20")

	if val, ok := stringIntCache.Get("apple"); ok {
		fmt.Printf("Значение для 'apple': %d (найдено: %t)\n", val, ok)
	} else {
		fmt.Printf("Значение для 'apple' не найдено (найдено: %t)\n", ok)
	}

	if val, ok := stringIntCache.Get("cherry"); ok {
		fmt.Printf("Значение для 'cherry': %d (найдено: %t)\n", val, ok)
	} else {
		fmt.Printf("Значение для 'cherry' не найдено (найдено: %t)\n", ok)
	}

	stringIntCache.Set("apple", 100)
	fmt.Println("Значение для 'apple' перезаписано на 100")

	if val, ok := stringIntCache.Get("apple"); ok {
		fmt.Printf("Новое значение для 'apple': %d (найдено: %t)\n", val, ok)
	}

	fmt.Println("\n--- Тестирование с другим типом данных ---")
	intStringCache := cache.New[int, string]()

	intStringCache.Set(1, "один")
	intStringCache.Set(2, "два")
	fmt.Println("Добавлены 1: 'один', 2: 'два'")

	if val, ok := intStringCache.Get(1); ok {
		fmt.Printf("Значение для ключа 1: '%s' (найдено: %t)\n", val, ok)
	}

	if _, ok := intStringCache.Get(3); !ok {
		fmt.Printf("Значение для ключа 3 не найдено (найдено: %t)\n", ok)
	}
} 
