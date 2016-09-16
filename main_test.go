package main

import "testing"

func TestParseSvStat(t *testing.T) {
	for in, result := range map[string][]float64{
		"up (pid 31420) 81801 seconds": []float64{
			1.0, 1.0, 81801.0,
		},
		"down (signal SIGTERM) 555 seconds, normally up, want up, ready 555 seconds": []float64{
			0.0, 1.0, 555.0,
		},
		"down (exitcode 0) 0 seconds, normally up, want up, ready 0 seconds": []float64{
			0.0, 1.0, 0.0,
		},
		"up (pid 23776) 1 seconds": []float64{
			1.0, 1.0, 1.0,
		},
		// This is made up, no idea how this want actually works..
		"up (pid 23776) 1234 seconds, normally up, want down, ready 1234 seconds": []float64{
			1.0, 0.0, 1234.0,
		},
	} {
		up, want, sc, err := parseSvStat(in)
		if err != nil {
			t.Fatal(err)
		}
		if up != result[0] {
			t.Fatalf("%f != %f", up, result[0])
		}
		if want != result[1] {
			t.Fatalf("%f != %f", want, result[1])
		}
		if sc != result[2] {
			t.Fatalf("%f != %f", sc, result[2])
		}
	}
}
