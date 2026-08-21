package match

import "strings"

func ParsePath(path string) (category, attr string, ok bool) {
        parts := strings.SplitN(path, ".", 2)
        if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
                return "", "", false
        }
        return parts[0], parts[1], true
}

func GetPath(bag Bag, path string) (any, bool) {
        cat, attr, ok := ParsePath(path)
        if !ok {
                return nil, false
        }
        return bag.Get(cat, attr)
}
