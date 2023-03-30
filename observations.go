package clusters

import (
	"fmt"
	"math"
)

// Coordinates is a slice of float64
type Coordinates []float64

// Distance returns the euclidean distance between two coordinates
func (c Coordinates) Distance(p2 Coordinates) float64 {
	var r float64
	for i, v := range c {
		r += math.Pow(v-p2[i], 2)
	}

	return r
}

// Observation is a data point (float64 between 0.0 and 1.0) in n dimensions
type Observation interface {
	Coordinates() Coordinates
	Weight() int
}

// Observations is a slice of observations
type Observations []Observation

// WeightedObservation implements Observation
type WeightedObservation struct {
	c      Coordinates
	weight int
}

func NewObservation(c Coordinates, weight int) WeightedObservation {
	if weight <= 0 {
		weight = 1
	}

	return WeightedObservation{
		c:      c,
		weight: weight,
	}
}

// Coordinates implements the Observation interface for a plain set of float64
// coordinates
func (o WeightedObservation) Coordinates() Coordinates {
	return o.c
}

func (o WeightedObservation) Weight() int {
	if o.weight <= 0 {
		return 1
	}

	return o.weight
}

// Center returns the center coordinates of a set of Observations
func (c Observations) Center() (Coordinates, error) {
	if len(c) == 0 {
		return nil, fmt.Errorf("there is no mean for an empty set of points")
	}

	cc := make([]float64, len(c[0].Coordinates()))
	totalPoints := 0
	for _, point := range c {
		for j, v := range point.Coordinates() {
			cc[j] += v * float64(point.Weight())
		}

		totalPoints += point.Weight()
	}

	var mean Coordinates
	for _, v := range cc {
		mean = append(mean, v/float64(totalPoints))
	}

	return mean, nil
}

// AverageDistance returns the average distance between o and all observations
func AverageDistance(o Observation, observations Observations) float64 {
	var d float64
	var l int

	for _, observation := range observations {
		dist := o.Coordinates().Distance(observation.Coordinates())
		if dist == 0 {
			continue
		}

		l += observation.Weight()
		d += dist * float64(observation.Weight())
	}

	if l == 0 {
		return 0
	}

	return d / float64(l)
}
