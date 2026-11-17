# How to Contribute

We'd love to accept your patches and contributions to this project.

## Before you begin

### Sign off your commits (Developer Certificate of Origin)

Contributions to this project must be accompanied by a sign-off under the
[Developer Certificate of Origin][dco] (DCO) version 1.1. You retain the
copyright to your contribution; the sign-off certifies that you have the right
to submit it under the project's license ([Apache License 2.0](LICENSE)).

Add the sign-off by committing with the `-s` flag:

```shell
git commit -s -m "Describe your change"
```

This appends a line like the following to the commit message, using the name
and email from your Git configuration:

```
Signed-off-by: Your Name <your.email@example.com>
```

The full text of the Developer Certificate of Origin, version 1.1:

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

### Review our Code of Conduct

This project follows the [Contributor Covenant][conduct], version 2.1. By
participating, you are expected to uphold it.

## Contribution process

We follow the [fork and branch workflow][workflow].

There will be three Git repositories involved:

1.  *upstream* - the Weft repository on GitHub, `github.com/brunobls/weft`.
2.  *origin* - your GitHub fork of `upstream`. This repository
    will typically be at a URL that looks like `github.com/_your_user_name_/weft`
3.  *local* - your local clone of `origin`

### First time setup

Follow these steps to get ready for making changes to Weft.  These
steps are only needed once and not for subsequent changes you might want to
make:

1.  Fork the Weft repository on GitHub to create `origin`.
    Visit the [Weft][github_weft] GitHub repository and click the `Fork` button.

2.  Make a `local` clone of your fork.

    ```shell
    git clone git@github.com:_your_user_name_/weft.git
    ```

3.  Add a remote pointing from `local` to `upstream`.

    ```shell
    cd weft
    git remote add upstream git@github.com:brunobls/weft.git
    ```

### Making changes

Here is a detailed outline of the steps needed to make changes to Weft.


1. Make a local branch in your clone and pull any recent changes into it.

   ```shell
   git switch -c my_branch  # Pick a name appropriate to your work
   git pull upstream main
   ```

2. Make changes and commit to local branch.

   ```shell
   # ... editing, testing, ... 
   git commit -s ...
   ```

3. Pull any changes that may have been made in the upstream repository
   main branch.

   ```shell
   git switch my_branch
   git pull --rebase upstream main
   ```

   Note that this command may result in merge conflicts. Fix those if
   needed.

4. Push your branch to the corresponding branch in your fork (the `origin` repository).

   ```shell
   git switch my_branch
   git push origin my_branch
   ```

5. Select the branch you are working on in the drop-down menu of branches on
   https://github.com/_your_user_name_/weft . Then hit the `Compare and pull
   request` button.

6. Respond to feedback, which may involve making new commits.
   If you made any changes, push them to github again.

   ```shell
   git switch my_branch
   git push origin my_branch
   ```

   Repeat as necessary until all feedback has been handled.

   Note: the preceding approach will cause the pull request to become a sequence
   of commits. Some people like to keep just a single commit that is amended as
   changes are made. If you are amending commits that had already been pushed,
   you will have to add `--force` to the `git push` command above.

7. Once reviewers are happy, pull any main branch changes that may
   have happened since step 3.
   
    ```shell
    git switch my_branch
    git pull --rebase upstream main
    ```

    If some changes were pulled, push again to the PR, but this time you will
    need to force push since the rebase above will have rewritten your commits.

    ```shell
    git switch my_branch
    git push --force origin my_branch
    ```

8.  Ask somebody who has permissions (or do it yourself if you
    have permissions) to merge your branch into the main branch
    of the `upstream` repository. The reviewer may do this without
    being asked.

    Select the `Squash and merge` option on https://github.com/brunobls/weft
    or use the command line instructions found on that page. Edit the commit message
    as appropriate for the squashed commit, and keep the `Signed-off-by` trailer(s)
    of the commits being squashed in the final message (see
    [Sign off your commits](#sign-off-your-commits-developer-certificate-of-origin)).

9.  Delete the branch from `origin`:

    ```
    git push origin --delete my_branch
    ```

10. Delete the branch from `local`

    ```
    git switch main
    git branch -D my_branch
    ```


### Code Reviews

All submissions, including submissions by project members, require review. We
use GitHub pull requests for this purpose. Consult [GitHub Help][github_help]
for more information on using pull requests.

To make sure changes are well coordinated, we ask you to discuss any significant
change prior to sending a pull request. To do so, either file a
[new issue][new_issue] or claim an [existing one][issues]. Issues and pull
requests on GitHub are the project's channel for coordination.

### License headers

Source files in this repository carry an Apache 2.0 license header:

*   A **new** file that does not come from Service Weaver gets the header
    `Copyright <year> The Weft Authors`.
*   A file that **already has** a `Copyright Google LLC` header keeps that line
    unchanged when you modify it (the Apache License requires keeping existing
    copyright notices).
*   A file copied from Service Weaver keeps its original header.

## What to contribute to?

Here are the current areas where the community contributions would be most
valuable:

* Implement deployers for platforms other than local, SSH, and GKE.
* Bug fixes for the existing libraries.
* Non-Go language support. (Note: please open an [issue][new_issue] first, to
  coordinate development. We will likely require a fair bit of API and
  language-integration design discussion before the implementation starts.)

## Attribution

Weft is a derivative work of [Service Weaver][upstream_weaver], originally
created by Google LLC and licensed under the Apache License 2.0. See the
[NOTICE](NOTICE) file.

[conduct]: CODE_OF_CONDUCT.md
[dco]: https://developercertificate.org/
[github_help]: https://help.github.com/articles/about-pull-requests/
[github_weft]: https://github.com/brunobls/weft
[issues]: https://github.com/brunobls/weft/issues
[new_issue]: https://github.com/brunobls/weft/issues/new
[upstream_weaver]: https://github.com/ServiceWeaver/weaver
[workflow]: https://docs.github.com/en/get-started/exploring-projects-on-github/contributing-to-a-project
