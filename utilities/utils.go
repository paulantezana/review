package utilities

import "fmt"

// Convert number to A,B,C… AA,AB,… AAA…, in base 26
func ConvertNumberToCharScheme(n uint) string  {
    baseChar := []rune("A")[0]
    letters := ""
    for ok := true; ok; ok = n > 0 {
        n -= 1
        letters = fmt.Sprintf("%c", uint(baseChar) + (n % 26)) + letters
        n = (n / 26) >> 0
    }
    return letters
}