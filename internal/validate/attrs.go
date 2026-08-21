package validate

import "fmt"

func AttrMaps(subject, resource, action, env map[string]any) error {
        if len(subject) == 0 && len(resource) == 0 && len(action) == 0 && len(env) == 0 {
                return fmt.Errorf("empty attribute bag")
        }
        return nil
}
