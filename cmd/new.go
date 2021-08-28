package new

import (
	"embed"
	"log"
	"os"
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

func CmdNew(dir string, tpls embed.FS) {
	templ := temp{
		Name: dir,
	}

	d, err := os.ReadDir(".")
	if err != nil {
		log.Panic(err)
	}

	for _, dd := range d {
		if dd.IsDir() && dd.Name() == dir {
			log.Panic("dir already exisits")
		}
	}
	os.Mkdir(dir, 0777)

	r, err := git.PlainInit(dir, false)
	if err != nil {
		log.Panic(err)
	}

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

	w, err := r.Worktree()
	if err != nil {
		log.Panic(err)
	}

	files := []string{"README.md", "main.go"}
	for _, v := range files {
		_, err = w.Add(v)
		if err != nil {
			log.Panic(err)
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
		log.Panic(err)
	}

	headRef, err := r.Head()
	if err != nil {
		log.Panic(err)
	}
	ref := plumbing.NewHashReference("refs/heads/main", headRef.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		log.Panic(err)
	}

	// err = w.Checkout(&git.CheckoutOptions{Branch: "main"})
	// if err != nil {
	// 	log.Panic(err)
	// }

	err = r.Storer.RemoveReference("refs/heads/master")
	if err != nil {
		log.Panic(err)
	}
}
