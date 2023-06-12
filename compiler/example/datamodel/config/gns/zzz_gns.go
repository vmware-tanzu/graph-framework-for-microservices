package gns

type AdditionalDescription struct {
	DiscriptionA string
	DiscriptionB string
	DiscriptionC string
	DiscriptionD string
}

type AdditionalStatus struct {
	StatusX int
	StatusY int
}

type TempConst1 string
type TempConst2 string
type TempConst3 string

const (
	Const3 TempConst3 = "Const3"
	Const2 TempConst2 = "Const2"
	Const1 TempConst1 = "Const1"
)
