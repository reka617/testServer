package manager

import (
	"encoding/json"
	"fmt"
	"os"

	n "github.com/hqpko/navmesh"
)

type Triangle struct {
	Indices [3]int `json:"indices"`
}

type NavMeshJsonData struct {
	Vertices  []*n.V3    `json:"vertices"`
	Triangles []Triangle `json:"triangles"`
}

type NavMeshManager struct {
	navMesh *n.NavMesh
}

var navMeshManager *NavMeshManager

func GetNavMeshManager() *NavMeshManager {
	if navMeshManager == nil {
		navMeshManager = &NavMeshManager{
			navMesh: &n.NavMesh{},
		}
		navMeshManager.LoadNavMeshData()
	}

	return navMeshManager
}

func (nm *NavMeshManager) PathFinding(srcX float64, srcY float64, srcZ float64,
	destX float64, destY float64, destZ float64) (*n.Path, error) {
	return nm.navMesh.FindingPath(&n.V3{X: srcX, Y: srcZ, Z: srcY}, &n.V3{
		X: destX,
		Y: destZ,
		Z: destY,
	})
}

func (nm *NavMeshManager) LoadNavMeshData() {
	file, err := os.Open("NavMeshData.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// JSON 디코더 생성
	decoder := json.NewDecoder(file)

	// NavMeshData 구조체 인스턴스 생성
	var navMeshData NavMeshJsonData

	// JSON 파일 디코딩
	err = decoder.Decode(&navMeshData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	nm.navMesh.Vertices = navMeshData.Vertices

	nm.navMesh.Triangles = make([][3]int32, len(navMeshData.Triangles))

	for i, triangle := range navMeshData.Triangles {
		nm.navMesh.Triangles[i] = [3]int32{
			int32(triangle.Indices[0]),
			int32(triangle.Indices[1]),
			int32(triangle.Indices[2]),
		}
	}

	nm.navMesh.Dijkstra.CreateMatrixFromMesh(navMeshData.Vertices, nm.navMesh.Triangles)

	print(len(nm.navMesh.Triangles))
}
