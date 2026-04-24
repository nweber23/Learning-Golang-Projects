//go:build amd64 && !purego

package Math

//go:noescape
func _matrixMultiplyVec4SSE(mat *Matrix, vecs []Vec4)

func MatrixMultiplyVec4Batch(m *Matrix, vec []Vec4) {
	mat := new(Matrix)
	*mat = (*m).Transpose()
	_matrixMultiplyVec4SSE(mat, vec)
}
