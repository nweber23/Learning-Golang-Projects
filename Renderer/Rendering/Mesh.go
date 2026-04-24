package Rendering

import (
	"Renderer/Math"
)

type UV struct {
	U, V float32
}

type Face struct {
	VertexIndices [3]int
	NormalIndices [3]int
	UVs           [3]UV
	Texture       *Texture
}

type Mesh struct {
	Name          string
	Vertices      []Math.Vec4
	VertexNormals []Math.Vec4
	FaceNormals   []Math.Vec4
	BoundingBox   [8]Math.Vec4
	Faces         []Face
}

func boundingBox(vertices []Math.Vec4) [8]Math.Vec4 {
	minX, minY, minZ := vertices[0].X, vertices[0].Y, vertices[0].Z
	maxX, maxY, maxZ := minX, minY, minZ
	for _, v := range vertices {
		minX = min(minX, v.X)
		minY = min(minY, v.Y)
		minZ = min(minZ, v.Z)
		maxX = max(maxX, v.X)
		maxY = max(maxY, v.Y)
		maxZ = max(maxZ, v.Z)
	}
	return [8]Math.Vec4{
		{minX, minY, minZ, 1},
		{minX, minY, maxZ, 1},
		{minX, maxY, minZ, 1},
		{minX, maxY, maxZ, 1},
		{maxX, minY, minZ, 1},
		{maxX, minY, maxZ, 1},
		{maxX, maxY, minZ, 1},
		{maxX, maxY, maxZ, 1},
	}
}

func NewMesh(vertices []Math.Vec4, vertexNormals []Math.Vec4, faces []Face) *Mesh {
	faceNormals := make([]Math.Vec4, len(faces))
	for i := range faces {
		v0 := vertices[faces[i].VertexIndices[0]].ToVec3()
		v1 := vertices[faces[i].VertexIndices[1]].ToVec3()
		v2 := vertices[faces[i].VertexIndices[2]].ToVec3()
		faceNormals[i] = v1.Sub(v0).Cross(v2.Sub(v0)).Normalize().ToVec4()
	}
	return &Mesh{
		Faces:         faces,
		Vertices:      vertices,
		VertexNormals: vertexNormals,
		FaceNormals:   faceNormals,
		BoundingBox:   boundingBox(vertices),
	}
}

type Object struct {
	*Mesh
	Rotation            Math.Vec3
	Translation         Math.Vec3
	Scale               Math.Vec3
	TransformedVertices []Math.Vec4
	WorldVertexNormals  []Math.Vec4
	WorldFaceNormals    []Math.Vec4
}

func NewObject(mesh *Mesh) *Object {
	return &Object{
		Mesh:                mesh,
		Scale:               Math.Vec3{1, 1, 1},
		TransformedVertices: make([]Math.Vec4, len(mesh.Vertices)),
		WorldFaceNormals:    make([]Math.Vec4, len(mesh.FaceNormals)),
		WorldVertexNormals:  make([]Math.Vec4, len(mesh.VertexNormals)),
	}
}
