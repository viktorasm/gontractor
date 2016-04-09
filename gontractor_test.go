package main

import "testing"

func TestE2E(t *testing.T) {
	g := NewGontractor()
	g.spec = "./test-resources/swagger.yaml"
	g.outDir = "test-output/testE2E"
	g.Execute()
}
