package dbscan

import (
	"fmt"
	"time"
	"log"
	"math"
	"math/rand"
	"testing"
)

type SimpleClusterable struct {
	position float64
}

func (s SimpleClusterable) Distance(c interface{}) float64 {
	distance := math.Abs(c.(SimpleClusterable).position - s.position)
	return distance
}

func (s SimpleClusterable) GetID() string {
	return fmt.Sprint(s.position)
}

func TestPutAll(t *testing.T) {
	testMap := make(map[string]Clusterable)
	clusterList := []Clusterable{
		SimpleClusterable{10},
		SimpleClusterable{12},
	}
	putAll(testMap, clusterList)
	mapSize := len(testMap)
	if mapSize != 2 {
		t.Errorf("Map does not contain expected size 2 but was %d", mapSize)
	}
}
 
//Test find neighbour function
func TestFindUnclusteredNeighbours(t *testing.T) {
	log.Println("Executing TestFindUnclusteredNeighbours")
	clusterList := []Clusterable{
		SimpleClusterable{0},
		SimpleClusterable{1},
		SimpleClusterable{-1},
		SimpleClusterable{1.5},
		SimpleClusterable{-0.5},
	}
	visited := make(map[string]bool)
	eps := 1.0
	neighbours := findUnclusteredNeighbours(clusterList[0], clusterList, visited, eps)

	assertEquals(t, 3, len(neighbours))
}

func TestMerge(t *testing.T) {
	log.Println("Executing TestMerge")
	expected := 6
	one := []Clusterable{
		SimpleClusterable{0},
		SimpleClusterable{1},
		SimpleClusterable{2.1},
		SimpleClusterable{2.2},
		SimpleClusterable{2.3},
	}

	two := []Clusterable{
		one[0],
		one[1],
		SimpleClusterable{2.4},
	}
	visitMap := make(map[string]bool)
	output := mergeWithCluster(one, two, visitMap)
	assertEquals(t, expected, len(output))
}

func TestExpandCluster(t *testing.T) {
	log.Println("Executing TestExpandCluster")
	expected := 4
	clusterList := []Clusterable{
		SimpleClusterable{0},
		SimpleClusterable{1},
		SimpleClusterable{2},
		SimpleClusterable{2.1},
		SimpleClusterable{5},
	}

	eps := 1.0
	minPts := 3
	visitMap := make(map[string]bool)
	cluster := make(Cluster, 0)
	cluster = expandCluster(cluster, clusterList, visitMap, minPts, eps)
	assertEquals(t, expected, len(cluster))
}

func TestClusterize(t *testing.T) {
	log.Println("Executing TestClusterize")
	clusterList := []Clusterable{
		SimpleClusterable{1},
		SimpleClusterable{0.5},
		SimpleClusterable{0},
		SimpleClusterable{5},
		SimpleClusterable{4.5},
		SimpleClusterable{4},
	}
	eps := 1.0
	minPts := 2
	clusters := Clusterize(clusterList, minPts, eps)
	assertEquals(t, 2, len(clusters))
	if 2 == len(clusters) {
		assertEquals(t, 3, len(clusters[0]))
		assertEquals(t, 3, len(clusters[1]))
	}
}

func TestClusterizeNoData(t *testing.T) {
	log.Println("Executing TestClusterizeNoData")
	clusterList := []Clusterable{}
	eps := 1.0
	minPts := 3
	clusters := Clusterize(clusterList, minPts, eps)
	assertEquals(t, 0, len(clusters))
}

// TestNumberOfPoints checks whether number of clustered points 
// is not greater than it was before clustering
func TestNumberOfPoints(t *testing.T) {
	// Use seed = 1486926237 to discover a bug in 7672773 commit
	seed := time.Now().Unix()
	rand.Seed(seed)
	log.Println("Executing TestNumberOfPoints with random seed:", seed)
	// Take random number of points(range [5, 105))
	totalPoints := rand.Intn(100) + 5
	clusterList := make([]Clusterable, totalPoints)
	for i := range clusterList {
		clusterList[i] = SimpleClusterable{rand.Float64()}
	}
	// Random epsilon(range [0.2, 1)) and minPts(range [3, 13))
	eps, minPts := rand.Float64() + 0.1, rand.Intn(10) + 2
	clusters := Clusterize(clusterList, minPts, eps)
	nClustered := 0
	for _, cluster := range clusters {
		nClustered += len(cluster)
	}	
	// Check number of points
	if !(nClustered <= totalPoints) {
		t.Errorf("Got the greater number of clustered points than it " + 
			"was in total.\nTotal number of points: %d\nNumber of " +
			"points after clustering: %d", totalPoints, nClustered)
	}
}

//Assert function. If  the expected value not equals result, function
//returns error.
func assertEquals(t *testing.T, expected, result int) {
	if expected != result {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}
