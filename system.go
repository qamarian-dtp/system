package system

import (
	"errors"
	"fmt"
	"gopkg.in/qamarian-etc/slices.v1"
	"strings"
)

func New () (*System) { // Creates a new system.
	return &System {[]string {}, map[string][]string {}, ""}
}

type System struct {
	systemElements []string // All elements in the system.
	dependencies map[string][]string /* The dependencies of individual elements in the
		system. The list of dependencies of each individual element, would be
		stored in this hash map, where the key of each record would be the ID of
		the element. */
	addedElements string /* A string that keeps track of what elements have been added
		to the system. It is just a redundant data meant to help speed up some
		certain operations of this data type. */
}

// Adds an element to the system.
//
// Inputs
//
// input 0: The new element to be added to the system. Value can not be an empty string.
//
// input 1: The IDs of the dependencies of the element. The ID of a dependency may not be
// an empty string.
//
// Outpts
//
// outpt 0: Possible errors include: ErrAlreadyAdded.
func (someSystem *System) AddElement (newElement string, dependencies []string) (error) {

	if newElement == "" {
		return errors.New ("Empty string can not be used as ID of an element.")
	}
	for _, dep := range dependencies {
		if dep == "" {
			return errors.New ("The ID of a dependency is an empty string.")
		}
	}
	if strings.Contains (someSystem.addedElements, newElement + "/") == true {
		return ErrAlreadyAdded
	}
	someSystem.systemElements = append (someSystem.systemElements, newElement)
	someSystem.dependencies [newElement] = dependencies
	someSystem.addedElements += newElement + "/"
	return nil
}

// This functions provides an order in which elements of the system could be safely
// initialized.
//
// Outpts
// outpt 0: A string slice of the IDs of the elements in the system. The ascending order
// of these IDs represents the "init order". If an error is encountered during the
// operation, value of this data would be nil.
//
// outpt 1: If operation succeeds, value would be nil. Otherwise, value would be the error
// that occured.
//
// outpt 2: When the value of outpt 1 is an error, value of this data would be a more
// precise description of the error. Possible values would look like the following:
//
// "Dependency 'x' is missing" - Value of outpt 2 when a dependency of an element is not
// in the system.
//
// "Element 'r' is part of the circle"- Value of outpt 2 when a cyclic dependency is
// detected.
func (someSystem *System) InitOrder () ([]string, error, string) {

	// Declaration of some data to be used for this operation. { ...
	elements := someSystem.systemElements
	initOrder := []string {}
	waitingList := []string {}
	// ... }

	/* The elements of this system are popped one-by-one, and added in an appropriate
		place, in the "init order" that is being generated. */
	for {
		if len (elements) == 0 {
			break
		}

		elementUnderProcessing := elements [0]
		var errX error = nil
		var errDescp string
		initOrder, elements, errX, errDescp = addToInitOrder (initOrder,
			elementUnderProcessing, waitingList, elements, someSystem)
		if errX != nil {
			return nil, errX, errDescp
		}
	}

	return initOrder, nil, ""
}

func addToInitOrder (initOrder []string, element string, waitingList []string,
		elements []string, someSystem *System) ([]string, []string, error,
		string) { /* This function is not meant to be used outside this package.
		The function simply takes an init order and an element, then adds the
		element to a safe place in the "init order".

	Inputs
	input 0: The init order where the element should be added.
	input 1: The element to be added.
	input 2: You may need to read the code to fully grasp the essence of this data.
		This data is a stack. When an element needs to be added to the init order,
		but has dependencies, the element is placed in this waiting list, and we
		try to add the dependencies to the init order first. Once the dependencies
		have been added to the init order, the element can then be popped from
		this stack and added to the init order.
	input 3: The system whose's init order is being worked on.

	Outpts
	outpt 0: A modified version of the init order. If this operation fails, the value
		of this data would be nil.
	outpt 1: If this operation succeeds, value of this data would be nil error. If
 		this operation should fail, value of this data would be an error.
	outpt 2: If this operation succeeds, value of this data would be an empty string.
		If this operation should fail, value of this data would be a more precise
		description of the error. */

	// Checking for existence of a circle.
	if slices.IsElementInStringSlice (waitingList, element) == true {
		return nil, nil, ErrCircleDetected, "Element '" + element +
			"' is part of the circle."
	}

	// If the element has no dependency, it is added to the init order, straightaway.
	if len (someSystem.dependencies [element]) == 0 {
		elements = slices.RemoveFromStringSlice (elements, element)
		initOrder := append (initOrder, element)
		return initOrder, elements, nil, ""
	}

	// If the element has any dependency, the dependencies are added first. { ...
	waitingList = append (waitingList, element) /* Placing the element in the waiting
		list. Once all its dependencies have been added to the init order, it
 		would be removed from this waiting list. */

	// Adding dependencies to the "init order".
	for _, dependency := range someSystem.dependencies [element] {
		/* If the dependency is already in the "init order", there is no need
			reading it. */
		if slices.IsElementInStringSlice (initOrder, dependency) == true {
			continue
		}

		// If dependency is not in the system, error is returned.
		if slices.IndexInStringSlice (elements, dependency) == -1 {
			return nil, nil, ErrElementMissing, fmt.Sprintf (
				"Dependency '%s' is missing", dependency)
		}

		// Adding dependency to the "init order". { ...
		var errZ error = nil
		var errDescp string
		initOrder, elements, errZ, errDescp = addToInitOrder (initOrder,
			dependency, waitingList, elements, someSystem)
		// ... }

		if errZ != nil {
			return nil, nil, errZ, errDescp
		}
	}

	/* At this stage all dependencies of the element must have been added to the init
		order. Now, the element will be removed from the waiting list, and added
		to the init order. */
	waitingList = slices.RemoveFromStringSlice (waitingList, element)
	elements = slices.RemoveFromStringSlice (elements, element)
	initOrder = append (initOrder, element)
	// ... }

	return initOrder, elements, nil, ""
}

var (
	ErrAlreadyAdded error = errors.New ("The element has already been added")
	ErrCircleDetected error = errors.New ("A circle has been detected")
	ErrElementMissing error = errors.New ("An element is missing")
)
