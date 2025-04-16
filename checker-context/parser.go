package checkercontext

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ExecutorInfo struct {
	Name string
	Args any
}

type SectionInfo struct {
	Id        string
	Executors ExecutorCollection
}

type ExecutorCollection []ExecutorInfo

type SectionCollection []SectionInfo

func (sc *SectionCollection) UnmarshalYAML(value *yaml.Node) error {
	for i := 1; i < len(value.Content); i += 2 {
		if value.Content[i-1].Tag != "!!str" {
			return fmt.Errorf("expected section name as !!str, but instead got %v", value.Content[i-1].Tag)
		}
		if value.Content[i].Tag != "!!seq" {
			return fmt.Errorf("expected executors as !!seq, but instead got %v", value.Content[i-1].Tag)
		}

		*sc = append(*sc, SectionInfo{Id: value.Content[i-1].Value})
		err := (*sc)[len(*sc)-1].Executors.UnmarshalYAML(value.Content[i])

		if err != nil {
			return err
		}
	}

	return nil
}

func (ec *ExecutorCollection) UnmarshalYAML(value *yaml.Node) error {
	for _, node := range value.Content {
		if len(node.Content) < 1 {
			return fmt.Errorf("expected executor name")
		}
		if node.Content[0].Tag != "!!str" {
			return fmt.Errorf("expected executor name as !!str, but instead got %v", node.Content[0].Tag)
		}

		executorName := node.Content[0].Value

		*ec = append(*ec, ExecutorInfo{
			Name: executorName,
		})

		if len(node.Content) < 2 {
			continue
		}

		argsFactory, ok := context.argsFactory[executorName]
		if !ok {
			return fmt.Errorf("executor (%v) was not properly registered", executorName)
		}

		args := argsFactory()

		err := node.Content[1].Decode(args)
		if err != nil {
			return err
		}

		(*ec)[len(*ec)-1].Args = args
	}

	return nil
}
