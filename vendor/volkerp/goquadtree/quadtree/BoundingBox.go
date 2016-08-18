/*
Copyright 2013 Volker Poplawski
*/


package quadtree

import "math"


// Use NewBoundingBox() to construct a BoundingBox object
type BoundingBox struct {
  MinX, MaxX, MinY, MaxY, MinZ, MaxZ float64
}


func NewBoundingBox(xA, xB, zA, zB float64) BoundingBox {
  return BoundingBox{
          MinX: math.Min(xA, xB),
          MaxX: math.Max(xA, xB),
          MinZ: math.Min(zA, zB),
          MaxZ: math.Max(zA, zB),
  }
}


// Make BoundingBox implement the BoundingBoxer interface
func (b BoundingBox) BoundingBox() BoundingBox {
  return b
}


func (b BoundingBox) SizeX() float64 {
  return b.MaxX - b.MinX
}


func (b BoundingBox) SizeY() float64 {
  return b.MaxY - b.MinY
}


// Returns true if o intersects this
func (b BoundingBox) Intersects(o BoundingBox) bool {
  return b.MinX < o.MaxX && b.MinY < o.MaxY &&
         b.MaxX > o.MinX && b.MaxY > o.MinY
}


// Returns true if o is within this
func (b BoundingBox) Contains(o BoundingBox) bool {
  return b.MinX <= o.MinX && b.MinY <= o.MinY &&
         b.MaxX >= o.MaxX && b.MaxY >= o.MaxY
}


