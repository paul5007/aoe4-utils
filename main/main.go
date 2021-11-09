package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	BuildTimeVillager      float32 = 20
	BuildTimeHouse         float32 = 20
	BuildTimeResourcePoint float32 = 20
	BuildTimeMilitary      float32 = 30
	BuildTimeMarket        float32 = 30
	BuildTimeLandmark1     float32 = 190
	BuildTimeLandmark2     float32 = 220
	BuildTimeTownCenter    float32 = 120
	BuildTimeSpearman      float32 = 15
	BuildTimeManAtArms     float32 = 22
	BuildTimeHorseman      float32 = 22
	BuildTimeKnight        float32 = 35
	BuildTimeArcher        float32 = 15
	BuildTimeCrossbow      float32 = 22
	CostVillager           float32 = 50
	CostHouse              float32 = 50
	CostFarm               float32 = 75
	CostMarket             float32 = 100
	CostResourcePoint      float32 = 50
	CostMilitaryBuilding   float32 = 150
	CostTownCenterWood     float32 = 400
	CostTownCenterStone    float32 = 300
	CostAge2Food           float32 = 400
	CostAge2Gold           float32 = 200
	CostAge3Food           float32 = 1200
	CostAge3Gold           float32 = 600
	CostSpearmanFood       float32 = 60
	CostSpearmanWood       float32 = 20
	CostManAtArmsFood      float32 = 100
	CostManAtArmsGold      float32 = 20
	CostHorsemanFood       float32 = 100
	CostHorsemanWood       float32 = 20
	CostKnightFood         float32 = 140
	CostKnightGold         float32 = 100
	CostArcherFood         float32 = 30
	CostArcherWood         float32 = 50
	CostCrossbowFood       float32 = 80
	CostCrossbowGold       float32 = 40
	IdleTimeWalk           float32 = 10
	HousingPer             float32 = 10
)

type BuildingType int

const (
	Barracks BuildingType = iota
	Stable
	ArcheryRange
	TownCenter
	House
	Landmark1
	Landmark2
)

func (b BuildingType) String() string {
	switch b {
	case Barracks:
		return "Barracks"
	case Stable:
		return "Stable"
	case ArcheryRange:
		return "ArcheryRange"
	case TownCenter:
		return "TownCenter"
	case House:
		return "House"
	case Landmark1:
		return "Landmark1"
	case Landmark2:
		return "Landmark2"
	default:
		return "UNKNOWN"
	}
}

type ResearchType int

const (
	NoResearch ResearchType = iota
	Forestry
	DoubleBroadax
	SpecializedPick
	Horticulture
	WheelBarrow
	SurvivalTechniques
)

type Building struct {
	Type      BuildingType
	UnitQueue []Unit
	BuildTime int
	FoodCost  float32
	WoodCost  float32
	GoldCost  float32
	StoneCost float32
}

type Research struct {
	Type         ResearchType
	ResearchTime int
}

const GatherSpeedBerry float32 = 0.56 * (1 - 0.1)
const GatherSpeedHunting float32 = 0.646 * (1 - 0.1)
const GatherSpeedSheep float32 = 0.646 * (1 - 0.1)
const GatherSpeedWood float32 = 0.565 * (1 - 0.1)
const GatherSpeedGold float32 = 0.646 * (1 - 0.1)
const GatherSpeedStone float32 = 0.646 * (1 - 0.1)

// const GatherSpeedBerry float32 = 0.56
// const GatherSpeedHunting float32 = 0.646
// const GatherSpeedSheep float32 = 0.646
// const GatherSpeedWood float32 = 0.565
// const GatherSpeedGold float32 = 0.646
// const GatherSpeedStone float32 = 0.646

var TotalFood float32 = 200
var TotalWood float32 = 150
var TotalGold float32 = 100
var TotalStone float32 = 0

var TotalHousing int = 10

var BuiltBerryMill bool = false
var BuiltHuntingMill bool = false
var BuiltLumberCamp bool = false
var BuiltGoldMine bool = false
var BuiltStoneMine bool = false

var Age1Military int = 0
var Age2 bool = false
var Age3 bool = false
var Age2Military int = 0
var Age2TownCenter = false

var TurnTime = 20
var Units []Unit
var Buildings []Building
var BuildOrder []string
var GameTime = 5
var reader = bufio.NewReader(os.Stdin)

var fastAge2 bool
var isFrench bool
var isEnglish bool

var unitMap = make(map[string]Unit)
var buildingMap = make(map[string]Building)
var researchMap = make(map[string]Research)

func setupMaps() {
	unitMap[Villager.String()] = Unit{Name: Villager, Job: New, FoodCost: 50, BuildTime: 20}
	unitMap[Spearman.String()] = Unit{Name: Spearman, Job: Military, FoodCost: 60, WoodCost: 20, BuildTime: 15}
	unitMap[ManAtArms.String()] = Unit{Name: ManAtArms, Job: Military, FoodCost: 100, GoldCost: 20, BuildTime: 22}
	unitMap[Horseman.String()] = Unit{Name: Horseman, Job: Military, FoodCost: 100, WoodCost: 20, BuildTime: 22}
	unitMap[Knight.String()] = Unit{Name: Knight, Job: Military, FoodCost: 140, GoldCost: 100, BuildTime: 35}
	unitMap[Archer.String()] = Unit{Name: Archer, Job: Military, FoodCost: 40, WoodCost: 50, BuildTime: 15}
	unitMap[Crossbow.String()] = Unit{Name: Crossbow, Job: Military, FoodCost: 80, WoodCost: 40, BuildTime: 22}

	buildingMap[TownCenter.String()] = Building{Type: TownCenter, BuildTime: 120, WoodCost: 400, StoneCost: 300}
	buildingMap[House.String()] = Building{Type: House, BuildTime: 20, WoodCost: 50}
	buildingMap[Landmark1.String()] = Building{Type: Landmark1, BuildTime: 190, FoodCost: 400, GoldCost: 200}
	buildingMap[Landmark2.String()] = Building{Type: Landmark2, BuildTime: 210, FoodCost: 1200, GoldCost: 600}
	buildingMap[Barracks.String()] = Building{Type: Barracks, BuildTime: 30, WoodCost: 150}
	buildingMap[Stable.String()] = Building{Type: Stable, BuildTime: 30, WoodCost: 150}
	buildingMap[ArcheryRange.String()] = Building{Type: ArcheryRange, BuildTime: 30, WoodCost: 150}

}

func main() {
	flag.BoolVar(&fastAge2, "fastAge2", false, "Get to Age2 as fast as possible")
	flag.BoolVar(&isFrench, "isFrench", false, "Playing French")
	flag.BoolVar(&isEnglish, "isEnglish", false, "Playing English")
	flag.Parse()

	if isFrench {
		TurnTime = BuildTimeVillager * 0.9
	}

	setupMaps()

	// Only add 300 food because we create 6 vills to start
	TotalFood += 300
	bo := fmt.Sprintf("------------------------------")
	BuildOrder = append(BuildOrder, bo)
	u := Unit{Job: Military}
	Units = append(Units, u)
	createVillager()
	createVillager()
	createVillager()
	createVillager()
	createVillager()
	createVillager()

	BuildOrder = append(BuildOrder, bo)

	for GameTime < 20*60 {
		checkHousing()

		checkThresholds()

		if GameTime%TurnTime == 0 {
			printGameTime()

			createVillager()

			printReport()
		}

		checkBuilders()

		gatherResources()

		GameTime += 1
	}

}

func checkHousing() {
	// Check if housing needed
	if len(Units) >= TotalHousing-1 && TotalWood >= 50 {
		// Used villager on Food to make house
		fmt.Println("BUILDING A HOUSE")
		for i, u := range Units {
			var ptr = &u
			if (u.Job == Sheep || u.Job == Berry) && u.IdleTime < 1 {
				ptr.IdleTime = BuildTimeHouse
				ptr.Job = Builder
				Units[i] = *ptr
				break
			}
		}
		TotalWood -= CostHouse
		TotalHousing += HousingPer
	}
}

func checkBuilders() {
	for i, u := range Units {
		var ptr = &u
		if u.Job == Builder && u.IdleTime < 1 {
			ptr.Job = Sheep
			Units[i] = *ptr
			break
		}
	}
}

func checkThresholds() {
	// Age Up
	if !Age2 && TotalFood >= CostAge2Food && TotalGold >= CostAge2Gold {
		Age2 = true
		fmt.Println("BUILDING AGE2 LANDMARK")
		TotalFood -= CostAge2Food
		TotalGold -= CostAge2Gold
		done := false
		var villCount int
		var tempVillCount int
		var tempIdleTime float32
		for !done {
			done = true
			fmt.Print("How many vills?: ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n")
			villCount, _ = strconv.Atoi(text)

			tempVillCount = villCount
			if villCount == 1 {
				tempIdleTime = BuildTimeLandmark1
			} else if villCount == 2 {
				tempIdleTime = BuildTimeLandmark1 * (1 - 0.24)
			} else if villCount == 3 {
				tempIdleTime = BuildTimeLandmark1 * (1 - 0.37)
			} else if villCount == 4 {
				tempIdleTime = BuildTimeLandmark1 * (1 - 0.47)
			} else if villCount == 5 {
				tempIdleTime = BuildTimeLandmark1 * (1 - 0.57)
			} else {
				done = false
			}
		}

		for i, u := range Units {
			var ptr = &u
			if (ptr.Job == Sheep || ptr.Job == Berry) && ptr.IdleTime < 1 {
				ptr.IdleTime = int(tempIdleTime)
				ptr.Job = Builder
				Units[i] = *ptr
				villCount -= 1
				if villCount == 0 {
					break
				}
			}
		}

		bo := fmt.Sprintf("Started Constructing Tier2 Landmark with %d Villagers: %02d:%02d", tempVillCount, GameTime/60, GameTime%60)
		BuildOrder = append(BuildOrder, bo)
		tempTime := GameTime + int(tempIdleTime)
		mins := tempTime / 60
		secs := tempTime % 60
		bo = fmt.Sprintf("Complete Constructing Tier2 Landmark: %02d:%02d", mins, secs)
		BuildOrder = append(BuildOrder, bo)

		reallocateWorkers()
	}

	// Age Again
	if !Age3 && TotalFood >= CostAge3Food && TotalGold >= CostAge3Gold {
		fmt.Println("BUILDING AGE3 LANDMARK")
		TotalFood -= CostAge3Food
		TotalGold -= CostAge3Gold
		for i, u := range Units {
			var ptr = &u
			if u.Job == Sheep || u.Job == Berry {
				u.IdleTime = BuildTimeLandmark2
				Units[i] = *ptr
				break
			}
		}

		bo := fmt.Sprintf("Started Constructing Tier3 Landmark: %02d:%02d", GameTime/60, GameTime%60)
		BuildOrder = append(BuildOrder, bo)
		tempTime := GameTime + BuildTimeLandmark2
		mins := tempTime / 60
		secs := tempTime % 60
		bo = fmt.Sprintf("Complete Constructing Tier3 Landmark: %02d:%02d", mins, secs)
		BuildOrder = append(BuildOrder, bo)

		reallocateWorkers()
	}

	// Second TC
	if !Age2TownCenter && TotalWood >= CostTownCenterStone && TotalStone >= CostTownCenterWood {
		fmt.Println("BUILDING AGE2 TOWNCENTER")
		TotalWood -= CostTownCenterWood
		TotalStone -= CostTownCenterStone
		for i, u := range Units {
			var ptr = &u
			if u.Job == Sheep || u.Job == Berry {
				u.IdleTime = BuildTimeTownCenter
				Units[i] = *ptr
				break
			}
		}

		bo := fmt.Sprintf("Started Constructing TownCenter: %02d:%02d", GameTime/60, GameTime%60)
		BuildOrder = append(BuildOrder, bo)
		tempTime := GameTime + BuildTimeTownCenter
		mins := tempTime / 60
		secs := tempTime % 60
		bo = fmt.Sprintf("Complete Constructing TownCenter: %02d:%02d", mins, secs)
		BuildOrder = append(BuildOrder, bo)

		reallocateWorkers()
	}
}

func printGameTime() {
	// Print current game time
	mins := GameTime / 60
	secs := GameTime % 60
	fmt.Printf("Current Time: %2d:%2d\n", mins, secs)
}

func createVillager() {
	// Create new villager if you have enough housing
	if TotalHousing > len(Units) && TotalFood > CostVillager {
		fmt.Println("Spawned New Villager")
		TotalFood -= CostVillager
		var u = &Unit{}
		changeJob(u)
		Units = append(Units, *u)
	}
}

func checkMakeBuilding() {
	if TotalWood >= CostMilitaryBuilding {
		done := false
		choseOne := false
		for !done && !choseOne {
			fmt.Print("Choose Building: barr, stab, arch, no: ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n")
			switch text {
			case "barr":
				TotalWood -= CostMilitaryBuilding
				b := Building{Type: Barracks, BuildTime: BuildTimeMilitary}
				Buildings = append(Buildings, b)
				choseOne = true
			case "stab":
				TotalWood -= CostMilitaryBuilding
				b := Building{Type: Stable, BuildTime: BuildTimeMilitary}
				Buildings = append(Buildings, b)
				choseOne = true
			case "arch":
				TotalWood -= CostMilitaryBuilding
				b := Building{Type: ArcheryRange, BuildTime: BuildTimeMilitary}
				Buildings = append(Buildings, b)
				choseOne = true
			case "no":
				done = true
			default:
				fmt.Println("Invalid option")
			}
		}
	}

}

func checkMakeUnit() {
	if len(Buildings) > 0 {
		done := false
		for !done {
			fmt.Print("Choose Building to Make Unit From: barr, stab, arch, no: ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n")
			switch text {
			case "barr":
				for i, b := range Buildings {
					var ptr = &b
					if ptr.Type == Barracks {
						makeNewUnit(ptr)
					}
					Buildings[i] = *ptr
				}
			case "stab":
				for i, b := range Buildings {
					var ptr = &b
					if ptr.Type == Stable {
						makeNewUnit(ptr)
					}
					Buildings[i] = *ptr
				}
			case "arch":
				for i, b := range Buildings {
					var ptr = &b
					if ptr.Type == ArcheryRange {
						makeNewUnit(ptr)
					}
					Buildings[i] = *ptr
				}
			case "no":
				done = true
			default:
				fmt.Println("Invalid option")
			}
		}
	}
}

func makeNewUnit(b *Building) {
	if b.Type == Barracks {
		fmt.Print("Choose a Unit: spr, man, no: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		switch text {
		case "spr":
			if TotalFood > CostSpearmanFood && TotalWood > CostSpearmanWood {
				u := Unit{Job: Military, IdleTime: BuildTimeSpearman}
				b.UnitQueue = append(b.UnitQueue, u)
			} else {
				fmt.Println("Not enough resources for Spearman")
			}
		case "man":
			if TotalFood > CostManAtArmsFood && TotalGold > CostManAtArmsGold {
				u := Unit{Job: Military, IdleTime: BuildTimeManAtArms}
				b.UnitQueue = append(b.UnitQueue, u)
			} else {
				fmt.Println("Not enough resources for Man-At-Arms")
			}
		case "no":
		default:
			fmt.Println("Invalid option")
		}
	} else if b.Type == Stable {
		fmt.Print("Choose a Unit: hor, kni, no: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		switch text {
		case "hor":
			if TotalFood > CostHorsemanFood && TotalWood > CostHorsemanWood {
				u := Unit{Job: Military, IdleTime: BuildTimeHorseman}
				b.UnitQueue = append(b.UnitQueue, u)
			} else {
				fmt.Println("Not enough resources for Horseman")
			}
		case "kni":
			if TotalFood > CostKnightFood && TotalGold > CostKnightGold {
				u := Unit{Job: Military, IdleTime: BuildTimeKnight}
				b.UnitQueue = append(b.UnitQueue, u)
			} else {
				fmt.Println("Not enough resources for Knight")
			}
		case "no":
		default:
			fmt.Println("Invalid option")
		}
	} else if b.Type == ArcheryRange {
		fmt.Print("Choose a Unit: arc, cro, no: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		switch text {
		case "arc":
		case "cro":
		case "no":
		default:
			fmt.Println("Invalid option")
		}
	}
}

func changeJob(u *Unit) {
	done := false
	for !done {
		fmt.Print("Choose Job: s, b, h, w, g, st, unit, build, reallo, no, print, stat: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		var prevJob Job
		prevUnit := *u
		prevJob = prevUnit.Job

		switch text {
		case "s":
			u.Job = Sheep
			u.IdleTime = 1
			done = true
		case "b":
			if !BuiltBerryMill {
				TotalWood -= CostResourcePoint
				u.IdleTime = 5 + BuildTimeResourcePoint
				BuiltBerryMill = true
			} else {
				u.IdleTime = 5
			}
			u.Job = Berry
			done = true
		case "h":
			if !BuiltHuntingMill {
				TotalWood -= CostResourcePoint
				u.IdleTime = 15 + BuildTimeResourcePoint
				BuiltHuntingMill = true
			} else {
				u.IdleTime = 15
			}
			u.Job = Hunting
			done = true
		case "w":
			if !BuiltLumberCamp {
				TotalWood -= CostResourcePoint
				u.IdleTime = 5 + BuildTimeResourcePoint
				BuiltLumberCamp = true
			} else {
				u.IdleTime = 5
			}
			u.Job = Wood
			done = true
		case "g":
			if !BuiltGoldMine {
				TotalWood -= CostResourcePoint
				u.IdleTime = 5 + BuildTimeResourcePoint
				BuiltGoldMine = true
			} else {
				u.IdleTime = 5
			}
			u.Job = Gold
			done = true
		case "st":
			if !BuiltStoneMine {
				TotalWood -= CostResourcePoint
				u.IdleTime = 5 + BuildTimeResourcePoint
				BuiltStoneMine = true
			} else {
				u.IdleTime = 5
			}
			u.Job = Stone
			done = true
		case "unit":
			checkMakeUnit()
		case "build":
			checkMakeBuilding()
		case "reallo":
			reallocateWorkers()
		case "no":
			if u.Job != New {
				done = true
			}
		case "print":
			printBuildOrder()
		case "stat":
			printStatus()
		default:
			fmt.Println("Not Valid")
		}

		if text != "reallo" && text != "no" && text != "print" {
			bo := fmt.Sprintf("Moved from %5s to %5s: %02d:%02d", prevJob, u.Job, GameTime/60, GameTime%60)
			BuildOrder = append(BuildOrder, bo)
		}

	}

}

func gatherResources() {
	// Gather resources if not idle
	for i, u := range Units {
		var ptr = &u
		if u.IdleTime < 1 {
			switch u.Job {
			case Sheep:
				TotalFood += GatherSpeedSheep
			case Berry:
				TotalFood += GatherSpeedBerry
			case Hunting:
				TotalFood += GatherSpeedHunting
			case Wood:
				TotalWood += GatherSpeedWood
			case Gold:
				TotalGold += GatherSpeedGold
			case Stone:
				TotalStone += GatherSpeedStone
			}
		} else {
			ptr.IdleTime = ptr.IdleTime - 1
		}
		Units[i] = *ptr
	}
}

func reallocateWorkers() {
	done := false
	for !done {
		fmt.Print("Reallocate Workers: s, b, h, w, g, st, print, no: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		switch text {
		case "s":
			for i, u := range Units {
				var ptr = &u
				if ptr.Job == Sheep {
					changeJob(ptr)
				}
				Units[i] = *ptr
			}
		case "b":
			for i, u := range Units {
				var ptr = &u
				if ptr.Job == Berry {
					changeJob(ptr)
				}
				Units[i] = *ptr
			}
		case "h":
			for i, u := range Units {
				var ptr = &u
				if ptr.Job == Hunting {
					changeJob(ptr)
				}
				Units[i] = *ptr
			}
		case "w":
			for i, u := range Units {
				var ptr = &u
				if ptr.Job == Wood {
					changeJob(ptr)
				}
				Units[i] = *ptr
			}
		case "g":
			for i, u := range Units {
				var ptr = &u
				if ptr.Job == Gold {
					changeJob(ptr)
				}
				Units[i] = *ptr
			}
		case "st":
			for i, u := range Units {
				var ptr = &u
				if ptr.Job == Stone {
					changeJob(ptr)
				}
				Units[i] = *ptr
			}
		case "print":
			printBuildOrder()
		case "no":
			done = true
		default:
			fmt.Println("Not Valid")
		}
	}
}

func printReport() {
	var foodPerSec float32
	var woodPerSec float32
	var goldPerSec float32
	var stonePerSec float32
	var villsOnFoodCount int
	var villsOnSheepCount int
	var villsOnBerryCount int
	var villsOnHuntCount int
	var villsOnWoodCount int
	var villsOnGoldCount int
	var villsOnStoneCount int
	var militaryCount int
	var builderCount int

	for _, u := range Units {
		if u.Job == Sheep {
			villsOnFoodCount += 1
			villsOnSheepCount += 1
			foodPerSec += GatherSpeedSheep
		}
		if u.Job == Berry {
			villsOnFoodCount += 1
			villsOnBerryCount += 1
			foodPerSec += GatherSpeedBerry
		}
		if u.Job == Hunting {
			villsOnFoodCount += 1
			villsOnHuntCount += 1
			foodPerSec += GatherSpeedHunting
		}
		if u.Job == Wood {
			villsOnWoodCount += 1
			woodPerSec += GatherSpeedWood
		}
		if u.Job == Gold {
			villsOnGoldCount += 1
			goldPerSec += GatherSpeedGold
		}
		if u.Job == Stone {
			villsOnStoneCount += 1
			stonePerSec += GatherSpeedStone
		}
		if u.Job == Military {
			militaryCount += 1
		}
		if u.Job == Builder {
			builderCount += 1
		}
	}

	fmt.Printf("TotalPop  : %5d\n", len(Units))
	fmt.Printf("Military  : %5d\n", militaryCount)
	fmt.Printf("Builders  : %5d\n", builderCount)
	fmt.Printf("TotalFood : %5.0f Vil:%2d PerMin: %5.2f\n", TotalFood, villsOnFoodCount, foodPerSec*60)
	fmt.Printf("TotalWood : %5.0f Vil:%2d PerMin: %5.2f\n", TotalWood, villsOnWoodCount, woodPerSec*60)
	fmt.Printf("TotalGold : %5.0f Vil:%2d PerMin: %5.2f\n", TotalGold, villsOnGoldCount, goldPerSec*60)
	fmt.Printf("TotalStone: %5.0f Vil:%2d PerMin: %5.2f\n", TotalStone, villsOnStoneCount, stonePerSec*60)
	fmt.Println("-----------------------------------------------")
	fmt.Printf("TotalOnSheep: %2d\n", villsOnSheepCount)
	fmt.Printf("TotalOnBerry: %2d\n", villsOnBerryCount)
	fmt.Printf("TotalOnHuntg: %2d\n", villsOnHuntCount)
	fmt.Println()
}

func printBuildOrder() {
	for _, s := range BuildOrder {
		fmt.Println(s)
	}
	printReport()
}

func printStatus() {
	fmt.Println(Units)
	fmt.Println(Buildings)
}
