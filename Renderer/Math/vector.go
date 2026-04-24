package Math

type Vec2 struct {
	X float32
	Y float32
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

func (v Vec2) Multiply(scalar float32) Vec2 {
	return Vec2{v.X * scalar, v.Y * scalar}
}

func (v Vec2) Divide(scalar float32) Vec2 {
	return Vec2{v.X / scalar, v.Y / scalar}
}

func (v Vec2) Length() float32 {
	return sqrt32(v.X*v.X + v.Y*v.Y)
}

func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) Normalize() Vec2 {
	return v.Divide(v.Length())
}

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func Vec3FromArray(array [3]float32) Vec3 {
	return Vec3{array[0], array[1], array[2]}
}

func (v Vec3) ToVec4() Vec4 {
	return Vec4{v.X, v.Y, v.Z, 1}
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v Vec3) Multiply(scalar float32) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec3) Divide(scalar float32) Vec3 {
	return Vec3{v.X / scalar, v.Y / scalar, v.Z / scalar}
}

func (v Vec3) Length() float32 {
	return sqrt32(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vec3) Normalize() Vec3 {
	return v.Divide(v.Length())
}

func (v Vec3) Cross(other Vec3) Vec3 {
	x := v.Y*other.Z - v.Z*other.Y
	y := v.Z*other.X - v.X*other.Z
	z := v.X*other.Y - v.Y*other.X
	return Vec3{x, y, z}
}

func (v Vec3) ToRadians() Vec3 {
	f := pi32 / 180
	return Vec3{v.X * f, v.Y * f, v.Z * f}
}

type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

func (v Vec4) ToVec3() Vec3 {
	return Vec3{v.X, v.Y, v.Z}
}

func (v Vec4) Add(other Vec4) Vec4 {
	return Vec4{v.X + other.X, v.Y + other.Y, v.Z + other.Z, v.W + other.W}
}

func (v Vec4) Sub(other Vec4) Vec4 {
	return Vec4{v.X - other.X, v.Y - other.Y, v.Z - other.Z, v.W - other.W}
}

func (v Vec4) Multiply(scalar float32) Vec4 {
	return Vec4{v.X * scalar, v.Y * scalar, v.Z * scalar, v.W * scalar}
}

func (v Vec4) Divide(scalar float32) Vec4 {
	return Vec4{v.X / scalar, v.Y / scalar, v.Z / scalar, v.W / scalar}
}

func (v Vec4) Length() float32 {
	return sqrt32(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)
}

func (v Vec4) Dot(other Vec4) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

func (v Vec4) Conjugate() Vec4 {
	return Vec4{-v.X, -v.Y, -v.Z, -v.W}
}

func (v Vec4) Normalize() Vec4 {
	return v.Divide(v.Length())
}
