package main
// 三个类型确定常量。
const n int = 1 << 64           // error: 溢出int
const r rune = 'a' + 0x7FFFFFFF // error: 溢出rune
const x float64 = 2e+308        // error: 溢出float64
func main() {}