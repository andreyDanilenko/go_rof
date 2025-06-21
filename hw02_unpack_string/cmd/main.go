package main

import (
	"fmt"

	unpack "github.com/andreyDanilenko/go_rof/tree/master/hw02_unpack_string"
)

func main() {
	str1, _ := unpack.Unpack("a0a4bc2d5e")
	str2, _ := unpack.Unpack("aaa0b")
	str3, _ := unpack.Unpack(`qwe\45`)
	str4, _ := unpack.Unpack("qwe\\5")
	str5, _ := unpack.Unpack(`qwe\\\3`)
	str6, _ := unpack.Unpack("aaф0b")
	str7, _ := unpack.Unpack("")

	fmt.Println("str1", str1)
	fmt.Println("str2", str2)
	fmt.Println("str3", str3)
	fmt.Println("str4", str4)
	fmt.Println("str5", str5)
	fmt.Println("str6", str6)
	fmt.Println("str7", str7)
}
