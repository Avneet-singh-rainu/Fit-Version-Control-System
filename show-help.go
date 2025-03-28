package main

import (
	"fmt"

	"github.com/fatih/color"
)

func Help() {

	bold := color.New(color.Bold)
	info := color.New(color.FgCyan, color.Bold).SprintFunc()
	success := color.New(color.FgGreen, color.Bold).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()
	cmdColor := color.New(color.FgMagenta, color.Bold).SprintFunc()


	fmt.Println(info("🚀 Fit CLI - Version 1.0"))
	fmt.Println(info("-------------------------"))


	fmt.Println(info("\n📋 Usage:"))
	fmt.Println(success("  fit help"), "      			- " + bold.Sprint("Show available commands"))
	fmt.Println(success("  fit init"), "      			- " + bold.Sprint("Initialize Fit repository"))
	fmt.Println(success("  fit add ."), "      			- " + bold.Sprint("Add all files for tracking"))
	fmt.Println(success("  fit add [filename]"), "			- " + bold.Sprint("Add specific file"))
	fmt.Println(success("  fit commit -m 'message'"), "		- " + bold.Sprint("Commit changes with a message"))
	fmt.Println(success("  fit cto [commit_hash]"), "		- " + bold.Sprint("Checkout to a specific commit"))
	fmt.Println(success("  fit log"), "       			- " + bold.Sprint("Show commit history"))


	fmt.Println(warn("\n🔍 Quick Start Example:"))
	fmt.Println(cmdColor("  $ fit init"))
	fmt.Println(cmdColor("  $ fit add ."))
	fmt.Println(cmdColor("  $ fit commit -m \"Initial commit\""))
	fmt.Println(cmdColor("  $ fit log"))


	fmt.Println(info("\n💡 Pro Tip: ") + "Use " + success("fit help [command]") + " for detailed command information \n\n")
}
