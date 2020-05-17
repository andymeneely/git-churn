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

To use it in the form of API, run the 
```
./git-churn
```
and you can find the swagger documentation on the link:
http://localhost:8080/swagger/index.html 
Click on "Try it out" to insert values into the query param and execute the APIs

Following the example APIs:
To get the churn metrics for a specific file:
http://localhost:8080/churn-metrics/file?repoUrl=https://github.com/andymeneely/git-churn&commitId=99992110e402f26ca9162f43c0e5a97b1278068a&whitespace=false&filepath=README.md

To get the aggregated churn metrics:
http://localhost:8080/churn-metrics/aggr?repoUrl=https://github.com/andymeneely/git-churn&commitId=99992110e402f26ca9162f43c0e5a97b1278068a&whitespace=true

For Developers:
If you make any Swagger API documentation changes, run the following command to update the swagger UI:
```
 $ swag init 
```
If it doesn't work try:
```
 $ GOPATH/bin/swag init 
```


In general, `git churn` works much like `git log`, with some additional options.

Show basic churn metrics for a specific commit and file:
```
  $ ./git-churn --help
  $ ./git-churn --repo https://github.com/ashishgalagali/SWEN610-project --commit c800ce62fc8a10d5fe69adb283f06296820522c1 --filepath src/main/java/com/webcheckers/ui/WebServer.java
  $ ./git-churn --repo https://github.com/andymeneely/git-churn --commit 00da33207bbb17a149d99301012006fbd86c80e4 --filepath testdata/file.txt --whitespace=false
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
  -w, --whitespace        Excludes whitespaces while calculating the churn metrics if set to false (default true)
```

# Sample Output

For a perticular file
```
{
  "FilePath": "src/main/java/com/webcheckers/ui/WebServer.java",
  "DeletedLinesCount": 13,
  "SelfChurnCount": 3,
  "InteractiveChurnCount": 10,
  "CommitAuthor": "ashishgalagali@gmail.com",
  "ChurnDetails": {
    "16123ab124432a058ed29e7d8fb2df52c310363b": "ashishgalagali@gmail.com",
    "9708c9a9da36928fd0b7143c74aa61694999fe5d": "ks3057@rit.edu",
    "979fe965043d49814c2fb7e7f5bae3461911b88b": "ashishgalagali@gmail.com",
    "b742aaf3e500712668d6f76c9736637436ee695e": "ks3057@rit.edu",
    "cef4dbea729fac483b43e130271c9e6efe93df33": "ks3057@rit.edu"
  },
  "FileDiffMetrics": {
    "Insertions": 17,
    "Deletions": 13,
    "LinesBefore": 154,
    "LinesAfter": 158,
    "File": "src/main/java/com/webcheckers/ui/WebServer.java",
    "NewFile": false,
    "DeleteFile": false
  }
}

```

For all files in a commit aggregated 
```
{
  "DeletedLinesCount": 110,
  "SelfChurnCount": 74,
  "InteractiveChurnCount": 36,
  "CommitAuthor": "ashishgalagali@gmail.com",
  "AggrDiffMetrics": {
    "Insertions": 225,
    "Deletions": 110,
    "LinesBefore": 3273,
    "LinesAfter": 3386,
    "FilesCount": 59,
    "NewFiles": 4,
    "DeletedFiles": 0
  }
}
```

# Metrics

* Lines added
* Lines deleted
* Churn (lines added + deleted)
* Number of authors
* Number of committers
* Inn
