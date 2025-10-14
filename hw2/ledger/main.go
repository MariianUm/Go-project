package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Transaction структура для представления транзакции
type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

// Budget структура для хранения информации о бюджете
type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period"`
}

// Глобальное хранилище транзакций в памяти
var transactions []Transaction

// Глобальное хранилище бюджетов
var budgets = make(map[string]Budget)

// AddTransaction добавляет новую транзакцию с проверкой бюджета
func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return fmt.Errorf("сумма транзакции не может быть нулевой")
	}
	
	// Проверка бюджета
	if budget, exists := budgets[tx.Category]; exists {
		currentTotal := calculateCategoryTotal(tx.Category)
		if currentTotal+tx.Amount > budget.Limit {
			return fmt.Errorf("бюджет превышен для категории '%s': текущие %.2f + новые %.2f > лимит %.2f", 
				tx.Category, currentTotal, tx.Amount, budget.Limit)
		}
	}
	
	tx.ID = len(transactions) + 1
	
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}
	
	transactions = append(transactions, tx)
	return nil
}

// calculateCategoryTotal вычисляет текущую сумму транзакций по категории
func calculateCategoryTotal(category string) float64 {
	total := 0.0
	for _, tx := range transactions {
		if tx.Category == category {
			total += tx.Amount
		}
	}
	return total
}

// ListTransactions возвращает все транзакции
func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

// SetBudget добавляет или обновляет бюджет в хранилище
func SetBudget(b Budget) error {
	if b.Category == "" {
		return errors.New("категория бюджета не может быть пустой")
	}
	if b.Limit <= 0 {
		return errors.New("лимит бюджета должен быть положительным")
	}
	
	budgets[b.Category] = b
	return nil
}

// GetBudget возвращает бюджет для категории
func GetBudget(category string) (Budget, bool) {
	budget, exists := budgets[category]
	return budget, exists
}

// ListBudgets возвращает все бюджеты
func ListBudgets() map[string]Budget {
	result := make(map[string]Budget)
	for k, v := range budgets {
		result[k] = v
	}
	return result
}

// LoadBudgets загружает бюджеты из JSON
func LoadBudgets(r io.Reader) error {
	var budgetList []Budget
	
	reader := bufio.NewReader(r)
	decoder := json.NewDecoder(reader)
	
	if err := decoder.Decode(&budgetList); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %v", err)
	}
	
	for _, budget := range budgetList {
		if err := SetBudget(budget); err != nil {
			return fmt.Errorf("ошибка установки бюджета для %s: %v", budget.Category, err)
		}
	}
	
	return nil
}

// GetCategoryTotal возвращает текущую сумму по категории
func GetCategoryTotal(category string) float64 {
	return calculateCategoryTotal(category)
}

func main() {
	fmt.Println("Сервис Ledger запущен с системой бюджетирования")

	// Установка бюджетов через код
	fmt.Println("\n1. Установка начальных бюджетов через код...")
	initialBudgets := []Budget{
		{Category: "Еда", Limit: 5000, Period: "monthly"},
		{Category: "Развлечения", Limit: 3000, Period: "monthly"},
		{Category: "Транспорт", Limit: 2000, Period: "monthly"},
	}
	
	for _, budget := range initialBudgets {
		if err := SetBudget(budget); err != nil {
			fmt.Printf("Ошибка установки бюджета: %v\n", err)
		} else {
			fmt.Printf("Бюджет установлен: %s - %.2f\n", budget.Category, budget.Limit)
		}
	}

	// Загрузка бюджетов из файла JSON
	fmt.Println("\n2. Загрузка дополнительных бюджетов из JSON файла...")
	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Printf("Предупреждение: не удалось открыть budgets.json: %v\n", err)
	} else {
		defer file.Close()
		if err := LoadBudgets(file); err != nil {
			fmt.Printf("Ошибка загрузки бюджетов из файла: %v\n", err)
		} else {
			fmt.Println("Бюджеты успешно загружены из JSON файла")
		}
	}

	// Выводим все установленные бюджеты
	fmt.Println("\n3. Текущие бюджеты:")
	for category, budget := range ListBudgets() {
		currentTotal := GetCategoryTotal(category)
		fmt.Printf("  %s: Лимит %.2f, Текущие расходы %.2f, Осталось %.2f\n", 
			category, budget.Limit, currentTotal, budget.Limit-currentTotal)
	}

	// Тестовые сценарии
	fmt.Println("\n4. Запуск тестовых сценариев...")

	// Сценарий 1: Успешное добавление транзакции в пределах бюджета
	fmt.Println("\n--- Сценарий 1: Транзакция в пределах бюджета ---")
	tx1 := Transaction{
		Amount:      1000,
		Category:    "Еда",
		Description: "Продукты",
		Date:        time.Now(),
	}
	
	if err := AddTransaction(tx1); err != nil {
		fmt.Printf("❌ Не удалось добавить транзакцию: %v\n", err)
	} else {
		fmt.Printf("✅ Транзакция успешно добавлена: %.2f для %s\n", tx1.Amount, tx1.Category)
	}

	// Сценарий 2: Еще одна успешная транзакция
	fmt.Println("\n--- Сценарий 2: Еще одна транзакция в пределах бюджета ---")
	tx2 := Transaction{
		Amount:      2000,
		Category:    "Еда",
		Description: "Ресторан",
	}
	
	if err := AddTransaction(tx2); err != nil {
		fmt.Printf("❌ Не удалось добавить транзакцию: %v\n", err)
	} else {
		fmt.Printf("✅ Транзакция успешно добавлена: %.2f для %s\n", tx2.Amount, tx2.Category)
	}

	// Сценарий 3: Транзакция, превышающая бюджет
	fmt.Println("\n--- Сценарий 3: Транзакция, превышающая бюджет ---")
	tx3 := Transaction{
		Amount:      2500,
		Category:    "Еда",
		Description: "Дорогой ужин",
	}
	
	if err := AddTransaction(tx3); err != nil {
		fmt.Printf("✅ Корректно отклонено: %v\n", err)
	} else {
		fmt.Printf("❌ Транзакция должна была быть отклонена!\n")
	}

	// Сценарий 4: Транзакция в другой категории (без бюджета)
	fmt.Println("\n--- Сценарий 4: Транзакция в категории без бюджета ---")
	tx4 := Transaction{
		Amount:      5000,
		Category:    "Здравоохранение",
		Description: "Медосмотр",
	}
	
	if err := AddTransaction(tx4); err != nil {
		fmt.Printf("❌ Не удалось добавить транзакцию: %v\n", err)
	} else {
		fmt.Printf("✅ Транзакция успешно добавлена: %.2f для %s\n", tx4.Amount, tx4.Category)
	}

	// Сценарий 5: Обновление бюджета и повторная попытка
	fmt.Println("\n--- Сценарий 5: Увеличение бюджета и повторная попытка ---")
	newBudget := Budget{Category: "Еда", Limit: 6000, Period: "monthly"}
	if err := SetBudget(newBudget); err != nil {
		fmt.Printf("Ошибка обновления бюджета: %v\n", err)
	} else {
		fmt.Printf("Бюджет обновлен: %s - %.2f\n", newBudget.Category, newBudget.Limit)
		
		// Пробуем снова добавить транзакцию
		if err := AddTransaction(tx3); err != nil {
			fmt.Printf("❌ Все еще отклонено: %v\n", err)
		} else {
			fmt.Printf("✅ Теперь принято с увеличенным бюджетом: %.2f для %s\n", tx3.Amount, tx3.Category)
		}
	}

	// Финальный статус
	fmt.Println("\n5. Финальное состояние:")
	fmt.Println("Бюджеты:")
	for category, budget := range ListBudgets() {
		currentTotal := GetCategoryTotal(category)
		fmt.Printf("  %s: Лимит %.2f, Текущие расходы %.2f, Осталось %.2f\n", 
			category, budget.Limit, currentTotal, budget.Limit-currentTotal)
	}

	fmt.Println("\nВсе транзакции:")
	allTx := ListTransactions()
	for _, tx := range allTx {
		fmt.Printf("ID: %d, Сумма: $%.2f, Категория: %s, Описание: %s, Дата: %s\n",
			tx.ID, tx.Amount, tx.Category, tx.Description, tx.Date.Format("2006-01-02"))
	}

	fmt.Printf("\nВсего транзакций: %d\n", len(allTx))
}
