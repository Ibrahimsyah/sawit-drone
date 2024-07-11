// app.go is intended to contains the core functionality of the app.
// The user entry point goes from here.

package main

import (
	"fmt"
	"os"
)

var (
	// Storing os.Exit to a variable to mock the test so we can simulate the os.Exit behavior.
	// We can not actually observe the os.Exit since it will automatically terminate all
	// process including the test, so the idea is to mock it into a panic
	// then we can catch and assert the value
	osExit = os.Exit
)

//go:generate mockgen -source=app.go -destination=app_mock.go -package=main
type UtilProvider interface {
	Scanln(target ...any)
}

type App struct {
	util UtilProvider
}

func NewApp(util UtilProvider) *App {
	return &App{
		util: util,
	}
}

func (a *App) Start() {
	length, width, count := 0, 0, 0
	a.util.Scanln(&length, &width, &count)
	if valid := a.validateInitialInputs(length, width, count); !valid {
		a.throwFail()
		return
	}

	// treeMap is a map of coordinate x,y to tree height
	treeMap := make(map[string]int)
	for i := 0; i < count; i++ {
		x, y, height := 0, 0, 0
		a.util.Scanln(&x, &y, &height)

		if height < 1 || height > 30 {
			a.throwFail()
			return
		}

		treeKey := a.generateTreeKey(x, y)
		treeMap[treeKey] = height
	}

	distance := a.calculateFlyDistance(length, width, treeMap)
	fmt.Println(distance)
}

// absInt is a method to return absolute integer value of the input
func (a *App) absInt(input int) int {
	if input < 0 {
		return -input
	}
	return input
}

// calculateFlyDistance is the core method to calculate total fly distance of the drone
// on both vertically and horizontally.
//
// It accepts length and width of the field, and the map of tree in the field.
// It returns an integer denoting the distance of the drone
func (a *App) calculateFlyDistance(length, width int, treeMap map[string]int) int {
	// initialize the distance the total horizontal fly distance
	// plus 1 at the beginning as the drone take off
	// and 1 at the end as the drone lands
	distance := 1 + a.calculateHorizontalDistance(length, width)*10 + 1

	// Explore every single plot on the field check whether there is a tree
	// on the current plot. Started from 1, 1
	x, y := 1, 1
	currentAltitude := 1 // The current drone altitude
	for x <= length && y <= width {
		x, y = a.getNextPlotCoordinate(length, x, y)
		key := a.generateTreeKey(x, y)

		// if there is no tree on the plot, decrease the drone altitude to 1
		treeHeight, found := treeMap[key]
		if !found {
			if currentAltitude != 1 {
				deltaAltitude := currentAltitude - 1
				currentAltitude = 1
				distance += a.absInt(deltaAltitude)
			}

			continue
		}

		// adjust the altitude
		deltaAltitude := treeHeight + 1 - currentAltitude
		currentAltitude = treeHeight + 1
		distance += a.absInt(deltaAltitude)
	}

	return distance
}

// calcualteHorizontalDistance calculates the horizontal distance of the drone in the field
// based on the given length and width of the field.
//
// It returns the distance of the drone will make to fly from the bottom-left-most
// to the top-right-most point horizontally
func (a *App) calculateHorizontalDistance(length, width int) int {
	// if the width is only 1, meaning the drone will only fly straight 1 time
	// then we return the result as the drone will not come back
	if width == 1 {
		return length - 1
	}

	// The distance of the drone will be determined by how big the area of the field
	// plus the step where the drone go north.
	// For the odd width, we need to add 1 more step as division by 2 will round the result down
	northSteps := width / 2
	if width%2 != 0 {
		northSteps += 1
	}

	return (length-1)*width + northSteps
}

// getNextPlotCoordinate is a method that will give the next coordinate x and y
// to drone for the next plot.
//
// It accepts the length of the field and also the current x y coordinate.
// It returns the x1 and y1 representing the next drone coordinate
func (a *App) getNextPlotCoordinate(length, x, y int) (x1 int, y1 int) {
	// if the x is 1, meaning it is in the west-most
	// we need to check if the y is odd or even to determine which direction
	// the drone will go (north or east)
	if x == 1 {
		if y%2 == 0 {
			return x, y + 1
		}

		return x + 1, y
	}

	// If the x is equal to length, then we check whether the drone has an odd or even y.
	// If the y is odd then we move the drone to north and keep the x (x, y + 1).
	// If the y is even then we move the drone to west (x - 1, y)
	if x == length {
		if y%2 == 0 {
			return x - 1, y
		}

		return x, y + 1
	}

	if y%2 == 0 {
		return x - 1, y
	}

	return x + 1, y
}

// generateTreeKey is helper method the generate the map key of the tree
// by the given coordinate x and y.
//
// It will return "x,y" as string
func (a *App) generateTreeKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// throwFail throws an error "FAIL" to stderr and then exits the program with status 1
func (a *App) throwFail() {
	fmt.Fprintln(os.Stderr, "FAIL")
	osExit(1)
}

// validateInitialInputs is a helper method to validate the width, length, and count given by user.
//
// It returns true if all the input valid.
// Otherwise, it returns false as there are some invalid input.
func (a *App) validateInitialInputs(length, width, count int) bool {
	return width >= 1 && width <= 50000 &&
		length >= 1 && length <= 50000 &&
		count >= 1 && count <= 50000
}
