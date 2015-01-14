package kdtree

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/nstehr/mosaicgen/db"
	"math"
	"sort"
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

type sortByRed []db.Photo
type sortByGreen []db.Photo
type sortByBlue []db.Photo

func (slice sortByRed) Len() int {
	return len(slice)
}

func (slice sortByRed) Less(i, j int) bool {
	return slice[i].AvgColor.R < slice[j].AvgColor.R
}

func (slice sortByRed) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice sortByGreen) Len() int {
	return len(slice)
}

func (slice sortByGreen) Less(i, j int) bool {
	return slice[i].AvgColor.G < slice[j].AvgColor.G
}

func (slice sortByGreen) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice sortByBlue) Len() int {
	return len(slice)
}

func (slice sortByBlue) Less(i, j int) bool {
	return slice[i].AvgColor.B < slice[j].AvgColor.B
}

func (slice sortByBlue) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func NewTree(photos []db.Photo) *Node {
	return createNode(photos, 0)
}

func (node *Node) NearestNeighbour(target colorful.Color) db.Photo {
	var match db.Photo
	bestDistance := 1.0
	nearestNeighborSearch(node, 0, target, &match, &bestDistance)
	return match
}

func nearestNeighborSearch(node *Node, depth int, target colorful.Color, match *db.Photo, bestDistance *float64) {
	if node == nil {
		return
	}
	cc := colorful.Color{float64(node.Photo.AvgColor.R) / 255.0, float64(node.Photo.AvgColor.G) / 255.0, float64(node.Photo.AvgColor.B) / 255.0}

	distance := cc.DistanceCIE94(target)

	if distance < *bestDistance {
		bestDistance = &distance
		*match = *node.Photo
	}

	axis := depth % dimensions

	var ti uint8
	var ni uint8

	switch axis {
	case r:
		ti, _, _ = target.RGB255()
		ni = node.Photo.AvgColor.R
	case g:
		_, ti, _ = target.RGB255()
		ni = node.Photo.AvgColor.G
	case b:
		_, _, ti = target.RGB255()
		ni = node.Photo.AvgColor.B

	}
	leftSearched := true
	if ti < ni {
		nearestNeighborSearch(node.leftNode, depth+1, target, match, bestDistance)
	} else {
		nearestNeighborSearch(node.rightNode, depth+1, target, match, bestDistance)
		leftSearched = false
	}

	if math.Abs(float64(ti)-float64(ni)) < *bestDistance {
		if leftSearched {
			nearestNeighborSearch(node.rightNode, depth+1, target, match, bestDistance)
		} else {
			nearestNeighborSearch(node.leftNode, depth+1, target, match, bestDistance)
		}
	}
}

func createNode(photos []db.Photo, depth int) *Node {

	if len(photos) <= 0 {
		return nil
	}

	axis := depth % dimensions
	switch axis {
	case r:
		sort.Sort(sortByRed(photos))
	case g:
		sort.Sort(sortByGreen(photos))
	case b:
		sort.Sort(sortByBlue(photos))
	}
	median := len(photos) / 2
	left := createNode(photos[0:median], depth+1)
	right := createNode(photos[median+1:], depth+1)

	return &Node{Photo: &photos[median], leftNode: left, rightNode: right}
}
