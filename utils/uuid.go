package utils

import (
	"fmt"
	"strings"
	"github.com/pborman/uuid"
)

// Uuid generate uuid
func Uuid() (string) {
	return strings.Join(strings.Split(uuid.New(), "-"), "")
}


func main() {
	fmt.Println(Uuid())
}
