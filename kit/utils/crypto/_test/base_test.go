package _test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Base64Decode(t *testing.T) {
	str := "eyJoZWFkIjp7InZlcnNpb24iOiIxLjAuMC4xIiwiY2xpZW50X2lkIjoiQUs3NzA2MEUyMTYwMDAzOSIsInNlc3NfaWQiOiIyMTM3LTI3MzkiLCJtc2dfc2VxIjoyMywiY21kIjoxNiwiY21kX3N0YXR1cyI6MSwic2lnbl9jb2RlIjowLCJzZWVkIjoibnVsbCIsInRpbWVzdGFtcCI6MTcyMTM3MjczMn0sInBheWxvYWQiOiJ7XCJkYXRhXCI6XCIvd2xRQUFNQkR3TURvUUUwb2dnRUFBQUFBQUFBQUtNSUJBQUFBQUFBQUFDa0NBUUFBQUFBQUFBQXBRZ0VBQUFBQUFBQUFLWUlCQUFBQUFBQUFBQ25DQVFBQUFBQUFBQUEvZ1VEUEJDYVpoWT1cIixcInNuXCI6XCJBSzc3MDYwRTIxNjAwMDM5XCIsXCJwblwiOlwiQTIzNDVcIn0ifQ=="

	// base64 转换 string
	ts, _ := base64.StdEncoding.DecodeString(str)

	t.Log(str)
	t.Log(string(ts))
}

// string 转换 []byte 时，长度与字节数相等。（中文一个字符是两个字节）
// 1. 先将每一个字符转换为对应ASCII码的十进制值。 [72 101 108 108 111 44 32 87 111 114 108 100 33]
func Test_StringToByte(t *testing.T) {
	str := "Hello, World!"

	// string 转换 []byte
	st := []byte(str)

	// []byte 转换 string
	ts := string(st)

	t.Log(str)
	t.Log(st) // [72 101 108 108 111 44 32 87 111 114 108 100 33]
	t.Log(ts)

	assert.Equal(t, str, ts)
}

// string 转换 hex 时，长度是原字符串的两倍。
// 1. 先将每一个字符转换为对应ASCII码的十进制值。 [72 101 108 108 111 44 32 87 111 114 108 100 33]
// 2. 然后将每一个字节的十进制转换为十六进制。 72(十进制) -> 48(十六进制) , 101(十进制) -> 65(十六进制) ...
func Test_StringToHex(t *testing.T) {
	str := "Hello, World!"

	// string 转换 hex
	st := hex.EncodeToString([]byte(str))

	// hex 转换 string
	ts, _ := hex.DecodeString(st)

	t.Log(str)
	t.Log(st) // 48656c6c6f2c20576f726c6421
	t.Log(ts)

	assert.Equal(t, str, string(ts))
}

// string 转换 base64时，长度是原字符串的 4/3 倍。
// 1. 先将每一个字符转换为对应ASCII码的十进制值。  [72 101 108 108 111 44 32 87 111 114 108 100 33]
// 2. 将每一个字节的十进制值转换为8位二进制。 72(十进制) -> 01001000(二进制) , 101(十进制) -> 01100101(二进制) ...
// 3. 将8位二进制拆分为6位一组，不足6位的在后面补0。 010010 000110 010101
// 4. 将每个 6 位的二进制数转换为对应的Base64 字符。 010010 -> S
// 5. 因为 "Hello, World!" 的长度是 13，不是 3 的倍数，所以添加两个填充字符 "="。
func Test_StringToBase64(t *testing.T) {
	str := "Hello, World!"

	// string 转换 base64
	st := base64.StdEncoding.EncodeToString([]byte(str))

	// base64 转换 string
	ts, _ := base64.StdEncoding.DecodeString(st)

	t.Log(str)
	t.Log(st) // SGVsbG8sIFdvcmxkIQ==
	t.Log(ts)

	assert.Equal(t, str, string(ts))
}

func Test_StringTo(t *testing.T) {
	bytes := []byte{'H', 'e', 'l', 'l', 'o', ',', ' ', 'W', 'o', 'r', 'l', 'd', '!'}
	str := string(bytes)
	fmt.Println(str)
}

/**
	ASCII Table
	Dec   Char     Dec   Char     Dec   Char     Dec   Char
	---------     ---------     ---------     ----------
	0   NUL       32   SPACE     64   @         96   `
    1   SOH       33   !         65   A         97   a
    2   STX       34   "         66   B         98   b
    3   ETX       35   #         67   C         99   c
    4   EOT       36   $         68   D        100   d
    5   ENQ       37   %         69   E        101   e
    6   ACK       38   &         70   F        102   f
    7   BEL       39   '         71   G        103   g
    8   BS        40   (         72   H        104   h
    9   TAB       41   )         73   I        105   i
   10   LF        42   *         74   J        106   j
   11   VT        43   +         75   K        107   k
   12   FF        44   ,         76   L        108   l
   13   CR        45   -         77   M        109   m
   14   SO        46   .         78   N        110   n
   15   SI        47   /         79   O        111   o
   16   DLE       48   0         80   P        112   p
   17   DC1       49   1         81   Q        113   q
   18   DC2       50   2         82   R        114   r
   19   DC3       51   3         83   S        115   s
   20   DC4       52   4         84   T        116   t
   21   NAK       53   5         85   U        117   u
   22   SYN       54   6         86   V        118   v
   23   ETB       55   7         87   W        119   w
   24   CAN       56   8         88   X        120   x
   25   EM        57   9         89   Y        121   y
   26   SUB       58   :         90   Z        122   z
   27   ESC       59   ;         91   [        123   {
   28   FS        60   <         92   \        124   |
   29   GS        61   =         93   ]        125   }
   30   RS        62   >         94   ^        126   ~
   31   US        63   ?         95   _        127   DEL

*/
