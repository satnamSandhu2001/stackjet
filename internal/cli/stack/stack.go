package stack

import (
	"slices"

	"github.com/satnamSandhu2001/stackjet/internal/cli/nodejs"
	"github.com/satnamSandhu2001/stackjet/pkg"
)

func IsValidStack(stack string) bool {
	return slices.Contains(pkg.Config().VALID_STACKS, stack)
}

func DeployStack(stack string) error {
	switch stack {
	case "nodejs":
		return nodejs.DeployStack()
	}

	return nil
}
