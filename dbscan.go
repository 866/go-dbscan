package dbscan

const (
	NOISE     = false
	CLUSTERED = true
)

type Clusterable interface {
	Distance(c interface{}) float64
	GetID() string
}

type Cluster []Clusterable

func Clusterize(objects []Clusterable, minPts int, eps float64) []Cluster {
	clusters := make([]Cluster, 0)
	visited := map[string]bool{}
	for _, point := range objects {
		if _, isVisited := visited[point.GetID()]; isVisited {
			continue
		}
		neighbours := findUnclusteredNeighbours(point, objects, visited, eps)
		if len(neighbours)+1 >= minPts {
			visited[point.GetID()] = CLUSTERED
			cluster := make(Cluster, 1)
			cluster[0] = point
			cluster = expandCluster(cluster, neighbours, visited, minPts, eps)

			if len(cluster) >= minPts {
				clusters = append(clusters, cluster)
			}
		} else {
			visited[point.GetID()] = NOISE
		}
	}
	return clusters
}

//Finds the neighbours from given array
//depends on Eps variable, which determines
//the distance limit from the point
func findUnclusteredNeighbours(point Clusterable, points []Clusterable, visited map[string]bool, eps float64) []Clusterable {
	neighbours := make([]Clusterable, 0)
	pointID := point.GetID()
	for _, potNeigb := range points {
		neigbID := potNeigb.GetID()
		pointState, isVisited := visited[neigbID]
		// Check whether the point is unclustered
		// Doing doublecheck(isVisited+pointState) 
		// because pointState could have value  
		// CLUSTERED by default		
		if (isVisited && pointState == CLUSTERED) {
			continue
		}
		if pointID != neigbID && potNeigb.Distance(point) <= eps {
			neighbours = append(neighbours, potNeigb)
		}
	}
	return neighbours
}

//Try to expand existing clutser
func expandCluster(cluster Cluster, neighbours []Clusterable, visited map[string]bool, minPts int, eps float64) Cluster {
	seed := make([]Clusterable, len(neighbours))
	copy(seed, neighbours)
	for _, point := range seed {
		pointState, isVisited := visited[point.GetID()]
		if !isVisited {
			currentNeighbours := findUnclusteredNeighbours(point, seed, visited, eps)
			if len(currentNeighbours)+1 >= minPts {
				cluster = mergeWithCluster(cluster, currentNeighbours, visited)
			}
		}
		if isVisited && pointState == NOISE {
			visited[point.GetID()] = CLUSTERED
			cluster = append(cluster, point)
		}
	}
	return cluster
}

func mergeWithCluster(clusterPoints []Clusterable, newPoints []Clusterable, visited map[string]bool) []Clusterable {
	mergeMap := make(map[string]Clusterable)
	// Mark all new points as clustered
	for _, point := range newPoints {
		visited[point.GetID()] = CLUSTERED
	}
	putAll(mergeMap, clusterPoints) 
	putAll(mergeMap, newPoints)
	merged := make([]Clusterable, 0)
	for _, val := range mergeMap {
		merged = append(merged, val)
	}

	return merged
}

//Function to add all values from list to map
//map keys is then the unique collecton from list
func putAll(m map[string]Clusterable, list []Clusterable) {
	for _, val := range list {
		m[val.GetID()] = val
	}
}
