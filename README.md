# File Server

## About the project

This project show files which containing char "a" on earlier position than other files. Files located in "temp" folder.

Currently program builds, pushes and releases a Docker container to Heroku on commit to main branch.

Working example: [Heroku Deployed Example](https://gofileserver.herokuapp.com/)

### Design

This project runs on [FileServer](https://pkg.go.dev/net/http#example-FileServer). FileServer read files chunk by chunk with concurreny methods to find character "a". 

### Status

The template project is in alpha status.

### Getting started

To run the project in your local use the following commands

```
docker build -t file_server .
```
```
docker run -p 8080:8080 file_server
```

### Layout

```tree
├── main.go
├── go.mod
├── Dockerfile
├── heroku.yml
├── .gitignore
├── README.md
├── temp
│   ├── 1.txt
│   ├── 2.txt
│   └── 3.txt
```

A brief description of the layout:

* `.gitignore` Varies per project, but all projects need to ignore `bin` directory.
* `README.md` Is a detailed description of the project.
* `heroku.yml` Contains yml commands for deployment to heroku
* `Dockerfile` Contains set of commands to dockerize the project
* `main.go` Main application codes.
* `temp/` Contains all test files.

## Notes

* Currently the program does not download files, only shows file/files with char "a" in an earlier position than other files.