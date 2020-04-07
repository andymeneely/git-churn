# git-churn

A fast tool for collecting code churn metrics from git repositories.

# Installation
You will need Go language installed on your system. Ref: https://golang.org/doc/install

```
  $ git clone github.com/andymeneely/git-churn
  $ cd git-churn
  $ go install github.com/andymeneely/git-churn
  $ go build
 ```

# Usage

In general, `git churn` works much like `git log`, with some additional options.

Show basic churn metrics for a specific commit and file:
```
  $ git-churn --help
  $ git-churn --repo https://github.com/andymeneely/git-churn --commit 00da33207bbb17a149d99301012006fbd86c80e4 --filepath testdata/file.txt --whitespace=false
  $ git-churn --repo https://github.com/andymeneely/git-churn --commit 00da33207bbb17a149d99301012006fbd86c80e4  --whitespace=false
```

To show the aggregated churn metrics for a specific commit:
```
 $ git-churn --repo https://github.com/andymeneely/git-churn --commit 00da33207bbb17a149d99301012006fbd86c80e4  --whitespace=false
```

# Options
```
Flags:
  -c, --commit string     Commit hash for which the metrics has to be computed
  -f, --filepath string   File path for the file on which the commit metrics has to be computed
  -h, --help              help for git-churn
  -r, --repo string       Git Repository URL on which the churn metrics has to be computed
  -w, --whitespace        Excludes whitespaces while calculating the churn metrics is set to false (default true)
```

# Metrics

* Lines added
* Lines deleted
* Churn (lines added + deleted)
* Number of authors
* Number of committers
* Inn
