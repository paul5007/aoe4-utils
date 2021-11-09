package unit

type NameType int

const (
	NoName NameType = iota
	Villager
	Spearman
	ManAtArms
	Horseman
	Knight
	Archer
	Crossbow
)

func (n NameType) String() string {
	switch n {
	case NoName:
		return "NoName"
	case Villager:
		return "Villagers"
	case Spearman:
		return "Spearman"
	case Horseman:
		return "Horseman"
	case Knight:
		return "Knight"
	case Archer:
		return "Archer"
	case Crossbow:
		return "Crossbow"
	default:
		return "UNKNOWN"
	}
}

type Job int

const (
	New Job = iota
	Sheep
	Berry
	Hunting
	Wood
	Gold
	Stone
	Military
	Builder
)

func (j Job) String() string {
	switch j {
	case New:
		return "New"
	case Sheep:
		return "Sheep"
	case Berry:
		return "Berry"
	case Hunting:
		return "Hunting"
	case Wood:
		return "Wood"
	case Gold:
		return "Gold"
	case Stone:
		return "Stone"
	case Military:
		return "Military"
	case Builder:
		return "Builder"
	default:
		return "UNKNOWN"
	}
}

type Unit struct {
	Name      NameType
	Job       Job
	IdleTime  int
	BuildTime int
	FoodCost  float32
	WoodCost  float32
	GoldCost  float32
}
