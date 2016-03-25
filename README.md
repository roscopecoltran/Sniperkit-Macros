### Wetting Your Appetite
Have you ever experienced headache to install libraries and dependencies?  
Ever had to deal with two incompatible versions of a program at once?  
Ever wished to try out a new language first, and install it only if it pleases you?  
Ever wished to develop for Linux when you use Mac OS or Windows?  
Ever wished to develop in Go from the folder of your choice?  
Ever wished to have a unified development tool, across all platforms, customizable to any languages?  

### Nut
**Nut** is a command line tool which offers a solution to common frustrations of developers. It hides the complexity of development environments, and extends them with customizable macros. Whether you develop in Swift, Go, Java, or C++, what you need is build/run/test the app. So just do it:

    $ nut build
    $ nut run
    $ nut test

**Nut** mounts the current folder in a [Docker](https://www.docker.com/) container, and executes commands on your behalf, according to the project configuration. The configuration is read from `nut.yml` file, in the current folder. You can choose the Docker image to use, declare volumes to mount, and define commands (called macros) such as *build*, *run*, and *test*.

### Nut File Syntax
#### Example
Here is an example of `nut.yml` to develop in Go. You can generate a sample configuration with  :

`nut --init`
```yaml
# nut.yml
project_name: nut
based_on:  # configuration can be inherited from:
  url: https://raw.githubusercontent.com/matthieudelaro/donut/master/go/nut.yml  # a URL (soon)
       # parameters set in this file will override basic configuration
  github: matthieudelaro/myenvironment  # a GitHub repository (soon)
  path: /home/matthieudelaro/environments/go/nut.yml  # a local file (soon)
  docker_image: golang:1.6  # directly from a Docker image
mount:  # declares folder to mount in the container
  main:  # give each folder any name that you like
  - .               # this folder (from your computer) will be mounted as
  - /go/src/project # this folder (in the container)
macros:  # macros define operations the Nut can perform
  build:  # call this one with "nut build"
    usage: build the project
    actions:  # a list of commands to run in the container
    - go build -o nut
  run:  # call this one with "nut build"
    usage: run the project in the container
    actions:
    - ./nut
container_working_directory: /go/src/project  # container folder in which macros will be executed
syntax_version: "2"  # Nut evolves quickly, its configuration file syntax as well,
                     # so nut files are versioned to ensure backward compatibility.
```
#### Guidelines
Nut aims to unify development tools, not to replace compilers.
Nut aims to unify development processes, not to expose language specific requirements.

So, when creating a `nut.yml` file, one should standard names for macros, such as:
- build
- run
- test
- debug
- deploy

As opposed to:
- javac (should be generalized with *run*)
- make (duplicate of *build*)
- do (hum... *Do* what?)
This will keep Nut easy to integrate in text editors and IDEs.

### What the Nut???
Achieved with Nut:  
- use Caffe with `nut checkModelDefinition`, `nut train`, `nut test`.
- compile CUDA code on a Mac Book Air, which has not any Nvidia GPU. Just `nut build`
- test code in a whole infrastructure, by defining a macro running *docker-compose* in a container.

### Milestones
- add support for port bindings
- `nut --init` project configuration from a GitHub repository
- plugin for Sublime Text, to call `nut run`, `nut build`, and `nut test` from the editor.

### Stay Tune
Wanna receive updates? Or share your thoughts? Please tell us on the [form](https://docs.google.com/forms/d/1reDwa7t2-8o_vPGuYg6QBCYHoDdge80dDbkBS9H72nM/viewform).

### Authors and Contributors
@matthieudelaro and @gdevillele

