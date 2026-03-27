package db

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
)

func printMigrationLog(
	logBuf string,
	namesMap map[string]string,
	err error,
) {
	logLines := strings.Split(logBuf, "\n")
	transformedLines := make([]string, len(logLines))
	applliedLinesMap := make(map[int]struct{})
	isRollingBack := false
	for i, line := range logLines {
		after, ok := strings.CutPrefix(line, "Applied: ")
		if !ok {
			after, ok = strings.CutPrefix(line, "Rolled back: ")
			if !ok {
				transformedLines[i] = line
				continue
			}
			isRollingBack = true
		}
		parts := strings.Split(after, " in ")
		after = parts[0]

		migrationFile := strings.TrimSpace(after)
		path, ok := namesMap[migrationFile]
		if !ok {
			transformedLines[i] = line
			continue
		}
		migrationFile = path + "/" + migrationFile

		if !strings.HasPrefix(migrationFile, "/") {
			path, _ := os.Getwd()
			migrationFile = path + "/" + migrationFile
		}

		migrationFile = strings.ReplaceAll(migrationFile, "/./", "/")
		migrationFile = strings.ReplaceAll(migrationFile, "//", "/")
		migrationFile = "file://" + migrationFile

		applliedLinesMap[i] = struct{}{}
		transformedLines[i] = migrationFile
	}

	maxMapKey := -1
	if len(applliedLinesMap) > 0 {
		maxMapKey = slices.Max(slices.Collect(maps.Keys(applliedLinesMap)))
	}
	for i, line := range transformedLines {
		if _, ok := applliedLinesMap[i]; ok {
			if err != nil && maxMapKey == i {
				fmt.Println(color.RedString("Error in Migration: %s", line))
				continue
			}
			if isRollingBack {
				fmt.Println(color.GreenString("Rolled back migration: %s", line))
				continue
			}
			fmt.Println(color.GreenString("Applied migration: %s", line))
		} else {
			fmt.Println(line)
		}
	}
}
