package main

import (
	"fmt"
	"strconv"
)

// Prime возвращает простые множители числа x
func Prime(x int) []int {
	var i int = 2
	var a []int
	for i*i <= x {
		for x%i == 0 {
			a = append(a, i)
			x /= i
		}
		i++
	}
	if x > 1 {
		a = append(a, x)
	}
	return a
}

// secret_key вычисляет секретный ключ по открытому ключу и модулю
func secret_key(open_key int, mod int) int {
	var Matr [2][3]int
	var mult []int

	Matr[0][0] = 1
	Matr[1][1] = 1
	Matr[0][2] = open_key

	mult = Prime(mod)
	if len(mult) < 2 {
		return -1 // Ошибка: модуль должен иметь как минимум два простых множителя
	}

	u := float32(mult[0])
	v := float32(mult[1])
	md := float32(mod)
	md2 := int(md * (1 - 1.0/u) * (1 - 1.0/v))
	Matr[1][2] = md2

	Matr = methodGaus(Matr)

	if Matr[0][2] == 1 {
		return Matr[0][0]
	} else if Matr[1][2] == 1 {
		return Matr[1][0]
	}
	return -1
}

// Printmat выводит матрицу для отладки
func Printmat(x [2][3]int) {
	for _, v := range x {
		fmt.Println(v)
	}
	fmt.Println("_____")
}

// methodGaus реализует расширенный алгоритм Евклида для нахождения обратного элемента
func methodGaus(x [2][3]int) [2][3]int {
	var mx, mn, ch int

	for !((x[0][2] == 0 && x[1][2] == 1) || (x[0][2] == 1 && x[1][2] == 0)) {
		if x[0][2] > x[1][2] {
			mx = 0
			mn = 1
		} else {
			mx = 1
			mn = 0
		}

		if x[mn][2] == 0 {
			break // Избегаем деления на ноль
		}

		ch = x[mx][2] / x[mn][2]
		for i := 0; i < 3; i++ {
			x[mx][i] = x[mx][i] - x[mn][i]*ch
		}
	}

	// Корректировка для получения положительного результата
	if x[0][2] == 1 {
		if x[0][0] > 0 {
			return x
		} else if x[1][0] > 0 {
			for x[0][0] <= 0 {
				for i := 0; i < 3; i++ {
					x[0][i] += x[1][i]
				}
			}
		} else if x[1][0] < 0 {
			x[1][0], x[1][1], x[1][2] = -x[1][0], -x[1][1], -x[1][2]
			for x[0][0] <= 0 {
				for i := 0; i < 3; i++ {
					x[0][i] += x[1][i]
				}
			}
		}
		return x
	}

	if x[1][2] == 1 {
		if x[1][0] > 0 {
			return x
		} else if x[0][0] > 0 {
			for x[1][0] <= 0 {
				for i := 0; i < 3; i++ {
					x[1][i] += x[0][i]
				}
			}
		} else if x[0][0] < 0 {
			x[0][0], x[0][1], x[0][2] = -x[0][0], -x[0][1], -x[0][2]
			for x[1][0] <= 0 {
				for i := 0; i < 3; i++ {
					x[1][i] += x[0][i]
				}
			}
		}
		return x
	}
	return x
}

// Enc реализует быстрое возведение в степень по модулю
func Enc(orig int, pow int64, mod int) int {
	bin := strconv.FormatInt(pow, 2)
	l1 := []rune(bin)
	k := len(l1)
	d1 := make([]int, k)
	d1[k-1] = orig % mod

	for i := len(l1) - 2; i >= 0; i-- {
		d1[i] = (d1[i+1] * d1[i+1]) % mod
	}

	var ml int = 1
	for i := len(l1) - 1; i >= 0; i-- {
		if l1[i] == '1' {
			ml = (ml * d1[i]) % mod
		}
	}
	return ml % mod
}

// encrypt шифрует текст с использованием открытого ключа
func encrypt(text int, openKey int64, mod int) int {
	return Enc(text, openKey, mod)
}

// decrypt расшифровывает текст с использованием секретного ключа
func decrypt(encryptedText int, secretKey int64, mod int) int {
	return Enc(encryptedText, secretKey, mod)
}

// generateSecretKey находит секретный ключ
func generateSecretKey(openKey int, mod int) int {
	return secret_key(openKey, mod)
}

func main() {
	var choice int
	var openKey, mod, text, secretKey int

	fmt.Println("=== RSA Шифрование ===")
	fmt.Println("1. Зашифровать сообщение")
	fmt.Println("2. Найти секретный ключ")
	fmt.Println("3. Полный цикл (генерация ключа, шифрование, расшифровка)")
	fmt.Print("Выберите действие (1-3): ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Print("Введите открытый ключ и модуль: ")
		fmt.Scanln(&openKey, &mod)
		fmt.Print("Введите число для шифрования: ")
		fmt.Scanln(&text)

		encrypted := encrypt(text, int64(openKey), mod)
		fmt.Printf("Зашифрованное число: %d\n", encrypted)

	case 2:
		fmt.Print("Введите открытый ключ и модуль: ")
		fmt.Scanln(&openKey, &mod)

		secretKey = generateSecretKey(openKey, mod)
		if secretKey == -1 {
			fmt.Println("Ошибка: не удалось найти секретный ключ")
		} else {
			fmt.Printf("Секретный ключ: %d\n", secretKey)
		}

	case 3:
		fmt.Print("Введите открытый ключ и модуль: ")
		fmt.Scanln(&openKey, &mod)
		fmt.Print("Введите число для шифрования: ")
		fmt.Scanln(&text)

		// Генерация секретного ключа
		secretKey = generateSecretKey(openKey, mod)
		if secretKey == -1 {
			fmt.Println("Ошибка: не удалось найти секретный ключ")
			return
		}
		fmt.Printf("Секретный ключ: %d\n", secretKey)

		// Шифрование
		encrypted := encrypt(text, int64(openKey), mod)
		fmt.Printf("Зашифрованное число: %d\n", encrypted)

		// Расшифровка
		decrypted := decrypt(encrypted, int64(secretKey), mod)
		fmt.Printf("Расшифрованное число: %d\n", decrypted)

		// Проверка
		if decrypted == text {
			fmt.Println("✓ Шифрование и расшифровка выполнены успешно!")
		} else {
			fmt.Println("✗ Ошибка: расшифрованное число не совпадает с исходным")
		}

	default:
		fmt.Println("Неверный выбор. Пожалуйста, выберите 1, 2 или 3.")
	}
}
