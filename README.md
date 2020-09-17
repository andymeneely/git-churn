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

To run all the test cases:

```
  $ go test -v ./...
```

# Usage

In general, `git churn` works much like `git log`, with some additional options.

The `--repo` flag takes either github URL of the repo in which case it will clone the repo into the local memory and performs the operations, or you can specify the path to the cloned repo on your system. Use `"."` if the working directory is the repo to be used
Show basic churn metrics for a specific commit (compared to it's parent) and file:
```
  $ ./git-churn --help
  $ ./git-churn --repo https://github.com/ashishgalagali/SWEN610-project --commit c800ce62fc8a10d5fe69adb283f06296820522c1 --filepath src/main/java/com/webcheckers/ui/WebServer.java
  $ ./git-churn --repo https://github.com/andymeneely/git-churn --commit 00da33207bbb17a149d99301012006fbd86c80e4 --filepath testdata/file.txt --whitespace=false
  $ /path/to/git-churn --repo /path/to/repo --commit c800ce62fc8a10d5fe69adb283f06296820522c1 --filepath src/main/java/com/webcheckers/ui/WebServer.java
```

To get churn metrics for a range of commit:
```
  $ ./git-churn --repo https://github.com/ashishgalagali/SWEN610-project --commit c800ce62fc8a10d5fe69adb283f06296820522c1..5a2bf9f4da3de056dde3d9a9c18859de124d2602 --filepath src/main/java/com/webcheckers/ui/WebServer.java 
  $ ./git-churn --repo https://github.com/ashishgalagali/SWEN610-project --commit c800ce62fc8a10d5fe69adb283f06296820522c1...5a2bf9f4da3de056dde3d9a9c18859de124d2602 --filepath src/main/java/com/webcheckers/ui/WebServer.java --whitespace=false
  $ ./git-churn --repo https://github.com/ashishgalagali/SWEN610-project --commit c800ce62fc8a10d5fe69adb283f06296820522c1...5a2bf9f4da3de056dde3d9a9c18859de124d2602 --whitespace=false 
  $ /path/to/git-churn --repo . --commit c800ce62fc8a10d5fe69adb283f06296820522c1..5a2bf9f4da3de056dde3d9a9c18859de124d2602 --filepath src/main/java/com/webcheckers/ui/WebServer.java

```

To show the aggregated churn metrics for a specific commit:
```
 $ ./git-churn --repo https://github.com/andymeneely/git-churn --commit 00da33207bbb17a149d99301012006fbd86c80e4  --whitespace=false
 $ ./git-churn --repo https://github.com/ashishgalagali/SWEN610-project --commit c800ce62fc8a10d5fe69adb283f06296820522c1...5a2bf9f4da3de056dde3d9a9c18859de124d2602 --whitespace=false 

```

# Options
```
Flags:
  -a, --aggregate string   Aggregate the churn metrics. "commit": Aggregates all files in a commit. "all": Aggregate all files all commits and all files (default "commit")
  -c, --commit string      Commit hash for which the metrics has to be computed
  -f, --filepath string    File path for the file on which the commit metrics has to be computed
  -h, --help               help for git-churn
  -j, --json               Writes the JSON output to a file within a folder named churn-details
  -p, --print              Prints the output in a human readable format (default true)
  -r, --repo string        Git Repository URL on which the churn metrics has to be computed
  -w, --whitespace         Excludes whitespaces while calculating the churn metrics is set to false (default true)
```

# Sample Output

For a commit range
```
{
  "BaseCommitId": "8b0c2116cea2bbcc8d0075e762b887200a1898e1",
  "CommitDetails": [
    {
      "CommitId": "3895dfa31c54adf83fdaffd90cf1b5fd4e5d7ff0",
      "CommitAuthor": "mcuadros@gmail.com",
      "DateTime": "2019-11-01 10:06:13 +0100 +0100",
      "CommitMessage": "Merge pull request #1235 from jmahler/master\n\nfix broken link (s/ftp/https/)",
      "ChurnMetrics": [
        {
          "FilePath": "_examples/ls-remote/main.go",
          "DeletedLinesCount": 1,
          "SelfChurnCount": 0,
          "InteractiveChurnCount": 1,
          "ChurnDetails": {
            "b4fba7ede146be79cf65b89975250cf6869fb409": "v.cocaud@gmail.com"
          }
        },
        {
          "FilePath": "_examples/merge_base/helpers.go",
          "DeletedLinesCount": 2,
          "SelfChurnCount": 0,
          "InteractiveChurnCount": 2,
          "ChurnDetails": {
            "66c4a36212ced976c33712ca4fb6abc6697f2654": "David.Pordomingo.F@gmail.com"
          }
        }
      ]
    },
    {
      "CommitId": "3ed21ff5df781c947aebcf1d602269b1206116e3",
      "CommitAuthor": "jmmahler@gmail.com",
      "DateTime": "2019-10-31 18:05:28 -0700 -0700",
      "CommitMessage": "fix broken link (s/ftp/https/)\n\nSigned-off-by: Jeremiah Mahler <jmmahler@gmail.com>\n",
      "ChurnMetrics": [
        {
          "FilePath": "_examples/ls-remote/main.go",
          "DeletedLinesCount": 1,
          "SelfChurnCount": 0,
          "InteractiveChurnCount": 1,
          "ChurnDetails": {
            "b4fba7ede146be79cf65b89975250cf6869fb409": "v.cocaud@gmail.com"
          }
        },
        {
          "FilePath": "_examples/merge_base/helpers.go",
          "DeletedLinesCount": 2,
          "SelfChurnCount": 0,
          "InteractiveChurnCount": 2,
          "ChurnDetails": {
            "66c4a36212ced976c33712ca4fb6abc6697f2654": "David.Pordomingo.F@gmail.com"
          }
        }
      ]
    }
  ]
}
```

For all files in a commit aggregated 
```
{
  "BaseCommitId": "99992110e402f26ca9162f43c0e5a97b1278068a",
  "AggCommitDetails": [
    {
      "CommitId": "180ec07da5d7a415b48fd3d9f7d5c9dd2925780e",
      "CommitAuthor": "ashishgalagali@gmail.com",
      "DateTime": "2020-03-28 00:59:14 -0400 -0400",
      "CommitMessage": "Merge pull request #19 from andymeneely/diffMetrics\n\nGetting git diff metrics for a given commit and file",
      "AggChurnMetrics": {
        "FilesCount": 4,
        "TotalDeletedLinesCount": 25,
        "TotalSelfChurnCount": 22,
        "TotalInteractiveChurnCount": 3
      }
    },
    {
      "CommitId": "3854e533318df4f5bb9a059c76ddd8bb2464a620",
      "CommitAuthor": "ashishgalagali@gmail.com",
      "DateTime": "2020-03-28 00:57:17 -0400 -0400",
      "CommitMessage": "Diff Merics whitespace excluded\n",
      "AggChurnMetrics": {
        "FilesCount": 4,
        "TotalDeletedLinesCount": 25,
        "TotalSelfChurnCount": 22,
        "TotalInteractiveChurnCount": 3
      }
    }
  ]
}
```

# Metrics

* Lines added
* Lines deleted
* Churn (lines added + deleted)
* Number of authors
* Number of committers
* Inn


# UML
Generated using https://www.dumels.com/

![Alt text](git-churn_UML.svg?raw=true "UML")

# Profiling

To generate profiling files, execute the following commands inside the metrics' directory where the test cases are present:

```
go test -cpuprofile cpu.prof
go test -memprofile mem.prof
```

In order to visualize and see the profiling data on to a web page, execute the pollowing commands:

```
go tool pprof -http=localhost:12345 cpu.prof 
go tool pprof -http=localhost:12346 mem.prof 
```
