## JBM

The java bullshit machine.
It is a simple JVM bytecode interpreter written in Golang.

It is a work in progress and is not yet functional, to see the progress check the [TODO](TODO.md) file.

### Tests

We have a lot of java classes written to test the interpreter, you can find them in the `tests/fixtures` directory.

To run the tests you will need the [make](https://www.gnu.org/software/make/) command available in your system, along with the `javac` and `java` commands.

To run all tests you can simply:

```bash
$ make test
```

### Building

We provide a `Makefile` to build the project, you can use the `make` command to build the project.
The output will be in the `bin` directory.
