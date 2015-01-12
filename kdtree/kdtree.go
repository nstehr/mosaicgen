package kdtree

import (
	"github.com/nstehr/mosaicgen/db"
	"math"
	"math/rand"
	"sort"
	"time"
)

const (
	dimensions = 3 //hardcode to three dimensions, which is what is needed for me
	r          = 0
	g          = 1
	b          = 2
)

type Node struct {
	Photo     *db.Photo
	leftNode  *Node
	rightNode *Node
}

type sample struct {
	photo   *db.Photo
	origIdx int
}

type sortByRed []sample
type sortByGreen []sample
type sortByBlue []sample

func (slice sortByRed) Len() int {
	return len(slice)
}

func (slice sortByRed) Less(i, j int) bool {
	return slice[i].photo.AvgColor.R < slice[j].photo.AvgColor.R
}

func (slice sortByRed) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice sortByGreen) Len() int {
	return len(slice)
}

func (slice sortByGreen) Less(i, j int) bool {
	return slice[i].photo.AvgColor.G < slice[j].photo.AvgColor.G
}

func (slice sortByGreen) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice sortByBlue) Len() int {
	return len(slice)
}

func (slice sortByBlue) Less(i, j int) bool {
	return slice[i].photo.AvgColor.B < slice[j].photo.AvgColor.B
}

func (slice sortByBlue) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func NewTree(photos []db.Photo) *Node {
	return createNode(photos, 0)
}

func createNode(photos []db.Photo, depth int) *Node {

	if len(photos) <= 0 {
		return nil
	}

	axis := depth % dimensions
	median := getMedian(photos, axis)
	left := createNode(photos[0:median.origIdx], depth+1)
	right := createNode(photos[median.origIdx+1:], depth+1)

	return &Node{Photo: median.photo, leftNode: left, rightNode: right}
}

func getMedian(photos []db.Photo, axis int) sample {
	//estimating the median by taking a fixed number
	//of randomly selected points and taking their median.  In this
	//case the fixed number is 60% the length of the slice
	rand.Seed(time.Now().UTC().UnixNano())
	length := len(photos)
	num := math.Ceil(0.6 * float64(length))
	var samples []sample
	for i := 0; i < int(num); i++ {
		idx := rand.Intn(len(photos))
		photo := photos[idx]
		samples = append(samples, sample{&photo, idx})
	}

	switch axis {
	case r:
		sort.Sort(sortByRed(samples))
	case g:
		sort.Sort(sortByGreen(samples))
	case b:
		sort.Sort(sortByBlue(samples))
	}

	return samples[len(samples)/2]
}
