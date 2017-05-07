/* main.go - Example for go-emacs API

Copyright (C) 2017 Yann Hodique <yann.hodique@gmail.com>.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or (at
your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.  */

package main

// int plugin_is_GPL_compatible;
import "C"

import (
	"github.com/sigma/go-emacs"
	git "gopkg.in/libgit2/git2go.v25"
)

func init() {
	emacs.Register(initModule)
}

func initModule(env emacs.Environment) {
	stdlib := env.StdLib()

	gitOpenFunc := env.MakeFunction(GitOpen, 1, "git-open", nil)
	gitOpenSym := stdlib.Intern("git-open")
	stdlib.Fset(gitOpenSym, gitOpenFunc)

	gitLsBranchesFunc := env.MakeFunction(GitLsBranches, 1, "git-ls-branches", nil)
	gitLsBranchesSym := stdlib.Intern("git-ls-branches")
	stdlib.Fset(gitLsBranchesSym, gitLsBranchesFunc)

	gitSym := stdlib.Intern("git")
	stdlib.Provide(gitSym)
}

func GitOpen(ctx emacs.FunctionCallContext) (emacs.Value, error) {
	path, err := ctx.GoStringArg(0)
	if err != nil {
		return nil, err
	}

	repo, err := git.OpenRepositoryExtended(path, 0, "")
	if err != nil {
		return nil, err
	}

	return ctx.Environment().MakeUserPointer(repo), nil
}

func GitLsBranches(ctx emacs.FunctionCallContext) (emacs.Value, error) {
	env := ctx.Environment()
	rawRepo, ok := env.ResolveUserPointer(ctx.UserPointerArg(0))
	if !ok {
		return emacs.Error("user-ptr does not exist")
	}

	repo, ok := rawRepo.(*git.Repository)
	if !ok {
		return emacs.Error("user-ptr is not a git repo")
	}

	iter, err := repo.NewBranchIterator(git.BranchLocal)
	if err != nil {
		return nil, err
	}

	branches := make([]emacs.Value, 0)

	var branchesRecorder git.BranchIteratorFunc = func(br *git.Branch, _ git.BranchType) error {
		name, err := br.Name()
		if err != nil {
			return err
		}
		branches = append(branches, env.String(name))
		return nil
	}
	iter.ForEach(branchesRecorder)

	return env.StdLib().List(branches...), nil
}

func main() {}
