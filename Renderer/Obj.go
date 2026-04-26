package main

import (
	"bufio"
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"path"
	"strings"

	"Renderer/Math"
	"Renderer/Rendering"
)

type ObjMaterial struct {
	Name  string
	MapKd string
}

type ObjContext struct {
	Vertices        []Math.Vec4
	Faces           []Rendering.Face
	TextureVertices []Rendering.UV
	VertexNormals   []Math.Vec4
	Textures        map[string]*Rendering.Texture

	VertexIndexOffset   int
	TextureVertexOffset int
	VertexNormalOffset  int
}

func (c *ObjContext) Clear() {
	c.VertexIndexOffset += len(c.Vertices)
	c.TextureVertexOffset += len(c.TextureVertices)
	c.VertexNormalOffset += len(c.VertexNormals)

	c.Vertices = nil
	c.Faces = nil
	c.TextureVertices = nil
	c.VertexNormals = nil
}

func parseVertex(line string) (Math.Vec4, error) {
	var x, y, z float32
	_, err := fmt.Sscanf(line, "v %f %f %f", &x, &y, &z)
	return Math.Vec4{x, y, z, 1}, err
}

func parseTextureVertex(line string) (Rendering.UV, error) {
	var x, y float32
	_, err := fmt.Sscanf(line, "vt %f %f", &x, &y)
	return Rendering.UV{x, y}, err
}

func parseVertexNormal(line string) (Math.Vec4, error) {
	var x, y, z float32
	_, err := fmt.Sscanf(line, "vn %f %f %f", &x, &y, &z)
	return Math.Vec4{x, y, z, 1}, err
}

type QuadFace struct {
	Triangles [2]Rendering.Face
	IsQuad    bool
}

func parseVertexRef(ref string) (v, vt, vn int, err error) {
	parts := strings.Split(ref, "/")
	if len(parts) < 1 {
		return 0, 0, 0, errors.New("invalid vertex reference")
	}

	_, err = fmt.Sscanf(parts[0], "%d", &v)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(parts) > 1 && parts[1] != "" {
		_, err = fmt.Sscanf(parts[1], "%d", &vt)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	if len(parts) > 2 && parts[2] != "" {
		_, err = fmt.Sscanf(parts[2], "%d", &vn)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	return v, vt, vn, nil
}

func parseFace(c *ObjContext, line string) (QuadFace, error) {
	fields := strings.Fields(line[2:])

	if len(fields) < 3 || len(fields) > 4 {
		return QuadFace{}, errors.New("face must have 3 or 4 vertices")
	}

	var qf QuadFace
	qf.IsQuad = len(fields) == 4

	refs := make([]struct{ v, vt, vn int }, len(fields))
	for i, field := range fields {
		v, vt, vn, err := parseVertexRef(field)
		if err != nil {
			return QuadFace{}, err
		}
		refs[i] = struct{ v, vt, vn int }{v, vt, vn}
	}

	face := Rendering.Face{}
	face.VertexIndices[0] = refs[0].v - c.VertexIndexOffset - 1
	face.VertexIndices[1] = refs[1].v - c.VertexIndexOffset - 1
	face.VertexIndices[2] = refs[2].v - c.VertexIndexOffset - 1

	if refs[0].vt > 0 {
		face.UVs[0] = c.TextureVertices[refs[0].vt-c.TextureVertexOffset-1]
		face.UVs[1] = c.TextureVertices[refs[1].vt-c.TextureVertexOffset-1]
		face.UVs[2] = c.TextureVertices[refs[2].vt-c.TextureVertexOffset-1]
	}

	if refs[0].vn > 0 {
		face.NormalIndices[0] = refs[0].vn - c.VertexNormalOffset - 1
		face.NormalIndices[1] = refs[1].vn - c.VertexNormalOffset - 1
		face.NormalIndices[2] = refs[2].vn - c.VertexNormalOffset - 1
	}

	qf.Triangles[0] = face

	if qf.IsQuad {
		face2 := Rendering.Face{}
		face2.VertexIndices[0] = refs[0].v - c.VertexIndexOffset - 1
		face2.VertexIndices[1] = refs[2].v - c.VertexIndexOffset - 1
		face2.VertexIndices[2] = refs[3].v - c.VertexIndexOffset - 1

		if refs[0].vt > 0 {
			face2.UVs[0] = c.TextureVertices[refs[0].vt-c.TextureVertexOffset-1]
			face2.UVs[1] = c.TextureVertices[refs[2].vt-c.TextureVertexOffset-1]
			face2.UVs[2] = c.TextureVertices[refs[3].vt-c.TextureVertexOffset-1]
		}

		if refs[0].vn > 0 {
			face2.NormalIndices[0] = refs[0].vn - c.VertexNormalOffset - 1
			face2.NormalIndices[1] = refs[2].vn - c.VertexNormalOffset - 1
			face2.NormalIndices[2] = refs[3].vn - c.VertexNormalOffset - 1
		}

		qf.Triangles[1] = face2
	}

	return qf, nil
}

func parseMtlLibFile(filename string) ([]ObjMaterial, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	var materials []ObjMaterial
	var mat *ObjMaterial

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(line, "newmtl "):
			if mat != nil {
				materials = append(materials, *mat)
			}
			name := strings.TrimPrefix(line, "newmtl ")
			mat = &ObjMaterial{Name: name}
		case strings.HasPrefix(line, "map_Kd "):
			mapKd := strings.TrimPrefix(line, "map_Kd ")
			mat.MapKd = mapKd
		}
	}

	if mat != nil {
		materials = append(materials, *mat)
	}

	return materials, nil
}

// LoadObjFile reads a mesh from an .obj file.
// Format description: https://people.computing.clemson.edu/~dhouse/courses/405/docs/brief-obj-file-format.html
func LoadObjFile(filename string, singleMesh bool) (meshes []*Rendering.Mesh, _ error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	dirname := path.Dir(filename)
	scanner := bufio.NewScanner(file)
	defaultTexture := Rendering.NewColorTexture(color.RGBA{255, 0, 255, 255})
	var currentTexture *Rendering.Texture

	c := &ObjContext{Textures: make(map[string]*Rendering.Texture)}
	textureFiles := make(map[string]*Rendering.Texture)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(line, "mtllib "):
			mtlLibFile := strings.TrimPrefix(line, "mtllib ")
			log.Printf("[INFO] found mtllib file: %s", mtlLibFile)

			materials, err := parseMtlLibFile(path.Join(dirname, mtlLibFile))
			if err != nil {
				return nil, fmt.Errorf("failed to parse material library: %s", err)
			}

			for _, m := range materials {
				if m.MapKd == "" {
					log.Printf("[INFO] using default texture for material: %s", m.Name)
					c.Textures[m.Name] = defaultTexture
				} else {
					if texture, ok := textureFiles[m.MapKd]; ok {
						c.Textures[m.Name] = texture
					} else {
						log.Printf("[INFO] loading texture: %s", m.MapKd)

						texturePath := m.MapKd
						if texturePath[0] != '/' {
							texturePath = path.Join(dirname, m.MapKd)
						}

						texture, err = Rendering.LoadTextureFile(texturePath)
						if err != nil {
							return nil, fmt.Errorf("failed to load texture: %s", err)
						}

						textureFiles[m.MapKd] = texture
						c.Textures[m.Name] = texture
					}
				}
			}

		case strings.HasPrefix(line, "o "):
			if len(c.Vertices) != 0 && !singleMesh {
				mesh := Rendering.NewMesh(c.Vertices, c.VertexNormals, c.Faces)
				meshes = append(meshes, mesh)
				c.Clear()
			}

		case strings.HasPrefix(line, "v "):
			v, err := parseVertex(line)
			if err != nil {
				return nil, err
			}
			c.Vertices = append(c.Vertices, v)

		case strings.HasPrefix(line, "vt "):
			vt, err := parseTextureVertex(line)
			if err != nil {
				return nil, err
			}
			c.TextureVertices = append(c.TextureVertices, vt)

		case strings.HasPrefix(line, "vn "):
			vn, err := parseVertexNormal(line)
			if err != nil {
				return nil, err
			}
			c.VertexNormals = append(c.VertexNormals, vn)

		case strings.HasPrefix(line, "usemtl "):
			mtlName := strings.TrimPrefix(line, "usemtl ")
			currentTexture = c.Textures[mtlName]

		case strings.HasPrefix(line, "f "):
			qf, err := parseFace(c, line)
			if err != nil {
				return nil, err
			}
			qf.Triangles[0].Texture = currentTexture
			c.Faces = append(c.Faces, qf.Triangles[0])
			if qf.IsQuad {
				qf.Triangles[1].Texture = currentTexture
				c.Faces = append(c.Faces, qf.Triangles[1])
			}
		}
	}

	if len(c.Vertices) != 0 {
		mesh := Rendering.NewMesh(c.Vertices, c.VertexNormals, c.Faces)
		meshes = append(meshes, mesh)
	}

	if len(meshes) == 0 {
		return nil, fmt.Errorf("obj file does not have any vertices data")
	}

	return meshes, nil
}

func LoadMeshFile(filename string, singleMesh bool) (meshes []*Rendering.Mesh, err error) {
	switch ext := path.Ext(filename); ext {
	case ".obj":
		meshes, err = LoadObjFile(filename, singleMesh)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported mesh format: %s", ext)
	}
	return meshes, nil
}
