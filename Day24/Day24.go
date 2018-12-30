package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strings"
	"sort"
)

func getMillis() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func check (e error) {
	if(e != nil) {
		panic(e);
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

type armyunit struct {
	dead bool
	hp int
}

type armygroup struct {
	units []armyunit
	initiative int
	weaknesses []string
	immunities []string
	attack string
	ap int
}

type army struct {
	is_infection bool
	groups []armygroup
}

type intTuple struct {
	a int
	b int
} 

func (a army) numGroups() int {
	count := 0
	for i:=0; i < len(a.groups); i++ {
		if(len(a.groups[i].units) == 0) {
			continue
		}
		count++
	}
	return count
}

func (ag *armygroup) defend(damage int)int {
	numUnits := ag.numUnits()
	if numUnits == 0 {
		return 0
	}
	deadUnits := damage/ag.units[0].hp
	if deadUnits > numUnits {
		deadUnits = numUnits
		ag.units = []armyunit{}
	} else {
		ag.units = ag.units[0:numUnits-deadUnits]
	}
	return deadUnits
}
func (ag armygroup) numUnits() int {
	num_units := 0
	for i:=0; i < len(ag.units); i++ {
		if(!ag.units[i].dead && ag.units[i].hp > 0) {
			num_units++
		}
	}
	return num_units
}

func (ag armygroup) EP() int {
	num_units := 0
	for i:=0; i < len(ag.units); i++ {
		if(!ag.units[i].dead && ag.units[i].hp > 0) {
			num_units++
		}
	}
	return num_units * ag.ap
}

func parseInput(lines []string)(army,army) {
	var infection army
	var immunesystem army
	infection.is_infection = true
	
	onImmune := false
	for i:=0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			continue
		}
		if lines[i] == "Immune System:" {
			onImmune = true
			continue
		}
		if lines[i] == "Infection:"  {
			onImmune = false
			continue
		}
		num_units := 0
		hp := 0
		initiative := 0
		ap := 0
		damage := ""
		pt1 := ""
		pt2 := ""
		pt3 := ""
		if strings.Index(lines[i], "(") > -1 {
			pt1 = lines[i][0:strings.Index(lines[i], "(")]
			pt2 = lines[i][strings.Index(lines[i], "(")+1:strings.Index(lines[i], ")")]
			pt3 = lines[i][strings.Index(lines[i], ")")+1:len(lines[i])]
		} else {
			pt1 = lines[i][0:strings.Index(lines[i]," with an attack")]
			pt3 = lines[i][strings.Index(lines[i]," with an attack"):len(lines[i])]
		}
		
		fmt.Sscanf(pt1, "%d units each with %d hit points", &num_units, &hp)
		fmt.Sscanf(pt3, " with an attack that does %d %s damage at initiative %d", &ap, &damage, &initiative)
		
		group := armygroup{initiative:initiative, attack:damage, ap:ap}
		if(pt2 != "") {
			parts := strings.Split(pt2, "; ")
			for _,p := range parts {
				if strings.Index(p, "immune") > -1 {
					p1 := p[len("immune to "):len(p)]
					p1p := strings.Split(p1, ", ")
					for _,p2 := range p1p {
						group.immunities = append(group.immunities, p2)
					}
				} else {
					p1 := p[len("weak to "):len(p)]
					p1p := strings.Split(p1, ", ")
					for _,p2 := range p1p {
						group.weaknesses = append(group.weaknesses, p2)
					}
				}
			}
		}
		for j:=0; j<num_units;j++ {
			group.units = append(group.units, armyunit{dead:false, hp:hp})
		}
		
		if(onImmune) {
			immunesystem.groups = append(immunesystem.groups, group)
		} else {
			infection.groups = append(infection.groups,group)
		}
		
	}
	
	return immunesystem,infection
}

func targetSelectionPriority(ar0 *army, ar1 *army) []intTuple {
	var order []intTuple
	for len(order) < ar0.numGroups() + ar1.numGroups() {
		// find highest priority group not already in order
		high_ep := 0
		high_init := 0
		high_index := 0
		high_army := 0
		for i:=0; i < len(ar0.groups); i++ {
			if(ar0.groups[i].numUnits() == 0) {
				continue
			}			
			found := false
			
			for j:=0; j < len(order); j++ {
				if i == order[j].b && order[j].a == 0 {
					found = true
					break
				}
			}
			if found {
				continue
			}
			if ar0.groups[i].EP() > high_ep || (ar0.groups[i].EP() == high_ep && ar0.groups[i].initiative > high_init) {
				high_index = i
				high_ep = ar0.groups[i].EP()
				high_init = ar0.groups[i].initiative
				high_army = 0
			}
		}
		for i:=0; i < len(ar1.groups); i++ {
			if(ar1.groups[i].numUnits() == 0) {
				continue
			}			
			found := false
			
			for j:=0; j < len(order); j++ {
				if i == order[j].b && order[j].a == 1 {
					found = true
					break
				}
			}
			if found {
				continue
			}
			if ar1.groups[i].EP() > high_ep || (ar1.groups[i].EP() == high_ep && ar1.groups[i].initiative > high_init) {
				high_index = i
				high_ep = ar1.groups[i].EP()
				high_init = ar1.groups[i].initiative
				high_army = 1
			}
		}
		order = append(order,intTuple{a:high_army,b:high_index})
	}
	return order
}

func calcDamage(src *army, dest *army, srcIndex int, destIndex int) int {
	ep := src.groups[srcIndex].EP()
	attack := src.groups[srcIndex].attack
	is_weak := false
	is_immune := false
	for i := 0; i < len(dest.groups[destIndex].weaknesses); i++ {
		if attack == dest.groups[destIndex].weaknesses[i] {
			is_weak = true
			break
		}
	}
	if is_weak {
		return ep * 2
	}
	for i := 0; i < len(dest.groups[destIndex].immunities); i++ {
		if attack == dest.groups[destIndex].immunities[i] {
			is_immune = true
			break
		}
	}
	if is_immune {
		return 0
	}
	return ep	
}

func initiativeOrder(immunesystem *army, infection *army) []intTuple {
	var order []intTuple
	var inits []int
	
	for i:=0; i < len(immunesystem.groups); i++ {
		if immunesystem.groups[i].numUnits() == 0 {
			continue
		}
		inits = append(inits, immunesystem.groups[i].initiative)
	}
	for i:=0; i < len(infection.groups); i++ {
		if infection.groups[i].numUnits() == 0 {
			continue
		}
		inits = append(inits, infection.groups[i].initiative)
	}
	sort.Ints(inits)
	for i:=len(inits) - 1; i >= 0; i-- {
		found := false
		for j:=0; j < len(immunesystem.groups); j++ {
			if immunesystem.groups[j].initiative == inits[i] {
				order = append(order, intTuple{a:0, b:j})
				found = true
				break
			}
		}
		if(found) {
			continue
		}
		for j:=0; j < len(infection.groups); j++ {
			if infection.groups[j].initiative == inits[i] {
				order = append(order, intTuple{a:1, b:j})
				break
			}
		}
	}
	return order
}

func fight(immunesystem *army, infection *army, verbose bool) {
	// target selection
	priority := targetSelectionPriority(immunesystem, infection)
	var immuneTargets []intTuple
	var infectionTargets []intTuple
	
	if(verbose) {
		fmt.Println("Immune system:")
		for i:=0; i < len(immunesystem.groups); i++ {
			fmt.Println("Group", i+1, immunesystem.groups[i].numUnits())
		}
		fmt.Println("Infections:")
		for i:=0; i < len(infection.groups); i++ {
			fmt.Println("Group", i+1, infection.groups[i].numUnits())
		}
	}
	
	for _,v := range priority {
		high_damage:=0
		high_index:=-1
		high_ep:=0
		high_init:=0
		if v.a == 0 {
			for i:= 0; i < len(infection.groups); i++ {
				if infection.groups[i].numUnits() == 0 {
					continue
				}
				damage := calcDamage(immunesystem, infection, v.b, i)
				if(damage == 0) {
					continue
				}
				if damage > high_damage || (damage == high_damage && infection.groups[i].EP() > high_ep) || (damage == high_damage && infection.groups[i].EP() == high_ep && infection.groups[i].initiative > high_init) {
					found := false 
					for j:=0; j < len(immuneTargets); j++ {
						if immuneTargets[j].b == i {
							found = true
							break
						}
						
					}
					if found {
						continue
					}
					high_index = i
					high_damage = damage
					high_ep = infection.groups[i].EP()
					high_init = infection.groups[i].initiative
				}
			}
			
		} else {
			for i:= 0; i < len(immunesystem.groups); i++ {
				if immunesystem.groups[i].numUnits() == 0 {
					continue
				}
				damage := calcDamage(infection, immunesystem, v.b, i)
				if(damage == 0) {
					continue
				}
				if damage > high_damage || (damage == high_damage && immunesystem.groups[i].EP() > high_ep) || (damage == high_damage && immunesystem.groups[i].EP() == high_ep && immunesystem.groups[i].initiative > high_init) {
					found := false 
					for j:=0; j < len(infectionTargets); j++ {
						if infectionTargets[j].b == i {
							found = true
							break
						}
						
					}
					if found {
						continue
					}
					high_index = i
					high_damage = damage
					high_ep = immunesystem.groups[i].EP()
					high_init = immunesystem.groups[i].initiative
				}
			}
		}
		
		if high_damage > 0 {
			if v.a == 0 {
				immuneTargets = append(immuneTargets, intTuple{v.b,high_index})
			} else {
				infectionTargets = append(infectionTargets, intTuple{v.b,high_index})
			}
		}
	}
	
	// fight
	initiativeOrder := initiativeOrder(immunesystem, infection)
	for _,v := range initiativeOrder {
		if v.a == 0 {
			target := -1
			for j:=0; j < len(immuneTargets); j++ {
				if immuneTargets[j].a == v.b {
					target = immuneTargets[j].b
					break
				}
			}
			
			if target > -1 {
				damage := calcDamage(immunesystem, infection, v.b, target)
				killed := infection.groups[target].defend(damage)
				if(verbose) {
					fmt.Printf("Immune System group %d attacks defending group %d for %d damage, killing %d units\n", v.b + 1, damage, target + 1, killed)
				}
			} else {
				if(verbose) {
					fmt.Printf("Immune System group %d does not attack\n", v.b + 1)
				}
			}
			
		} else {
			target := -1
			for j:=0; j < len(infectionTargets); j++ {
				if infectionTargets[j].a == v.b {
					target = infectionTargets[j].b
					break
				}
			}
			
			if target > -1 {
				damage := calcDamage(infection, immunesystem, v.b, target)
				killed := immunesystem.groups[target].defend(damage)
				if(verbose) {
					fmt.Printf("Infection group %d attacks defending group %d for %d damage, killing %d units\n", v.b + 1, damage, target + 1, killed)
				}
			} else {
				if(verbose) {
					fmt.Printf("Infection group %d does not attack\n", v.b + 1)
				}
			}
		}
	}
}

func boostArmy(a *army, boost int) {
	for i:=0; i < len(a.groups); i++ {
		a.groups[i].ap += boost
	}
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// Parse input
	immunesystem, infection := parseInput(lines)
	
	turn := 0
	sum_infection := 0
	sum_immune := 0
	verbose := false
	for {
		if(verbose) {
			fmt.Printf("### TURN %d ###\n", turn + 1)
		}
		fight(&immunesystem, &infection, verbose)
		turn++
		sum_infection = 0
		sum_immune = 0
		for i := 0; i < len(immunesystem.groups); i++ {
			sum_immune += immunesystem.groups[i].numUnits()
		}
		for i := 0; i < len(infection.groups); i++ {
			sum_infection += infection.groups[i].numUnits()
		}
		if(sum_immune == 0 || sum_infection == 0) {
			break
		}
	}
	
	fmt.Println("Result A:", sum_immune + sum_infection)
	
	
	boost := 1
	turn = 0
	verbose = false
	
	for {
		immunesystem, infection = parseInput(lines)
		boostArmy(&immunesystem, boost)
		last_sum_immune := 0
		last_sum_infection := 0
		stalemate := false
		for {
			if(verbose) {
				fmt.Printf("### TURN %d ###\n", turn + 1)
			}
			fight(&immunesystem, &infection, verbose)
			turn++
			sum_infection = 0
			sum_immune = 0
			for i := 0; i < len(immunesystem.groups); i++ {
				sum_immune += immunesystem.groups[i].numUnits()
			}
			for i := 0; i < len(infection.groups); i++ {
				sum_infection += infection.groups[i].numUnits()
			}
			if(sum_immune == 0 || sum_infection == 0) {
				break
			}
			if(last_sum_immune == sum_immune && last_sum_infection == sum_infection) {
				stalemate = true
				break
			}
			last_sum_immune = sum_immune
			last_sum_infection = sum_infection
		}
		if( sum_immune > 0 && !stalemate) {
			break
		}
		
		if(sum_infection > 10000) {
			boost+=10
		} else {
			boost++
		}
	}
	
	fmt.Println("Result B:", sum_immune)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}