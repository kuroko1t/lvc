package main

import (
	//"time"
	"os"
	"os/exec"
	"strings"
	"fmt"
	"log"
	"path/filepath"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/plumbing"
	"flag"
)

func getHashList(url, fromhash string) []plumbing.Hash {
	//r, err := git.PlainClone("./", false, &git.CloneOptions{
	// 	URL: "https://github.com/kuroko1t/nne",
	//})
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/" + url,
	})
	if err != nil {
		log.Fatal(err)
	}
	ref, err := r.Head()
	if err != nil {
		log.Fatal(err)
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	allHashList := make([]plumbing.Hash, 0)
	hash := plumbing.NewHash(fromhash)
	err = cIter.ForEach(func(c *object.Commit) error {
		allHashList = append(allHashList, c.Hash)
		return nil
	})
	fromHashList := make([]plumbing.Hash, 0)
	for _, h := range allHashList {
		fromHashList = append(fromHashList, h)
		if h == hash {
			break
		}
	}
	return fromHashList
}

func checkoutProject(url string, hash plumbing.Hash) {
	path := strings.Split(url, "/")[1]
	var r *git.Repository
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		r, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: "https://github.com/kuroko1t/nne",
		})
		if err != nil {
			log.Fatal(err)
		}
	 } else {
		r, err = git.PlainOpen(path)
		if err != nil {
			log.Fatal(err)
		 }
	 }
	w, err := r.Worktree()
	err = w.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func build(url string) {
	path := strings.Split(url, "/")[1]
	prevDir, _ := filepath.Abs(".")
	os.Chdir(path)
	cmd := exec.Command("python", "setup.py", "install")
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	os.Chdir(prevDir)
}

func cmdrun(cmdstring string) {
	cmdlist := strings.Split(cmdstring, " ")
	//first:= cmdlist[0]
	var argscmd []string
	for i, c := range cmdlist {
		if i != 0 {
			argscmd = append(argscmd, c)
		}
	}
	cmd := exec.Command(cmdlist[0], argscmd...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		log.Fatal("Please set, giturl, hash, cmd")
	}
	url := args[0] //"kuroko1t/nne"
	hashlist := getHashList(url, args[1])
	cmd := args[2]
	for _, h := range hashlist {
		fmt.Println("HASH:[",h, "]")
		checkoutProject(url, h)
		build(url)
		cmdrun(cmd)
	}
}
