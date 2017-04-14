# goflycheck

The ```goflycheck``` program is a wrapper around the ```go build``` tool to provide on the fly check by reading the modified file from stdin. It copies all the files belong to the same package to a temp folder and run the ```go build``` in the temp folder and output the raw output of ```go build```

# Install
```
go get -u github.com/dzhou121/goflycheck
```
