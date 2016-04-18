package entities

type SnakeParts struct {
  x int
  y int
}

type Snake struct {
  Id int16
  Username string
  Skin string
  XPos int32
  YPos int32
  XPosHead int32
  YPosHead
  Speed int16
  X int32
  D int32
  Parts SnakeParts
}
