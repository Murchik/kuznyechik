package crypt

import (
	"crypto/cipher"
	"strconv"
	"unsafe"
)

var CipherInitialized = false

const BlockSize = 16

var Pi_table = [256]uint8{

	0xFC, 0xEE, 0xDD, 0x11, 0xCF, 0x6E, 0x31, 0x16,
	0xFB, 0xC4, 0xFA, 0xDA, 0x23, 0xC5, 0x04, 0x4D,
	0xE9, 0x77, 0xF0, 0xDB, 0x93, 0x2E, 0x99, 0xBA,
	0x17, 0x36, 0xF1, 0xBB, 0x14, 0xCD, 0x5F, 0xC1,
	0xF9, 0x18, 0x65, 0x5A, 0xE2, 0x5C, 0xEF, 0x21,
	0x81, 0x1C, 0x3C, 0x42, 0x8B, 0x01, 0x8E, 0x4F,
	0x05, 0x84, 0x02, 0xAE, 0xE3, 0x6A, 0x8F, 0xA0,
	0x06, 0x0B, 0xED, 0x98, 0x7F, 0xD4, 0xD3, 0x1F,
	0xEB, 0x34, 0x2C, 0x51, 0xEA, 0xC8, 0x48, 0xAB,
	0xF2, 0x2A, 0x68, 0xA2, 0xFD, 0x3A, 0xCE, 0xCC,
	0xB5, 0x70, 0x0E, 0x56, 0x08, 0x0C, 0x76, 0x12,
	0xBF, 0x72, 0x13, 0x47, 0x9C, 0xB7, 0x5D, 0x87,
	0x15, 0xA1, 0x96, 0x29, 0x10, 0x7B, 0x9A, 0xC7,
	0xF3, 0x91, 0x78, 0x6F, 0x9D, 0x9E, 0xB2, 0xB1,
	0x32, 0x75, 0x19, 0x3D, 0xFF, 0x35, 0x8A, 0x7E,
	0x6D, 0x54, 0xC6, 0x80, 0xC3, 0xBD, 0x0D, 0x57,
	0xDF, 0xF5, 0x24, 0xA9, 0x3E, 0xA8, 0x43, 0xC9,
	0xD7, 0x79, 0xD6, 0xF6, 0x7C, 0x22, 0xB9, 0x03,
	0xE0, 0x0F, 0xEC, 0xDE, 0x7A, 0x94, 0xB0, 0xBC,
	0xDC, 0xE8, 0x28, 0x50, 0x4E, 0x33, 0x0A, 0x4A,
	0xA7, 0x97, 0x60, 0x73, 0x1E, 0x00, 0x62, 0x44,
	0x1A, 0xB8, 0x38, 0x82, 0x64, 0x9F, 0x26, 0x41,
	0xAD, 0x45, 0x46, 0x92, 0x27, 0x5E, 0x55, 0x2F,
	0x8C, 0xA3, 0xA5, 0x7D, 0x69, 0xD5, 0x95, 0x3B,
	0x07, 0x58, 0xB3, 0x40, 0x86, 0xAC, 0x1D, 0xF7,
	0x30, 0x37, 0x6B, 0xE4, 0x88, 0xD9, 0xE7, 0x89,
	0xE1, 0x1B, 0x83, 0x49, 0x4C, 0x3F, 0xF8, 0xFE,
	0x8D, 0x53, 0xAA, 0x90, 0xCA, 0xD8, 0x85, 0x61,
	0x20, 0x71, 0x67, 0xA4, 0x2D, 0x2B, 0x09, 0x5B,
	0xCB, 0x9B, 0x25, 0xD0, 0xBE, 0xE5, 0x6C, 0x52,
	0x59, 0xA6, 0x74, 0xD2, 0xE6, 0xF4, 0xB4, 0xC0,
	0xD1, 0x66, 0xAF, 0xC2, 0x39, 0x4B, 0x63, 0xB6,
}

var Pi_inverse_table = [256]uint8{

	0xA5, 0x2D, 0x32, 0x8F, 0x0E, 0x30, 0x38, 0xC0,
	0x54, 0xE6, 0x9E, 0x39, 0x55, 0x7E, 0x52, 0x91,
	0x64, 0x03, 0x57, 0x5A, 0x1C, 0x60, 0x07, 0x18,
	0x21, 0x72, 0xA8, 0xD1, 0x29, 0xC6, 0xA4, 0x3F,
	0xE0, 0x27, 0x8D, 0x0C, 0x82, 0xEA, 0xAE, 0xB4,
	0x9A, 0x63, 0x49, 0xE5, 0x42, 0xE4, 0x15, 0xB7,
	0xC8, 0x06, 0x70, 0x9D, 0x41, 0x75, 0x19, 0xC9,
	0xAA, 0xFC, 0x4D, 0xBF, 0x2A, 0x73, 0x84, 0xD5,
	0xC3, 0xAF, 0x2B, 0x86, 0xA7, 0xB1, 0xB2, 0x5B,
	0x46, 0xD3, 0x9F, 0xFD, 0xD4, 0x0F, 0x9C, 0x2F,
	0x9B, 0x43, 0xEF, 0xD9, 0x79, 0xB6, 0x53, 0x7F,
	0xC1, 0xF0, 0x23, 0xE7, 0x25, 0x5E, 0xB5, 0x1E,
	0xA2, 0xDF, 0xA6, 0xFE, 0xAC, 0x22, 0xF9, 0xE2,
	0x4A, 0xBC, 0x35, 0xCA, 0xEE, 0x78, 0x05, 0x6B,
	0x51, 0xE1, 0x59, 0xA3, 0xF2, 0x71, 0x56, 0x11,
	0x6A, 0x89, 0x94, 0x65, 0x8C, 0xBB, 0x77, 0x3C,
	0x7B, 0x28, 0xAB, 0xD2, 0x31, 0xDE, 0xC4, 0x5F,
	0xCC, 0xCF, 0x76, 0x2C, 0xB8, 0xD8, 0x2E, 0x36,
	0xDB, 0x69, 0xB3, 0x14, 0x95, 0xBE, 0x62, 0xA1,
	0x3B, 0x16, 0x66, 0xE9, 0x5C, 0x6C, 0x6D, 0xAD,
	0x37, 0x61, 0x4B, 0xB9, 0xE3, 0xBA, 0xF1, 0xA0,
	0x85, 0x83, 0xDA, 0x47, 0xC5, 0xB0, 0x33, 0xFA,
	0x96, 0x6F, 0x6E, 0xC2, 0xF6, 0x50, 0xFF, 0x5D,
	0xA9, 0x8E, 0x17, 0x1B, 0x97, 0x7D, 0xEC, 0x58,
	0xF7, 0x1F, 0xFB, 0x7C, 0x09, 0x0D, 0x7A, 0x67,
	0x45, 0x87, 0xDC, 0xE8, 0x4F, 0x1D, 0x4E, 0x04,
	0xEB, 0xF8, 0xF3, 0x3E, 0x3D, 0xBD, 0x8A, 0x88,
	0xDD, 0xCD, 0x0B, 0x13, 0x98, 0x02, 0x93, 0x80,
	0x90, 0xD0, 0x24, 0x34, 0xCB, 0xED, 0xF4, 0xCE,
	0x99, 0x10, 0x44, 0x40, 0x92, 0x3A, 0x01, 0x26,
	0x12, 0x1A, 0x48, 0x68, 0xF5, 0x81, 0x8B, 0xC7,
	0xD6, 0x20, 0x0A, 0x08, 0x00, 0x4C, 0xD7, 0x74,
}

var L_vector = [16]uint8{0x94, 0x20, 0x85, 0x10, 0xC2, 0xC0, 0x01, 0xFB, 0x01, 0xC0, 0xC2, 0x10, 0x85, 0x20, 0x94, 0x01}

var LS_enc_lookup [16][256][16]uint8

var L_inv_lookup [16][256][16]uint8

var SL_dec_lookup [16][256][16]uint8

func GF2_mul(x, y uint8) uint8 {
	var z uint8

	z = 0
	for y != 0 {
		if y&1 == 1 {
			z = z ^ x
		}
		if x&0x80 != 0 {
			x = (x << 1) ^ 0xC3
		} else {
			x = x << 1
		}
		y = y >> 1
	}

	return z
}

func L(block [16]uint8) [16]uint8 {

	var i, j int
	var x uint8

	for j = 0; j < 16; j++ {
		x = block[15]
		for i = 14; i >= 0; i-- {
			block[i+1] = block[i]
			x = x ^ GF2_mul(block[i], L_vector[i])
		}
		block[0] = x
	}

	return block
}

func L_inv(block [16]uint8) [16]uint8 {
	var i, j int
	var x uint8

	for j = 0; j < 16; j++ {
		x = block[0]
		for i = 0; i < 15; i++ {
			block[i] = block[i+1]
			x = x ^ GF2_mul(block[i], L_vector[i])
		}
		block[15] = x
	}

	return block
}

func StretchKey(key [32]uint8) [10][16]uint8 {
	var i, k int
	var C, x, y, z [16]uint8
	var rkeys [10][16]uint8

	for i = 0; i < 16; i++ {
		x[i] = key[i]
		y[i] = key[i+16]
	}

	rkeys[0] = x
	rkeys[1] = y

	for i = 1; i <= 32; i++ {

		for k = range C {
			C[k] = 0
		}
		C[15] = uint8(i)
		C = L(C)

		for k = range z {
			z[k] = Pi_table[(x[k] ^ C[k])]
		}
		z = L(z)
		for k = range z {
			z[k] = z[k] ^ y[k]
		}

		y = x
		x = z

		if i%8 == 0 {
			rkeys[(i >> 2)] = x
			rkeys[(i>>2)+1] = y
		}
	}
	return rkeys
}

func GetDecryptRoundKeys(rkeys [10][16]uint8) [10][16]uint8 {
	var rkeys_L [10][16]uint8

	for k := 1; k < 10; k++ {
		rkeys_L[k] = L_inv(rkeys[k])
	}
	rkeys_L[0] = rkeys[0]
	return rkeys_L
}

func Encrypt(key [32]uint8, block [16]uint8) [16]uint8 {
	var ct [16]uint8
	var rkeys [10][16]uint8

	if !CipherInitialized {
		InitCipher()
	}

	rkeys = StretchKey(key)

	ct = Encrypt_K(rkeys, block)

	return ct
}

func Encrypt_K(rkeys [10][16]uint8, block [16]uint8) [16]uint8 {
	var i, j, k int
	var ct, r [16]uint8

	ct = block
	for i = 0; i < 9; i++ {
		*(*uint64)(unsafe.Pointer(&ct[0])) = *(*uint64)(unsafe.Pointer(&ct[0])) ^ *(*uint64)(unsafe.Pointer(&rkeys[i][0]))
		*(*uint64)(unsafe.Pointer(&ct[8])) = *(*uint64)(unsafe.Pointer(&ct[8])) ^ *(*uint64)(unsafe.Pointer(&rkeys[i][8]))

		for k = range r {
			r[k] = LS_enc_lookup[0][ct[0]][k]
		}

		for j = 1; j <= 15; j++ {
			*(*uint64)(unsafe.Pointer(&r[0])) = *(*uint64)(unsafe.Pointer(&r[0])) ^ *(*uint64)(unsafe.Pointer(&LS_enc_lookup[j][ct[j]][0]))
			*(*uint64)(unsafe.Pointer(&r[8])) = *(*uint64)(unsafe.Pointer(&r[8])) ^ *(*uint64)(unsafe.Pointer(&LS_enc_lookup[j][ct[j]][8]))

		}
		ct = r
	}

	*(*uint64)(unsafe.Pointer(&ct[0])) = *(*uint64)(unsafe.Pointer(&ct[0])) ^ *(*uint64)(unsafe.Pointer(&rkeys[9][0]))
	*(*uint64)(unsafe.Pointer(&ct[8])) = *(*uint64)(unsafe.Pointer(&ct[8])) ^ *(*uint64)(unsafe.Pointer(&rkeys[9][8]))

	return ct

}

func Decrypt_L(key [32]uint8, block [16]uint8) [16]uint8 {
	var i, j, k int
	var pt, r [16]uint8
	var rkeys [10][16]uint8

	rkeys = StretchKey(key)
	pt = block

	for i = 9; i > 0; i-- {
		for k = range pt {
			pt[k] = pt[k] ^ rkeys[i][k]
		}

		for k = range r {
			r[k] = L_inv_lookup[0][pt[0]][k]
		}
		for j = 1; j <= 15; j++ {
			for k = range r {
				r[k] = r[k] ^ L_inv_lookup[j][pt[j]][k]
			}
		}
		pt = r
		for k = range pt {
			pt[k] = Pi_inverse_table[pt[k]]
		}
	}
	for k = range pt {
		pt[k] = pt[k] ^ rkeys[0][k]
	}
	return pt
}

func Decrypt(key [32]uint8, block [16]uint8) [16]uint8 {
	var rkeys [10][16]uint8

	if !CipherInitialized {
		InitCipher()
	}

	rkeys = GetDecryptRoundKeys(StretchKey(key))
	pt := Decrypt_K(rkeys, block)
	return pt
}

func Decrypt_K(rkeys [10][16]uint8, block [16]uint8) [16]uint8 {
	var i, j, k int
	var pt, r [16]uint8

	pt = block
	for k = range r {
		r[k] = L_inv_lookup[0][pt[0]][k]
	}
	for j = 1; j <= 15; j++ {
		*(*uint64)(unsafe.Pointer(&r[0])) = *(*uint64)(unsafe.Pointer(&r[0])) ^ *(*uint64)(unsafe.Pointer(&L_inv_lookup[j][pt[j]][0]))
		*(*uint64)(unsafe.Pointer(&r[8])) = *(*uint64)(unsafe.Pointer(&r[8])) ^ *(*uint64)(unsafe.Pointer(&L_inv_lookup[j][pt[j]][8]))
	}
	pt = r

	for i = 9; i > 1; i-- {
		*(*uint64)(unsafe.Pointer(&pt[0])) = *(*uint64)(unsafe.Pointer(&pt[0])) ^ *(*uint64)(unsafe.Pointer(&rkeys[i][0]))
		*(*uint64)(unsafe.Pointer(&pt[8])) = *(*uint64)(unsafe.Pointer(&pt[8])) ^ *(*uint64)(unsafe.Pointer(&rkeys[i][8]))

		for k = range r {
			r[k] = SL_dec_lookup[0][pt[0]][k]
		}
		for j = 1; j <= 15; j++ {
			*(*uint64)(unsafe.Pointer(&r[0])) = *(*uint64)(unsafe.Pointer(&r[0])) ^ *(*uint64)(unsafe.Pointer(&SL_dec_lookup[j][pt[j]][0]))
			*(*uint64)(unsafe.Pointer(&r[8])) = *(*uint64)(unsafe.Pointer(&r[8])) ^ *(*uint64)(unsafe.Pointer(&SL_dec_lookup[j][pt[j]][8]))
		}
		pt = r
	}

	//for k = range pt {
	//		pt[k] = pt[k] ^ rkeys[1][k]			// XOR with K_2
	//		pt[k] = Pi_inverse_table[pt[k]]	// Inverse Pi
	//		pt[k] = pt[k] ^ rkeys[0][k]			// XOR with K_1
	//}

	*(*uint64)(unsafe.Pointer(&pt[0])) = *(*uint64)(unsafe.Pointer(&pt[0])) ^ *(*uint64)(unsafe.Pointer(&rkeys[1][0]))
	*(*uint64)(unsafe.Pointer(&pt[8])) = *(*uint64)(unsafe.Pointer(&pt[8])) ^ *(*uint64)(unsafe.Pointer(&rkeys[1][8]))

	for k = range pt {
		pt[k] = Pi_inverse_table[pt[k]]
	}

	*(*uint64)(unsafe.Pointer(&pt[0])) = *(*uint64)(unsafe.Pointer(&pt[0])) ^ *(*uint64)(unsafe.Pointer(&rkeys[0][0]))
	*(*uint64)(unsafe.Pointer(&pt[8])) = *(*uint64)(unsafe.Pointer(&pt[8])) ^ *(*uint64)(unsafe.Pointer(&rkeys[0][8]))

	return pt
}

func InitCipher() {
	var i, j, k int
	var x [16]uint8

	if CipherInitialized {
		return
	}

	for i = 0; i < 16; i++ {
		for j = 0; j < 256; j++ {

			for k = range x {
				x[k] = 0
			}
			x[i] = Pi_table[j]
			x = L(x)
			LS_enc_lookup[i][j] = x

			for k = range x {
				x[k] = 0
			}
			x[i] = uint8(j)
			x = L_inv(x)
			L_inv_lookup[i][j] = x

			for k = range x {
				x[k] = 0
			}
			x[i] = Pi_inverse_table[j]
			x = L_inv(x)
			SL_dec_lookup[i][j] = x

		}
	}
	CipherInitialized = true
	return
}

type kuznecCipher struct {
	enc_keys [10][16]uint8
	dec_keys [10][16]uint8
}

type KeySizeError int

func (k KeySizeError) Error() string {
	return "Kuznyechik cipher: invalid key size! Must be 32 bytes - got: " + strconv.Itoa(int(k))
}

func NewCipher(key []byte) (cipher.Block, error) {
	var t_key [32]uint8

	if len(key) != 32 {
		return nil, KeySizeError(len(key))
	}

	c := *(new(kuznecCipher))
	copy(t_key[:], key[:32])
	c.enc_keys = StretchKey(t_key)
	c.dec_keys = GetDecryptRoundKeys(c.enc_keys)
	if !CipherInitialized {
		InitCipher()
	}

	return &c, nil
}

func (c *kuznecCipher) BlockSize() int {
	return BlockSize
}

func (c *kuznecCipher) Encrypt(dst, src []byte) {
	var ct_block [16]uint8

	if len(src) < BlockSize {
		panic("Kuznyechik cipher: input length less than full block!")
	}
	if len(dst) < BlockSize {
		panic("Kuznyechik cipher: output length less than full block!")
	}
	copy(ct_block[:], src[:16])
	ct_block = Encrypt_K(c.enc_keys, ct_block)
	copy(dst, ct_block[:])

}

func (c *kuznecCipher) Decrypt(dst, src []byte) {
	var pt_block [16]uint8

	if len(src) < BlockSize {
		panic("Kuznyechik cipher: input length less than full block!")
	}
	if len(dst) < BlockSize {
		panic("Kuznyechik cipher: output length less than full block!")
	}

	copy(pt_block[:], src[:16])
	pt_block = Decrypt_K(c.dec_keys, pt_block)
	copy(dst, pt_block[:])

}
