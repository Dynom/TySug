package finder

import (
	"math"
	"math/rand"
	"testing"
)

func BenchmarkRefresh(b *testing.B) {
	initialRefs := []string{
		"zjekuraalf", "irsqktljsw", "lavebqqrbg", "oozfqpnsvl", "zbxcijlhsi", "qqpthspmrk", "woxvyehskj", "wzaiegvxcw",
		"tvdcnodnju", "lnugrvpmtq", "pzpwriltwg", "utadyfzvzk", "tgehqtxmnf", "minfuuzytg", "gicfpqgpot", "dfpqvcrgmv",
		"zbpxhrwdkz", "ajjjefnapl", "gxxaqeqzhu", "njnfdoxrzx", "spbvavebgx", "vhtrygoulv", "ttfawftqek", "pehivzzdsp",
	}

	newRefs := []string{
		"bfzwtgdrny", "wtrpemlhle", "aiirvbdmlv", "tsvjdodfba", "hfzlzvamtz", "vvionkxydg", "jxogoheean", "onisaflgmk",
		"ffzwzeizsg", "ebayntksly", "kfzcguvnid", "vysjvkvplw", "crdkwxxjbr", "dkafiwdmiy", "sushudfnkq", "onblgyuggq",
		"czcrrbouvq", "qukpyfwqxb", "ejasnahxzx", "zwjvdnhrfa", "qqexohqgve", "rquqvvdqjx", "mzjiilpwpz", "mlxyxeozye",
		"fxcckjbsyk", "gpmbrumbqv",
	}

	combinedRefs := append(initialRefs, newRefs...)

	b.Run("only new refs", func(b *testing.B) {
		findr, err := New(initialRefs, WithAlgorithm(NewJaroWinklerDefaults()), WithPrefixBuckets(false))
		if err != nil {
			b.Errorf("Failed setting up test %s", err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			findr.Refresh(newRefs)
		}
	})

	b.Run("extra refs", func(b *testing.B) {
		findr, err := New(initialRefs, WithAlgorithm(NewJaroWinklerDefaults()), WithPrefixBuckets(false))
		if err != nil {
			b.Errorf("Failed setting up test %s", err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			findr.Refresh(combinedRefs)
		}
	})
}

// Preventing the compiler to inline
var ceilA, ceilB int

func BenchmarkCeilOrNoCeil(b *testing.B) {
	inputLen := 64
	threshold := 0.195
	b.Run("No Ceil", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ceilA = int((float64(inputLen) * threshold) + 0.555)
		}
	})

	b.Run("Ceil", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ceilB = int(math.Ceil(float64(inputLen) * threshold))
		}
	})

	if ceilA != ceilB {
		b.Errorf("Implementation failure, a:%d != b:%d", ceilA, ceilB)
	}
}

func BenchmarkSliceOrMap(b *testing.B) {
	// With sets of more than 20 elements, maps become more efficient. (Not including setup costs)
	size := 20
	hashMap := make(map[int]int, size)
	list := make([]int, size)

	for i := size - 1; i > 0; i-- {
		hashMap[i] = i
		list[i] = i
	}

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = hashMap[i]
		}
	})
	b.Run("List", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range list {
				_ = v
			}
		}
	})
}

func BenchmarkFindWithBucket(b *testing.B) {
	refs := generateRefs(1000, 20)
	alg := NewJaroWinkler(.7, 4)

	testRef := generateRef(20)
	b.ReportAllocs()
	b.Run("find with bucket", func(b *testing.B) {
		f, _ := New(refs,
			WithAlgorithm(alg),
			WithLengthTolerance(0),
			WithPrefixBuckets(true),
		)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f.Find(testRef)
		}
	})

	b.Run("find without bucket", func(b *testing.B) {
		f, _ := New(refs,
			WithAlgorithm(alg),
			WithLengthTolerance(0),
			WithPrefixBuckets(false),
		)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f.Find(testRef)
		}
	})
}

const numToAllocate = 1024 * 1024

var (
	refsCopySrc   = make([]string, numToAllocate)
	refsAppendSrc = make([]string, numToAllocate)
	refsCopyDst   []string
	refsAppendDst []string
)

func BenchmarkCopyOrAppend(b *testing.B) {
	refsCopySrc[0] = "a"
	refsCopySrc[len(refsCopySrc)-1] = "z"
	refsAppendSrc[0] = "a"
	refsAppendSrc[len(refsAppendSrc)-1] = "z"

	b.Run("equal size copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refsCopyDst = make([]string, numToAllocate)
			copy(refsCopyDst, refsCopySrc)
		}

		if first, last := refsCopyDst[0], refsCopyDst[len(refsCopyDst)-1]; first != "a" || last != "z" {
			b.Errorf("result length: %d first and last index value doesn't match: %s - %s", len(refsCopyDst), first, last)
		}
	})

	b.Run("equal size append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refsAppendDst = append(refsAppendSrc[:0:0], refsAppendSrc...)
		}

		if first, last := refsAppendDst[0], refsAppendDst[len(refsAppendDst)-1]; first != "a" || last != "z" {
			b.Errorf("result length: %d first and last index value doesn't match: %s - %s", len(refsAppendDst), first, last)
		}
	})

	// "dst smaller copy" can't work, since the result won't contain all items or requires logic which makes the
	// implementation slower then append

	b.Run("dst smaller append", func(b *testing.B) {
		refsAppendDst = make([]string, numToAllocate/2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			refsAppendDst = append(refsAppendSrc[:0:0], refsAppendSrc...)
		}

		if first, last := refsAppendDst[0], refsAppendDst[len(refsAppendDst)-1]; first != "a" || last != "z" {
			b.Errorf("result length: %d first and last index value doesn't match: %s - %s", len(refsAppendDst), first, last)
		}
	})

	b.Run("dst larger copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refsCopyDst = make([]string, numToAllocate*2)
			copy(refsCopyDst, refsCopySrc)

			// Necessary branch to shrink to the right length
			if len(refsCopyDst) > len(refsCopySrc) {
				refsCopyDst = refsCopyDst[:len(refsCopySrc)]
			}
		}

		if first, last := refsCopyDst[0], refsCopyDst[len(refsCopyDst)-1]; first != "a" || last != "z" {
			b.Errorf("result length: %d first and last index value doesn't match: %s - %s", len(refsCopyDst), first, last)
		}
	})

	b.Run("dst larger append", func(b *testing.B) {
		refsAppendDst = make([]string, numToAllocate*2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			refsAppendDst = append(refsAppendSrc[:0:0], refsAppendSrc...)
		}

		if first, last := refsAppendDst[0], refsAppendDst[len(refsAppendDst)-1]; first != "a" || last != "z" {
			b.Errorf("result length: %d first and last index value doesn't match: %s - %s", len(refsAppendDst), first, last)
		}
	})
}

func generateRefs(refNum, length uint64) []string {
	refs := make([]string, refNum)
	for i := uint64(0); i < refNum; i++ {
		refs[i] = generateRef(length)
	}

	return refs
}

func generateRef(length uint64) string {
	const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := uint64(0); i < length; i++ {
		b[i] = alnum[rand.Intn(len(alnum))]
	}
	return string(b)
}
