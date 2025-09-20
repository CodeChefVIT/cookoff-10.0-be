package submissions

import "fmt"

func RuntimeMut(language_id int) (int, error) {
	var runtime_mut int
	switch language_id {
	case 50, 54, 60, 73, 63:
		runtime_mut = 1
	case 51, 62:
		runtime_mut = 2
	case 68:
		runtime_mut = 3
	case 71:
		runtime_mut = 5
	default:
		return 0, fmt.Errorf("invalid language ID: %d", language_id)
	}
	return runtime_mut, nil
}
