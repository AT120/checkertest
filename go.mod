module backend-testing-module-checker

go 1.21.4

require (
	backend-testing-module-shared v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)

replace backend-testing-module-shared v1.0.0 => ../shared
