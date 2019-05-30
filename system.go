package system

/* This package implements the ADT (abstract data type) "system". Always use func "system.New ()"
	when creating a new data of this type. */

import (
	"errors"
	"fmt"
	"github.com/qamarian-etc/slices"
	"strings"
)

func New () (*System) { // Creates a new system.
	return &System {[]string {}, map[string][]string {}, ""}
}

type System struct {
	systemElements []string // All elements in the system.

	dependencies map[string][]string /* The dependencies of individual elements in the system.
		The list of dependencies of each individual element, would be stored in this hasp
		map, where the key of each record would be the ID of the element. */

	addedElements string /* A string that keeps track of what elements have been added to the
		system. It is just a redundant data meant to help speed up some certain operations
		of this data type. */
}

func (someSystem *System) AddElement (newElement string, dependencies []string) (error) { /* Adds
	an element to the system.

	INPUT
	input 0: The new element to be added to the system.

	OUTPT
	outpt 0: If the element was successfully added, value would be nil error. But, if the
		element has already been added to the system, value would be error
		"system.ErrAlreadyAdded". */

	if strings.Contains (someSystem.addedElements, newElement + "/") == true {
		return ErrAlreadyAdded
	} else {
		someSystem.systemElements = append (someSystem.systemElements, newElement)
		someSystem.dependencies [newElement] = dependencies
		someSystem.addedElements += newElement + "/"
		return nil
	}
}

func (someSystem *System) InitOrder () ([]string, error, string) { /* This functions provides an
	order in which elements of this system can be safely initialized.

	OUTPT
	outpt 0: A string slice of the IDs of the elements in the system. The order of these IDs
		represents the "init order". If an error is encountered during the operation, this
		data would be nil.

	outpt 1: If operation succeeds, value would be nil. Otherwise, value would be the error
		that occured.

	outpt 2: When the value of outpt 1 is an error, value of this data would be a more precise
		description of the error. Possible values look like the following:

		"Dependency 'x' is missing" - Value of outpt 2 when the dependency of an element
			is not in the system.

		"Element 'r' is part of the circle"- Value of outpt 2 when a cyclic
			dependency is detected. */

	// Declaration of some data to be used for this operation. { ...
	elements := someSystem.systemElements
	initOrder := []string {}
	waitingList := []string {}
	// ... }

	/* The elements of this system are popped one-by-one, and added in an appropriate place, in
		the "init order" that is being generated. */
	for {
		elementUnderProcessing := elements [0]
		var errX error = nil
		var errDescp string
		initOrder, elements, errX, errDescp = addToInitOrder (initOrder,
			elementUnderProcessing, waitingList, elements, someSystem) /* Adding an
				element to the "init order" */
		if errX != nil {
			return nil, errX, errDescp
		}
		if len (elements) == 0 {
			break
		}
	}

	return initOrder, nil, ""
}

func addToInitOrder (initOrder []string, element string, waitingList []string,
		elements []string, someSystem *System) ([]string, []string, error, string) { /* This
	function is not meant to be used outside this package. The function simply takes an init
	order and an element, then adds the element to a safe place in the "init order".

	INPUT
	input 0: The init order where the element should be added.
	input 1: The element to be added.
	input 2: You may need to read the code to fully grasp the essence of this data. This data is
		a stack. When an element needs to be added to the init order, but has dependencies,
		the element is placed in this waiting list, and we try to add the dependencies to the
		init order first. Once the dependencies have been added to the init order, the
		element can then be popped from this stack and added to the init order.
	input 3: The system whose's init order is being worked on.

	OUTPT
	outpt 0: A modified version of the init order. If this operation fails, the value of this
		data would be nil.
	outpt 1: If this operation succeeds, value of this data would be nil error. If this operation
		should fail, value of this data would be an error.
	ouptt 2: If this operation succeeds, value of this data would be an empty string. If this
		operation should fail, value of this data would be a more precise description of the
		error. */

	// Checking for existence of a circle.
	if slices.IsElementInStringSlice (waitingList, element) == true {
		return nil, nil, ErrCircleDetected, "Element '" + element + "' is part of the circle"
	}

	// If the element has no dependency, it is added to the init order, straightaway.
	if len (someSystem.dependencies [element]) == 0 {
		elements = slices.RemoveFromStringSlice (elements, element)
		initOrder := append (initOrder, element)
		return initOrder, elements, nil, ""
	}

	// If the element has any dependency, the dependencies are added first. { ...
	waitingList = append (waitingList, element) /* Placing the element in the waiting list.
		Once all its dependencies have been added to the init order, it would be removed
		from this waiting list. */

	// Adding dependencies to the "init order".
	for _, dependency := range someSystem.dependencies [element] {
		// If the dependency is already in the "init order", there is no need readding it.
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
		initOrder, elements, errZ, errDescp = addToInitOrder (initOrder, dependency,
			waitingList, elements, someSystem)
		// ... }

		if errZ != nil {
			return nil, nil, errZ, errDescp
		}
	}

	/* At this stage all dependencies of the element must have been added to the init order. Now,
		the element will be removed from the waiting list, and added to the init order. */
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
