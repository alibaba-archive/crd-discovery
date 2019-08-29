package utils

import (
	"testing"
)

func TestLoadKubeConfigOrDie(t *testing.T) {
	_ = LoadKubeConfigOrDie()
}

