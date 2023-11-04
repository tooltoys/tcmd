package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Tree() {
	app := tview.NewApplication()

	rootDir := "./cmd"
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetBackgroundColor(0)

	add := func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			fileName := file.Name()
			if file.IsDir() {
				fileName = fmt.Sprintf("%s %s", "[]", file.Name())
			} else {
				fileName = fmt.Sprintf("%s %s ", "</>", file.Name())
			}

			node := tview.NewTreeNode(fileName).
				SetReference(filepath.Join(path, file.Name()))
			node.SetColor(tcell.ColorRed)

			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		f, _ := os.Stat(reference.(string))
		if f.IsDir() {
			children := node.GetChildren()
			if len(children) == 0 {
				// Load and show files in this directory.
				path := reference.(string)
				add(node, path)
			} else {
				// Collapse if visible, expand if collapsed.
				node.SetExpanded(!node.IsExpanded())
			}
			return
		} else {
			f, err := os.Open(reference.(string))
			if err != nil {
				app.Stop()
				os.Stdout.Write([]byte(err.Error()))
			}
			defer f.Close()

			data, _ := io.ReadAll(f)
			var script Script
			err = json.Unmarshal(data, &script)
			if err != nil {
				app.Stop()
				os.Stdout.Write([]byte(err.Error()))
			}

			cmd := exec.Command(script.Cmd, script.Inputs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stdout
			app.Stop()
			if err := cmd.Run(); err != nil {
				os.Stdout.Write([]byte(err.Error()))
			}
		}
	})

	if err := app.SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

type Script struct {
	Cmd         string   `json:"cmd"`
	Description string   `json:"description"`
	Inputs      []string `json:"inputs"`
}
