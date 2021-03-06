package samparser

import (
	"bufio"
	"github.com/Hanbin/AberrantSplice/Internal/genodatastruct"
	"log"
	"os"
	"strconv"
	"strings"
	//"sync"
)

//ParseSam filter reads and parse them into SamRec struct
//and generating a channel of iterator
func ParseSam(sam string, filter func(genodatastruct.SamRec) bool) <-chan genodatastruct.SamRec {
	out := make(chan genodatastruct.SamRec, 100)
	go func() {
		samF, err := os.Open(sam)
		if err != nil {
			log.Fatal(err)
		}
		defer samF.Close()

		scanner := bufio.NewScanner(samF)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "@") {
				continue
			}
			fields := strings.SplitN(line, "\t", 7)
			flag, _ := strconv.Atoi(fields[1])
			if flag == 4 { //unmapped
				continue
			}
			mapq, _ := strconv.Atoi(fields[4])
			pos, _ := strconv.Atoi(fields[3])
			temp := genodatastruct.SamRec{
				Flag:       int64(flag),
				MAPQ:       mapq,
				Pos:        pos,
				Chromosome: genodatastruct.ChroSym(fields[2]),
				CIGAR:      fields[5],
			}
			if filter(temp) {
				out <- temp
			}
		}
		close(out)
	}()
	return out
}

//concurrent does not imporove equal process time processes of I/O and split
////Fan in sam lines
//func SamReader(sam string) <-chan string {
//	out := make(chan string, 100)
//	go func() {
//		samF, err := os.Open(sam)
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer samF.Close()
//		scanner := bufio.NewScanner(samF)
//		for scanner.Scan() {
//			line := scanner.Text()
//			if strings.HasPrefix(line, "@") {
//				continue
//			}
//			out <- line
//		}
//		close(out)
//	}()
//	return out
//}
//
////Unit job of spliting sam line and generate Sam rec structure
//func MakeSamRec(in <-chan string, out chan genodatastruct.SamRec, filter func(genodatastruct.SamRec) bool) {
//	for line := range in {
//		fields := strings.SplitN(line, "\t", 7)
//		flag, _ := strconv.Atoi(fields[1])
//		if flag == 4 { //unmapped
//			continue
//		}
//		mapq, _ := strconv.Atoi(fields[4])
//		pos, _ := strconv.Atoi(fields[3])
//		temp := genodatastruct.SamRec{
//			Flag:       int64(flag),
//			MAPQ:       mapq,
//			Pos:        pos,
//			Chromosome: genodatastruct.ChroSym(fields[2]),
//			CIGAR:      fields[5],
//		}
//		if filter(temp) {
//			out <- temp
//		}
//	}
//	close(out)
//}
//
////ParseSam filter reads and parse them into SamRec struct
////and generating a channel of iterator
//func ParseSamConcur(sam string, filter func(genodatastruct.SamRec) bool) <-chan genodatastruct.SamRec {
//	readline := SamReader(sam)
//	spliter := []chan genodatastruct.SamRec{}
//	n := 1
//	for i := 0; i < n; i++ {
//		o := make(chan genodatastruct.SamRec, 10)
//		go MakeSamRec(readline, o, filter)
//		spliter = append(spliter, o)
//	}
//	//fanout
//	var wg sync.WaitGroup
//	out := make(chan genodatastruct.SamRec)
//	output := func(c chan genodatastruct.SamRec) {
//		for m := range c {
//			out <- m
//		}
//		wg.Done()
//	}
//	wg.Add(n)
//	for _, worker := range spliter {
//		go output(worker)
//	}
//	go func() {
//		wg.Wait()
//		close(out)
//	}()
//	return out
//}
//
