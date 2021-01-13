package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/axiomhq/hyperminhash"
	pcgr "github.com/dgryski/go-pcgr"
)

func estimateError(got, exp uint64) float64 {
	var delta uint64
	if got > exp {
		delta = got - exp
	} else {
		delta = exp - got
	}
	return float64(delta) / float64(exp)
}

func main() {

	rnd := pcgr.New(time.Now().UnixNano(), 0)
	for k := int64(1000); k < 1000000000; k *= 10 {

		fmt.Println("\n### Max Cardinality", k)
		fmt.Println("\n| Set1 | HMH1 | HLL1 | Set2 | HMH2 | HLL2 | S1 ∪ S2 | HMH1 ∪ HMH2 | HLL1 ∪ HLL2 | S1 ∩ S2 | HMH1 ∩ HMH2 | HLL1+HLL2-(HLL1∪HLL2) |")
		fmt.Println("|---|---|---|---|---|---|---|---|---|---|---|---|")

		for j := 1; j <= 5; j++ {

			size1 := rnd.Int63() % k
			size2 := rnd.Int63() % k
			sk1 := hyperminhash.New()
			sk2 := hyperminhash.New()

			ll1 := hyperloglog.New()
			ll2 := hyperloglog.New()

			maxCol := size1
			if maxCol > size2 {
				maxCol = size2
			}

			cols := rnd.Int63() % maxCol
			intersections := 0
			set := make(map[int]uint8)

			for i := 0; i < int(size1); i++ {
				set[i]++
				k := []byte(strconv.Itoa(i))
				sk1.Add(k)
				ll1.Insert(k)
			}

			for i := int(size1 - cols); i < int(size1-cols+size2); i++ {
				set[i]++
				if set[i] > 1 {
					intersections++
				}
				k := []byte(strconv.Itoa(i))
				sk2.Add(k)
				ll2.Insert(k)
			}

			card1 := sk1.Cardinality()
			card2 := sk2.Cardinality()
			ints1 := sk1.Intersection(sk2)
			m := sk1.Merge(sk2)
			mcard := m.Cardinality()

			var ll1card, ll2card, llmcard, llints uint64
			ll1card = ll1.Estimate()
			ll2card = ll2.Estimate()
			err := ll1.Merge(ll2)
			if err == nil {
				llmcard = ll1.Estimate()
			}
			if llmcard > 0 {
				if llmcard < ll1card+ll2card {
					llints = ll1card + ll2card - llmcard
				}
			} else {
				llmcard, llints = 1, 0
			}

			row := fmt.Sprintf(
				"| %d | %d | %d | %d | %d | %d | %d | %d | %d | **%d** (%f%%) | **%d** (%f%%) | **%d** (%f%%) |",
				size1, card1, ll1card,
				size2, card2, ll2card,
				len(set), mcard, llmcard,
				cols, float64(float64(100*cols)/float64(len(set))),
				ints1, 100*float64(ints1)/float64(mcard),
				llints, float64(llints)/float64(llmcard)*100,
			)

			fmt.Println(row)
		}

		fmt.Println("")

	}

}
