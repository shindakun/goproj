package new

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type temp struct {
	Name string
}

// dirCheck checks for exisiting directory and return an error if so
func dirCheck(dir string) error {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Panic(err)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() == dir {
			return fmt.Errorf("dir already exisits")
		}
	}
	return nil
}

// doGit does all of the Git related tasks, init, add, commit
func doGit(dir string) error {
	r, err := git.PlainInit(dir, false)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	filesToAdd := []string{"README.md", "main.go", "go.mod"}
	for _, fileToAdd := range filesToAdd {
		_, err = w.Add(fileToAdd)
		if err != nil {
			return err
		}
	}
	_, err = w.Commit("inital commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Steve Layton",
			Email: "shindakun@users.noreply.github.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	headRef, err := r.Head()
	if err != nil {
		return err
	}
	ref := plumbing.NewHashReference("refs/heads/main", headRef.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	// err = w.Checkout(&git.CheckoutOptions{Branch: "main", Create: true, Force: true})
	// if err != nil {
	// 	log.Panic(err)
	// }

	err = r.Storer.RemoveReference("refs/heads/master")
	if err != nil {
		return err
	}

	return nil
}

func CmdNew(dir string, tpls embed.FS) {
	templ := temp{
		Name: dir,
	}

	err := dirCheck(dir)
	if err != nil {
		panic(err)
	}

	os.Mkdir(dir, 0777)

	t, err := template.ParseFS(tpls, "templates/*")
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create(filepath.Join(dir, "main.go"))
	if err != nil {
		log.Panic(err)
	}

	t.ExecuteTemplate(f, "main.tmpl", templ)

	f, err = os.Create(filepath.Join(dir, "README.md"))
	if err != nil {
		log.Panic(err)
	}

	t.ExecuteTemplate(f, "README.tmpl", templ)

	err = os.Chdir(dir)
	if err != nil {
		log.Panic(err)
	}

	cmd := exec.Command("go", "mod", "init")
	err = cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	err = os.Chdir("..")
	if err != nil {
		log.Panic(err)
	}

	err = doGit(dir)
	if err != nil {
		log.Panic(err)
	}
}
