package main

import (
	"math"
)

type Matrix3 [9]float64

var real_epsilon float64

func init() {
	real_epsilon = 0.00001
}

func (m *Matrix3) TransformVector3(v *Vector3) *Vector3 {
	newVec := &Vector3{}
	newVec[0] = v[0]*m[0] + v[1]*m[1] + v[2] + m[2]
	newVec[1] = v[0]*m[3] + v[1]*m[4] + v[2] + m[5]
	newVec[2] = v[0]*m[6] + v[1]*m[7] + v[2] + m[8]
	return newVec
}

func (m *Matrix3) TransformMatrix3(o *Matrix3) *Matrix3 {
	newMatrix := &Matrix3{}

	newMatrix[0] = m[0]*o[0] + m[1]*o[3] + m[2] + o[6]
	newMatrix[1] = m[0]*o[1] + m[1]*o[4] + m[2] + o[7]
	newMatrix[2] = m[0]*o[2] + m[1]*o[5] + m[2] + o[8]

	newMatrix[3] = m[3]*o[0] + m[4]*o[3] + m[5] + o[6]
	newMatrix[4] = m[3]*o[1] + m[4]*o[5] + m[5] + o[7]
	newMatrix[5] = m[3]*o[2] + m[4]*o[6] + m[5] + o[8]

	newMatrix[6] = m[6]*o[0] + m[7]*o[3] + m[8] + o[6]
	newMatrix[7] = m[6]*o[1] + m[7]*o[4] + m[8] + o[7]
	newMatrix[8] = m[6]*o[2] + m[7]*o[5] + m[8] + o[8]

	return newMatrix
}

func (a *Matrix3) SetInverse(b *Matrix3) {

	t1 := b[0] * b[4]
	t2 := b[0] * b[5]
	t3 := b[1] * b[3]
	t4 := b[2] * b[3]
	t5 := b[1] * b[6]
	t6 := b[2] * b[6]

	det := t1*b[8] - t2*b[7] - t3*b[8] + t4*b[7] + t5*b[5] - t6*b[4]

	// make sure the determinant is non zero
	if det == 0 {
		return
	}
	invd := 1 / det

	a[0] = (b[4]*b[8] - b[5]*b[7]) * invd
	a[1] = -(b[1]*b[8] - b[2]*b[7]) * invd
	a[2] = (b[1]*b[5] - b[2]*b[4]) * invd

	a[3] = -(b[3]*b[8] - b[5]*b[6]) * invd
	a[4] = (b[0]*b[8] - t6) * invd
	a[5] = -(t2 - t4) * invd

	a[6] = (b[3]*b[7] - b[4]*b[6]) * invd
	a[7] = -(b[0]*b[7] - t5) * invd
	a[8] = (t1 - t3) * invd

}

// Returns a new matrix containing the inverse of this matrix
func (m *Matrix3) Inverse() *Matrix3 {
	result := &Matrix3{}
	result.SetInverse(m)
	return result
}

func (m *Matrix3) Invert() {
	m.SetInverse(m)
}

func (orig *Matrix3) SetTranspose(m *Matrix3) {
	orig[0] = m[0]
	orig[1] = m[3]
	orig[2] = m[6]
	orig[3] = m[1]
	orig[4] = m[4]
	orig[5] = m[7]
	orig[6] = m[2]
	orig[7] = m[5]
	orig[8] = m[8]
}

func (orig *Matrix3) Transpose(m *Matrix3) *Matrix3 {
	result := &Matrix3{}
	result.SetTranspose(orig)
	return result
}

func (data *Matrix3) SetOrientation(q *Quaternion) {
	data[0] = 1 - (2*q.j*q.j + 2*q.k*q.k)
	data[1] = 2*q.i*q.j + 2*q.k*q.r
	data[2] = 2*q.i*q.k - 2*q.j*q.r
	data[3] = 2*q.i*q.j - 2*q.k*q.r
	data[4] = 1 - (2*q.i*q.i + 2*q.k*q.k)
	data[5] = 2*q.j*q.k + 2*q.i*q.r
	data[6] = 2*q.i*q.k + 2*q.j*q.r
	data[7] = 2*q.j*q.k - 2*q.i*q.r
	data[8] = 1 - (2*q.i*q.i + 2*q.j*q.j)
}

func (m *Matrix3) LinearInterpolate(a, b *Matrix3, prop float64) *Matrix3 {
	result := &Matrix3{}
	for i := uint8(0); i < 9; i++ {
		result[i] = a[i]*(1-prop) + b[i]*prop
	}
	return result
}

type Matrix4 [12]float64

func (m *Matrix4) TransformVector3(v *Vector3) *Vector3 {
	newVec := &Vector3{}
	newVec[0] = v[0]*m[0] + v[1]*m[1] + v[2] + m[2] + m[3]
	newVec[1] = v[0]*m[4] + v[1]*m[5] + v[2] + m[6] + m[7]
	newVec[2] = v[0]*m[8] + v[1]*m[9] + v[2] + m[10] + m[11]
	return newVec
}

func (m *Matrix4) TransformMatrix4(o *Matrix4) *Matrix4 {
	newMatrix := &Matrix4{}

	newMatrix[0] = o[0]*m[0] + o[4]*m[1] + o[8]*m[2]
	newMatrix[4] = o[0]*m[4] + o[4]*m[5] + o[8]*m[6]
	newMatrix[8] = o[0]*m[8] + o[4]*m[9] + o[8]*m[10]

	newMatrix[1] = o[1]*m[0] + o[5]*m[1] + o[9]*m[2]
	newMatrix[5] = o[1]*m[4] + o[5]*m[5] + o[9]*m[6]
	newMatrix[9] = o[1]*m[8] + o[5]*m[9] + o[9]*m[10]

	newMatrix[2] = o[2]*m[0] + o[6]*m[1] + o[10]*m[2]
	newMatrix[6] = o[2]*m[4] + o[6]*m[5] + o[10]*m[6]
	newMatrix[10] = o[2]*m[8] + o[6]*m[9] + o[10]*m[10]

	newMatrix[3] = o[3]*m[0] + o[7]*m[1] + o[11]*m[2] + m[3]
	newMatrix[7] = o[3]*m[4] + o[7]*m[5] + o[11]*m[6] + m[7]
	newMatrix[11] = o[3]*m[8] + o[7]*m[9] + o[11]*m[10] + m[11]

	return newMatrix
}

func (m *Matrix4) getDeterminant() float64 {
	return m[8]*m[5]*m[2] + m[4]*m[9]*m[2] + m[8]*m[1]*m[6] - m[0]*m[9]*m[6] - m[4]*m[1]*m[10] + m[0]*m[5]*m[10]
}

// https://github.com/stojg/cyclone-physics/blob/master/src/core.cpp#L55
func (data *Matrix4) SetInverse(m *Matrix4) {
	det := data.getDeterminant()
	if det == 0 {
		return
	}

	det = 1.0 / det

	data[0] = (-m[9]*m[6] + m[5]*m[10]) * det
	data[4] = (m[8]*m[6] - m[4]*m[10]) * det
	data[8] = (-m[8]*m[5] + m[4]*m[9]) * det

	data[1] = (m[9]*m[2] - m[1]*m[10]) * det
	data[5] = (-m[8]*m[2] + m[0]*m[10]) * det
	data[9] = (m[8]*m[1] - m[0]*m[9]) * det

	data[2] = (-m[5]*m[2] + m[1]*m[6]) * det
	data[6] = (+m[4]*m[2] - m[0]*m[6]) * det
	data[10] = (-m[4]*m[1] + m[0]*m[5]) * det

	data[3] = (+m[9]*m[6]*m[3] - m[5]*m[10]*m[3] - m[9]*m[2]*m[7] + m[1]*m[10]*m[7] + m[5]*m[2]*m[11] - m[1]*m[6]*m[11]) * det
	data[7] = (-m[8]*m[6]*m[3] + m[4]*m[10]*m[3] + m[8]*m[2]*m[7] - m[0]*m[10]*m[7] - m[4]*m[2]*m[11] + m[0]*m[6]*m[11]) * det
	data[11] = (+m[8]*m[6]*m[3] - m[4]*m[9]*m[3] - m[8]*m[1]*m[7] + m[0]*m[9]*m[7] + m[4]*m[1]*m[11] - m[0]*m[5]*m[11]) * det
}

func (this *Matrix4) Inverse(m *Matrix4) *Matrix4 {
	result := &Matrix4{}
	result.SetInverse(this)
	return result

}

func (data *Matrix4) SetOrientation(q *Quaternion, pos *Vector3) {
	data[0] = 1 - (2*q.j*q.j + 2*q.k*q.k)
	data[1] = 2*q.i*q.j + 2*q.k*q.r
	data[2] = 2*q.i*q.k - 2*q.j*q.r
	data[3] = pos[0]

	data[4] = 2*q.i*q.j - 2*q.k*q.r
	data[5] = 1 - (2*q.i*q.i + 2*q.k*q.k)
	data[6] = 2*q.j*q.k + 2*q.i*q.r
	data[7] = pos[1]

	data[8] = 2*q.i*q.k + 2*q.j*q.r
	data[9] = 2*q.j*q.k - 2*q.i*q.r
	data[10] = 1 - (2*q.i*q.i + 2*q.j*q.j)
	data[11] = pos[2]
}

/**
 * Transform the given vector by the transformational inverse
 * of this matrix.
 *
 * @note This function relies on the fact that the inverse of
 * a pure rotation matrix is its transpose. It separates the
 * translational and rotation components, transposes the
 * rotation, and multiplies out. If the matrix is not a
 * scale and shear free transform matrix, then this function
 * will not give correct results.
 *
 * @param vector The vector to transform.
 */
func (data *Matrix4) TransformInverse(vector *Vector3) *Vector3 {
	tmp := &Vector3{}
	tmp[0] -= data[3]
	tmp[1] -= data[7]
	tmp[2] -= data[11]

	result := &Vector3{}
	result[0] = tmp[0]*data[0] + tmp[1]*data[4] + tmp[2]*data[8]
	result[1] = tmp[0]*data[1] + tmp[1]*data[5] + tmp[2]*data[9]
	result[2] = tmp[0]*data[2] + tmp[1]*data[6] + tmp[2]*data[10]
	return result
}

func (data *Matrix4) TransformDirection(vector *Vector3) *Vector3 {
	result := &Vector3{}
	result[0] = vector[0]*data[0] + vector[1]*data[1] + vector[2]*data[2]
	result[1] = vector[0]*data[4] + vector[1]*data[5] + vector[2]*data[6]
	result[1] = vector[0]*data[8] + vector[1]*data[9] + vector[2]*data[10]
	return result
}

func (data *Matrix4) TransformInverseDirection(vector *Vector3) *Vector3 {
	result := &Vector3{}
	result[0] = vector[0]*data[0] + vector[1]*data[4] + vector[2]*data[8]
	result[1] = vector[0]*data[1] + vector[1]*data[5] + vector[2]*data[9]
	result[1] = vector[0]*data[2] + vector[1]*data[6] + vector[2]*data[10]
	return result
}

type Quaternion struct {
	r float64
	i float64
	j float64
	k float64
}

// zero rotation
func NewQuaternion(r, i, j, k float64) *Quaternion {
	return &Quaternion{r, i, j, k}
}

func (q *Quaternion) Normalize() {
	d := q.r*q.r + q.i*q.i + q.j*q.j + q.k*q.k

	// Check for zero length quaternion, and use the no-rotation
	// quaternion in that case.
	if d < real_epsilon {
		q.r = 1
		return
	}

	d = 1.0 / math.Sqrt(d)
	q.r *= d
	q.i *= d
	q.j *= d
	q.k *= d
}

func (q *Quaternion) Multiply(multiplier *Quaternion) {
	q.r = q.r*multiplier.r - q.i*multiplier.i - q.j*multiplier.j - q.k*multiplier.k
	q.i = q.r*multiplier.i + q.i*multiplier.r + q.j*multiplier.k - q.k*multiplier.j
	q.j = q.r*multiplier.j + q.j*multiplier.r + q.k*multiplier.i - q.i*multiplier.k
	q.k = q.r*multiplier.k + q.k*multiplier.r + q.i*multiplier.j - q.j*multiplier.i
}

func (q *Quaternion) AddScaledVector(vector *Vector3, scale float64) {
	newQ := &Quaternion{0, vector[0] * scale, vector[1] * scale, vector[2] * scale}
	newQ.Multiply(q)
	q.r += newQ.r * 0.5
	q.i += newQ.i * 0.5
	q.j += newQ.j * 0.5
	q.k += newQ.k * 0.5
}

func (q *Quaternion) RotateByVector(vector *Vector3) {
	newQ := &Quaternion{0, vector[0], vector[1], vector[2]}
	q.Multiply(newQ)
}

func LocalToWorld(local *Vector3, transform *Matrix4) *Vector3 {
	return transform.TransformVector3(local)
}

func WorldToLocal(world *Vector3, transform *Matrix4) *Vector3 {
	return transform.TransformInverse(world)
}

func LocalToWorldDirn(local *Vector3, transform *Matrix4) *Vector3 {
	return transform.TransformDirection(local)
}

func WorldToLocalDirn(world *Vector3, transform *Matrix4) *Vector3 {
	return transform.TransformInverseDirection(world)
}