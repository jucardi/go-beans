package testx

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func getSize() (int, int) {
	cmd := exec.Command("stty", "size")
	data, err := cmd.Output()
	if err != nil {
		return -1, -1
	}
	result := strings.Split(string(data), " ")
	row, _ := strconv.Atoi(result[0])
	col, _ := strconv.Atoi(result[1])
	return row, col
}

func getPosStr(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}
