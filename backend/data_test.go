package main

import "testing"

func TestCalcDensity(t *testing.T) {
	tests := []struct {
		name       string
		population int64
		area       float64
		want       float64
	}{
		{name: "regular", population: 13250000, area: 2561.5, want: 5172.75},
		{name: "zero area", population: 1000, area: 0, want: 0},
		{name: "negative area", population: 1000, area: -5, want: 0},
		{name: "negative population", population: -1, area: 10, want: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calcDensity(test.population, test.area); got != test.want {
				t.Fatalf("calcDensity(%d, %v) = %v, want %v", test.population, test.area, got, test.want)
			}
		})
	}
}

func TestDatasetIntegrity(t *testing.T) {
	dataset := getDataset()
	if len(dataset) == 0 {
		t.Fatal("dataset is empty")
	}

	ids := make(map[string]struct{}, len(dataset))
	for _, entity := range dataset {
		if entity.ID == "" || entity.Name == "" {
			t.Fatalf("entity has empty id or name: %+v", entity)
		}
		if _, exists := ids[entity.ID]; exists {
			t.Fatalf("duplicate entity id: %q", entity.ID)
		}
		ids[entity.ID] = struct{}{}

		if entity.Type != "country" && entity.Type != "city" {
			t.Fatalf("entity %q has invalid type %q", entity.ID, entity.Type)
		}
		if len(entity.History) < 2 {
			t.Fatalf("entity %q has too few observations", entity.ID)
		}

		previousYear := 0
		for _, stat := range entity.History {
			if stat.Year <= previousYear {
				t.Fatalf("entity %q history is not sorted by year", entity.ID)
			}
			previousYear = stat.Year
			if stat.Population <= 0 || stat.AreaKm2 <= 0 {
				t.Fatalf("entity %q has non-positive population or area in %d", entity.ID, stat.Year)
			}
			if want := calcDensity(stat.Population, stat.AreaKm2); stat.Density != want {
				t.Fatalf("entity %q density in %d = %v, want %v", entity.ID, stat.Year, stat.Density, want)
			}
		}
	}
}
