package main

import (
	"fmt"

	kuznec "github.com/Murchik/kuznyechik/crypt"
)

// GOST R 34.12-2015: https://tc26.ru/standard/gost/GOST_R_3412-2015.pdf

func main() {
	var test_K = [32]uint8{0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

	var test_PT = [16]uint8{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88}
	var reference_CT = [16]uint8{0x7F, 0x67, 0x9D, 0x90, 0xBE, 0xBC, 0x24, 0x30, 0x5A, 0x46, 0x8D, 0x42, 0xB9, 0xD4, 0xED, 0xCD}

	kuznec.InitCipher()

	fmt.Printf("Message: %X\n", test_PT)

	test_CT := kuznec.Encrypt(test_K, test_PT)

	fmt.Printf("Cipher text: %X - ", test_CT)
	if test_CT != reference_CT {
		fmt.Printf("FAILED! [Not equal to reference cipher text!]\n")
	} else {
		fmt.Printf("OK\n")
	}

	test_2PT := kuznec.Decrypt(test_K, test_CT)

	fmt.Printf("Message decrypted: %X - ", test_2PT)
	if test_2PT != test_PT {
		fmt.Printf("FAILED! [PT != D(E(PT,K),K)]\n")
	} else {
		fmt.Printf("OK\n")
	}
}
