package genodatastruct

import (
	"sort"
)

//Allgene is a preset value for choosing behaviour of selecting all gene to record
var Allgene = []string{"all"}

//Gene is a struct type for Gtf features
type Gene struct { //a gene struct records name, genomic locus and alternative transcripts of it
	GeneName           string
	Chromosome, Strand string
	Coordinate         Coor
	Transcripts        []*Transcript
	Attributes         map[string]string //split by ";" and process to key/val pair
}

//Transcript is a struct type record the exons of the transcripts
type Transcript struct { //name, genomic locus, exon
	TranscriptName     string
	Chromosome, Strand string
	Coordinate         Coor
	Exons              []Coor //start and end locations of an exon
	Attributes         map[string]string
}

//Coordinate is the (start, end) coordinates of a genomic feature
type Coor struct {
	Start, End int
}

//###########################
//METHODS
//###########################
func MergeRegions(regions []Coor) []Coor {
	//sort regions by start coor before merging
	sort.Slice(regions, func(i, j int) bool { return regions[i].Start <= regions[j].Start })
	//merge overlap region into a larger region
	merged := []Coor{regions[0]} //initiate
	for _, reg := range regions {
		rightmost := &merged[len(merged)-1]
		if reg.Start <= rightmost.End && reg.End >= rightmost.End {
			// overlap and need to extend
			rightmost.End = reg.End
		} else if reg.Start > rightmost.End {
			merged = append(merged, reg)
		}
	}
	return merged
}

//Interval region take in a SORTED region list and return the intervals
//demarcated by the list
func IntervalRegions(regions []Coor) []Coor {
	interval := make([]Coor, len(regions)-1)
	for i := 0; i < len(interval); i++ {
		interval[i] = Coor{regions[i].End + 1, regions[i+1].Start - 1}
	}
	return interval
}

//MergeExons merges all the exons of a gene to
//produce a joint set of exons of a gene locus
func (g Gene) MergeExons() []Coor {
	exonCollect := []Coor{}
	for _, transcript := range g.Transcripts {
		exonCollect = append(exonCollect, transcript.Exons...)
	}
	merged := MergeRegions(exonCollect)
	return merged
}

//IntervalofExons takes the interval of merged exons
//they are typically high confident introns
func (g Gene) IntervalOfExons() []Coor {
	merged := g.MergeExons()
	return IntervalRegions(merged)
}
