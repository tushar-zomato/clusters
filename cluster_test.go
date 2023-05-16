package clusters

import (
	"testing"
)

func TestDistance(t *testing.T) {
	p1 := Coordinates{2, 2}
	p2 := Coordinates{3, 5}

	d := p1.Distance(p2)
	if d != 10 {
		t.Errorf("Expected distance of 10, got %f", d)
	}
}

func TestCenter(t *testing.T) {
	var o Observations
	o = append(o, WeightedObservation{c: Coordinates{1, 1}, weight: 5})
	o = append(o, WeightedObservation{c: Coordinates{3, 2}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{5, 3}, weight: 10})

	m, err := o.Center()
	if err != nil {
		t.Errorf("Could not retrieve center: %v", err)
		return
	}

	if m[0] != 3.625 || m[1] != 2.3125 {
		t.Errorf("Expected coordinates [3.625 2.3125], got %v", m)
	}
}

func TestAverageDistance(t *testing.T) {
	var o Observations
	o = append(o, WeightedObservation{c: Coordinates{1, 1}, weight: 5})
	o = append(o, WeightedObservation{c: Coordinates{3, 2}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{5, 3}, weight: 9})

	d := AverageDistance(o[0], o[1:])
	if d != 18.5 {
		t.Errorf("Expected average distance of 18.5, got %v", d)
	}

	d = AverageDistance(o[1], Observations{o[1]})
	if d != 0 {
		t.Errorf("Expected average distance of 0, got %v", d)
	}
}

func TestClusters(t *testing.T) {
	var o Observations
	o = append(o, WeightedObservation{c: Coordinates{1, 1}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{3, 2}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{5, 3}, weight: 1})

	c, err := New(2, o)
	if err != nil {
		t.Errorf("Error seeding clusters: %v", err)
		return
	}

	if len(c) != 2 {
		t.Errorf("Expected 2 clusters, got %d", len(c))
		return
	}

	c[0].Append(o[0])
	c[1].Append(o[1])
	c[1].Append(o[2])
	c.Recenter()

	if n := c.Nearest(o[1]); n != 1 {
		t.Errorf("Expected nearest cluster 1, got %d", n)
	}

	nc, d := c.Neighbour(o[0], 0)
	if nc != 1 {
		t.Errorf("Expected neighbouring cluster 1, got %d", nc)
	}

	if d != 12.5 {
		t.Errorf("Expected neighbouring cluster distance 12.5, got %f", d)
	}

	if pp := c[1].PointsInDimension(0); pp[0] != 3 || pp[1] != 5 {
		t.Errorf("Expected [3 5] as points in dimension 0, got %v", pp)
	}

	if pp := c.CentersInDimension(0); pp[0] != 1 || pp[1] != 4 {
		t.Errorf("Expected [1 4] as centers in dimension 0, got %v", pp)
	}

	c.Reset()
	if len(c[0].Observations) > 0 {
		t.Errorf("Expected empty cluster 1, found %d observations", len(c[0].Observations))
	}
}

func TestFarApartClusters(t *testing.T) {
	var o Observations
	o = append(o, WeightedObservation{c: Coordinates{1, 1}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{3, 2}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{5, 3}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{10, 3}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{5, 19}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{21, 13}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{8, 9}, weight: 1})
	o = append(o, WeightedObservation{c: Coordinates{7, 1}, weight: 1})

	cFADistTotal := 0.0
	cDistTotal := 0.0
	for i := 0; i < 50; i++ {
		cFA, _ := NewFarApart(4, o)
		for j := 1; j < 4; j++ {
			cFADistTotal += cFA[j-1].Center.Distance(cFA[j].Center)
		}

		c, _ := New(4, o)
		for j := 1; j < 4; j++ {
			cDistTotal += c[j-1].Center.Distance(c[j].Center)
		}
	}

	if cDistTotal > cFADistTotal {
		t.Errorf("Expected far apart distance to be more, found FA: %f Regular: %f", cFADistTotal, cDistTotal)
	}
}
