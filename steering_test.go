package main

//func deg2rad(degree float64) float64 {
//	return degree * (math.Pi / 180)
//}
//
//func rad2deg(radians float64) float64 {
//	return radians * (180 / math.Pi)
//}
//
//func TestAlignNoRotation(t *testing.T) {
//	character := NewEntity()
//	target := NewEntity()
//
//	var alignNoRotationTests = []struct {
//		character *Quaternion
//		target    *Quaternion
//		expected  *Vector3
//	}{
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(0)), &Vector3{0, 0, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(34)), QuaternionFromAxisAngle(VectorY(), deg2rad(34)), &Vector3{0, 0, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(90)), QuaternionFromAxisAngle(VectorY(), deg2rad(90)), &Vector3{0, 0, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(180)), QuaternionFromAxisAngle(VectorY(), deg2rad(180)), &Vector3{0, 0, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(234)), QuaternionFromAxisAngle(VectorY(), deg2rad(234)), &Vector3{0, 0, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(270)), QuaternionFromAxisAngle(VectorY(), deg2rad(270)), &Vector3{0, 0, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(270)), QuaternionFromAxisAngle(VectorY(), deg2rad(270)), &Vector3{0, 0, 0}},
//	}
//
//	for i := range alignNoRotationTests {
//		character.Orientation = alignNoRotationTests[i].character
//		target.Orientation = alignNoRotationTests[i].target
//
//		character.Body.calculateDerivedData(character)
//		target.Body.calculateDerivedData(target)
//
//		align := NewAlign(character, target, 0.5, 0.01, 0.1)
//		steering := align.GetSteering()
//
//		if !steering.angular.Equals(alignNoRotationTests[i].expected) {
//			t.Errorf("Expected %v, but got %v for test %d", alignNoRotationTests[i].expected, steering.angular, i+1)
//		}
//	}
//}
//
//func TestAlignRotation(t *testing.T) {
//	character := NewEntity()
//	target := NewEntity()
//
//	var alignTests = []struct {
//		character *Quaternion
//		target    *Quaternion
//		expected  *Vector3
//	}{
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(45)), &Vector3{0, 15.707963267948966, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(90)), &Vector3{0, 15.707963267948966, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(179)), &Vector3{0, 15.707963267948966, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(180)), &Vector3{0, 15.707963267948966, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(181)), &Vector3{0, -15.707963267948966, 0}},
//		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(270)), &Vector3{0, -15.707963267948966, 0}},
//	}
//
//	for i := range alignTests {
//		character.Orientation = alignTests[i].character
//		target.Orientation = alignTests[i].target
//
//		character.Body.calculateDerivedData(character)
//		target.Body.calculateDerivedData(target)
//
//		align := NewAlign(character, target, 0.001, 0.002, 0.1)
//		steering := align.GetSteering()
//
//		if !steering.angular.Equals(alignTests[i].expected) {
//			t.Errorf("Expected %v, but got %v for test %d", alignTests[i].expected, steering.angular, i+1)
//		}
//	}
//
//}
//
//func TestFace_calculateOrientation(t *testing.T) {
//	character := NewEntity()
//	target := NewEntity()
//
//	face := NewFace(character, target)
//
//	var alignTests = []struct {
//		target   *Vector3
//		expected *Quaternion
//	}{
//		{&Vector3{1, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(0))},
//		{&Vector3{-1, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(180))},
//		{&Vector3{0, 0, 1}, QuaternionFromAxisAngle(VectorY(), deg2rad(-90))},
//		{&Vector3{0, 0, -1}, QuaternionFromAxisAngle(VectorY(), deg2rad(90))},
//		{&Vector3{1, 0, 1}, QuaternionFromAxisAngle(VectorY(), deg2rad(-45))},
//		{&Vector3{-1, 0, -1}, QuaternionFromAxisAngle(VectorY(), deg2rad(135))},
//		{&Vector3{1, 0, -1}, QuaternionFromAxisAngle(VectorY(), deg2rad(45))},
//		{&Vector3{-1, 0, 1}, QuaternionFromAxisAngle(VectorY(), deg2rad(-135))},
//		{&Vector3{-0.5, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(180))},
//		{&Vector3{0.5, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(0))},
//		{&Vector3{0.5, 0, 0.5}, QuaternionFromAxisAngle(VectorY(), deg2rad(-45))},
//		{&Vector3{-0.5, 0, -0.5}, QuaternionFromAxisAngle(VectorY(), deg2rad(135))},
//		{&Vector3{1, 0, 0.99}, QuaternionFromAxisAngle(VectorY(), deg2rad(-44.712083933442905))},
//	}
//
//	for i, test := range alignTests {
//
//		actual := face.calculateOrientation(test.target)
//		if !actual.Equals(test.expected) {
//			t.Errorf("Expected %v, but got %v for test %d", alignTests[i].expected, actual, i+1)
//		}
//	}
//}
//
//func TestPath_getParam(t *testing.T) {
//	points := []*Vector3{
//		&Vector3{0, 0, 0},
//		&Vector3{1, 0, 0},
//		&Vector3{2, 0, 0},
//		&Vector3{3, 0, 0},
//		&Vector3{4, 0, 0},
//		&Vector3{5, 0, 0},
//	}
//
//	path := &Path{
//		points: points,
//	}
//
//	param := path.getParam(&Vector3{}, 0)
//	if param != 0 {
//		t.Errorf("Expected param to be 0, got %d", param)
//	}
//
//	param = path.getParam(&Vector3{}, 2)
//	if param != 0 {
//		t.Errorf("Expected param to be 0, got %d", param)
//	}
//
//	param = path.getParam(&Vector3{3, 0, 0}, 2)
//	if param != 3 {
//		t.Errorf("Expected param to be 3, got %d", param)
//	}
//}
//
//func TestPath_getPosition(t *testing.T) {
//	points := []*Vector3{
//		&Vector3{0, 0, 0},
//		&Vector3{1, 0, 0},
//		&Vector3{2, 0, 0},
//		&Vector3{3, 0, 0},
//		&Vector3{4, 0, 0},
//		&Vector3{5, 0, 0},
//	}
//
//	path := &Path{
//		points: points,
//	}
//
//	expected := points[0]
//	pos := path.getPosition(0)
//	if !pos.Equals(expected) {
//		t.Errorf("Expected position to be %v, got %v", expected, pos)
//	}
//
//	expected = points[1]
//	pos = path.getPosition(1)
//	if !pos.Equals(expected) {
//		t.Errorf("Expected position to be %v, got %v", expected, pos)
//	}
//}
