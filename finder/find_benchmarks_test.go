package finder

import (
	"math"
	"math/rand"
	"testing"
)

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
	var hashMap = make(map[int]int, size)
	var list = make([]int, size)

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
			WithPrefixBuckets(false),
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
			WithPrefixBuckets(true),
		)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f.Find(testRef)
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

	var b = make([]byte, length)
	for i := uint64(0); i < length; i++ {
		b[i] = alnum[rand.Intn(len(alnum))]
	}
	return string(b)
}
