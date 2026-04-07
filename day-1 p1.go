package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

// ======================= BASIC STRUCTS =======================

// Structs bundle related data together (like classes)
// Used to create custom data types
type Person struct {
	Name  string
	Age   int
	Email string
}

// Value receiver → works on copy (read-only)
func (p Person) Greet() string {
	return "Hi, I'm " + p.Name
}

// Pointer receiver → modifies original struct
func (p *Person) Birthday() {
	p.Age++
}

// ======================= EMBEDDING =======================

// Embedding allows reuse of fields/methods (like inheritance)
type Animal struct {
	Name string
}

func (a Animal) Speak() string {
	return a.Name + " speaks"
}

type Dog struct {
	Animal
	Breed string
}

// ======================= STRUCT TAGS =======================

// Struct tags provide metadata for JSON/DB/etc.
// `omitempty` skips field if empty
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

// ======================= FUNCTIONS =======================

// Multiple return values → common Go pattern (value + error)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// Variadic function → accepts variable number of arguments
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Named return values → declared in signature
// Returned implicitly using 'return'
func minMax(nums []int) (min, max int) {
	min, max = nums[0], nums[0]
	for _, v := range nums {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

// ======================= CLOSURES =======================

// Closure remembers outer variables (state)
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// ======================= DEFER =======================

// defer executes at end of function (LIFO order)
// Commonly used for cleanup
func readFile(name string) {
	f, _ := os.Open(name)
	defer f.Close()
}

// ======================= INTERFACES =======================

// Interface defines method signatures
// Types implement implicitly (no 'implements' keyword)
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle implements Shape
type Rectangle struct{ W, H float64 }

func (r Rectangle) Area() float64      { return r.W * r.H }
func (r Rectangle) Perimeter() float64 { return 2 * (r.W + r.H) }

// Circle implements Shape
type Circle struct{ Radius float64 }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Interface as parameter → polymorphism
func printInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// ======================= STRINGER =======================

// Stringer interface customizes printing
// If String() exists → fmt uses it automatically
func (p Person) String() string {
	return fmt.Sprintf("%s (age %d)", p.Name, p.Age)
}

// ======================= TYPE ASSERTION =======================

// Type assertion extracts concrete type from interface
func describeShape(s Shape) {
	switch v := s.(type) {
	case Circle:
		fmt.Println("Circle radius:", v.Radius)
	case Rectangle:
		fmt.Println("Rectangle:", v.W, v.H)
	default:
		fmt.Println("Unknown shape")
	}
}

// ======================= EMPTY INTERFACE =======================

// interface{} (or any) accepts any type
// Use only when type is unknown (avoid overuse)
func printAny(v any) {
	fmt.Println(v)
}

// ======================= ERROR HANDLING =======================

// Error is an interface → any type with Error() is error
func sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("cannot sqrt negative")
	}
	return math.Sqrt(x), nil
}

// fmt.Errorf wraps error with context (%w)
func loadConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("loadConfig: %w", err)
	}
	return string(data), nil
}

// ======================= CUSTOM ERRORS =======================

// Custom error → carries structured data
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ======================= PANIC & RECOVER =======================

// panic → stops execution (use only for critical failures)
// recover → catches panic inside defer
func safeRun(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	fn()
}

// ======================= REAL WORLD EXAMPLE =======================

// Custom error for business logic
type InsufficientFundsError struct {
	Requested float64
	Available float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("need %.2f, have %.2f", e.Requested, e.Available)
}

type BankAccount struct {
	Owner   string
	Balance float64
}

func (b *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if amount > b.Balance {
		return &InsufficientFundsError{amount, b.Balance}
	}
	b.Balance -= amount
	return nil
}

// ======================= MAIN =======================

func main() {

	// Type casting → explicit in Go
	var a int = 42
	b := float64(a)
	c := string(rune(65))
	fmt.Println(a, b, c)

	// Slice → dynamic array
	nums := []int{1, 2, 3}
	nums = append(nums, 4)
	fmt.Println(nums)

	// Map → key-value store
	scores := map[string]int{"Alice": 90}
	val, ok := scores["Bob"]
	if ok {
		fmt.Println(val)
	} else {
		fmt.Println("not found")
	}

	// Functions
	res, _ := divide(10, 2)
	fmt.Println(res)

	fmt.Println(sum(1, 2, 3))
	min, max := minMax([]int{1, 5, 2})
	fmt.Println(min, max)

	// Closure
	cntr := counter()
	fmt.Println(cntr(), cntr())

	// Defer
	defer fmt.Println("last")

	// Struct
	p1 := Person{Name: "Madhu", Age: 21}
	p1.Birthday()
	fmt.Println(p1)

	// Interface
	printInfo(Rectangle{5, 3})
	printInfo(Circle{4})

	describeShape(Circle{3})

	// Empty interface
	printAny(42)
	printAny("hello")

	// Error handling
	_, err := sqrt(-1)
	if err != nil {
		fmt.Println(err)
	}

	// Custom error usage
	acc := &BankAccount{"Madhu", 1000}
	err = acc.Withdraw(1500)

	var ife *InsufficientFundsError
	if errors.As(err, &ife) {
		fmt.Println("Short by:", ife.Requested-ife.Available)
	}

	// JSON
	data, _ := json.Marshal(User{ID: 1, Name: "Madhu"})
	fmt.Println(string(data))
}
