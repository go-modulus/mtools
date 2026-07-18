package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

func PrintLogo() {
	str := `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢠⣄⡄⠀⠀⠀⠀⢀⣤⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣤⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢸⣿⣿⡄⠀⠀⠀⣾⣿⣿⠀⠀⢀⣠⣠⣀⡀⠀⠀⢀⣠⣠⣀⢺⣿⠀⣀⣀⠀⠀⢀⣀⡀⠸⣿⡇⢀⣀⡀⠀⠀⣀⣀⠀⢀⣀⣄⣄⣀⠀⠀⠀⠀⣠⣿⣏⣈⠀⣀⡄⣀⣄⠀⢀⣠⣠⣀⡀⠀⢀⣀⣀⣠⣠⡀⢀⣠⣄⡀⠀⠀⢀⣠⣠⣀⡀⠀⣀⣀⠀⠀⣀⣀⠀⠀⣀⣀⠀⢀⣠⣠⣀⡀⠀⢀⣀⣀⣠⡀⢸⣿⠂⠀⣀⣠⡀
⠀⠀⠀⠀⢸⣿⠹⣿⡄⠀⣼⡟⢽⣿⠀⣴⣿⠋⠉⠻⣿⡆⢰⣿⠟⠉⠛⣿⣿⠀⣺⣿⠀⠀⢸⣿⡇⠸⣿⡇⢨⣿⡇⠀⠀⣿⣿⠀⣾⣟⡉⠙⠿⠆⠀⠀⠀⠙⣿⣏⠉⠠⣿⡿⠋⠋⠸⠟⠋⡉⣻⣿⡀⢼⣿⠏⠙⢻⣿⡛⠉⢻⣿⠄⣰⣿⢋⢉⢻⣿⡄⢹⣿⡀⢸⡿⣿⡀⢸⣿⠃⣴⣿⠛⠉⠻⣿⡆⢸⣿⡟⠙⠁⢺⣿⣡⣾⠟⠁⠀
⠀⠀⠀⠀⢸⣿⠀⢻⣷⣼⡿⠁⢽⣿⠀⢿⣿⠀⠀⢠⣿⡏⢸⣿⡅⠀⠀⣾⣿⠀⣺⣿⠀⠀⢸⣿⡇⠸⣿⡇⠰⣿⡇⠀⠀⣿⣷⠀⣈⡛⠛⠿⣷⣆⠀⠀⠀⢈⣿⣇⠀⠀⣿⡯⠀⠀⣴⣿⠛⠛⢻⣿⡂⢸⣿⠂⠀⢸⣿⠄⠀⢸⣿⠅⢿⣿⠛⠛⠛⣛⠃⠀⣿⣇⣿⠇⢻⣧⣾⡟⠀⢿⣿⠀⠀⢀⣿⡯⠸⣿⡇⠀⠀⢸⣿⠟⣿⣦⠀⠀
⠀⠀⠀⠀⠸⡿⠀⠀⠻⠿⠁⠀⠽⠿⠀⠈⠻⢿⢶⠿⠟⠁⠈⠻⢿⢾⠾⠻⠿⠀⠘⠿⡿⡾⠻⠿⠇⠸⠿⠇⠈⠻⢿⢾⠞⠿⠿⠀⠛⠿⡶⡾⠟⠃⠀⠀⠀⠠⠿⡇⠀⠀⡿⠯⠀⠀⠙⠿⡶⠾⠻⢿⠂⢸⢿⠁⠀⠸⡿⠂⠀⠸⡿⠅⠈⠻⢷⢶⠿⠟⠁⠀⠸⠿⠟⠀⠈⠿⠿⠀⠀⠈⠻⢿⢶⠿⠟⠁⠸⠿⠇⠀⠀⠸⡿⠂⠈⠻⡷⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`
	fmt.Println(color.CyanString(str))
}

func ExecCommand(ctx context.Context, name string, arg ...string) error {
	cmdCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	commandStr := name + " " + strings.Join(arg, " ")
	fmt.Printf("Running command %s:\n", color.BlueString(commandStr))

	out, err := exec.CommandContext(cmdCtx, name, arg...).CombinedOutput()

	coloredOut := color.HiBlackString(string(out))
	if err != nil {
		coloredOut = color.RedString(string(out))
	}
	if len(out) != 0 {
		fmt.Println(coloredOut)
	}
	return err
}
