package main

import (
	"fmt"
	"path/filepath"
	"sort"
	//"go/printer"
	"io"
	"os"
	//"strings"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "")
}

// func printDir(out *io.Writer, path string, printFiles bool, prefix string) error { что такое *
func printDir(out io.Writer, path string, printFiles bool, prefix string) error {
	// Читаем содержимое директории
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Сортируем записи по имени
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	// Убираем файлы, если флаг printFiles == false
	if !printFiles {
		var dirs []os.DirEntry
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry)
			}
		}
		entries = dirs
	}

	// Обрабатываем каждый элемент директории
	for i, entry := range entries {
		// Определяем символы для соединения (последний элемент или нет)
		var connector string
		if i == len(entries)-1 {
			connector = "└───"
		} else {
			connector = "├───"
		}

		printBytes := !entry.IsDir() && printFiles
		if !printBytes {
			fmt.Fprintf(out, "%s%s%s\n", prefix, connector, entry.Name())
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				return err
			}
			size := fileInfo.Size()
			if size > 0 {
				fmt.Fprintf(out, "%s%s%s (%db)\n", prefix, connector, entry.Name(), size)

			} else {
				fmt.Fprintf(out, "%s%s%s (%s)\n", prefix, connector, entry.Name(), "empty")
			}
			// Выводим название директории или файла
			//fmt.Fprintf(out, "%s%s%s (%db)\n", prefix, connector, entry.Name(), fileInfo.Size()) что такое db
		}

		// Если это директория, вызываем функцию рекурсивно
		if entry.IsDir() {
			var newPrefix string
			if i == len(entries)-1 {
				newPrefix = prefix + "\t" // если последний элемент, добавляем отступ
			} else {
				newPrefix = prefix + "│\t" // для остальных элементов
			}
			printDir(out, filepath.Join(path, entry.Name()), printFiles, newPrefix)
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
