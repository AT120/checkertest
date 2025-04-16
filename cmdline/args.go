package cmdline

import (
	"fmt"
	"os"
)

type CheckerArgs struct {
	TestFile   string
	AnswerFile string
}

func (c *CheckerArgs) fieldToFlag(flag string) *string {
	switch flag {
	case "--test":
		return &c.TestFile
	case "--answer":
		return &c.AnswerFile
	default:
		return nil
	}
}

func PrintUsage() {
	fmt.Println("usage: checker --test /path/to/testfile --answer /path/to/answer")
}

func ParseCmdlineArgs() (*CheckerArgs, error) {
	if len(os.Args) < 5 {
		return nil, fmt.Errorf("not enough arguments")
	}

	args := CheckerArgs{}

	for i := 1; i < len(os.Args); i += 2 {
		field := args.fieldToFlag(os.Args[i])
		if field == nil {
			return nil, fmt.Errorf("invalid flag: %v", os.Args[i])
		}

		*field = os.Args[i+1] //TODO: смешно - смеемся
	}

	return &args, nil
}
