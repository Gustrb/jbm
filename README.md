## JBM

The java bullshit machine.
It is a simple JVM bytecode interpreter written in Golang.

It is a work in progress and is not yet functional, to see the progress check the [TODO](TODO.md) file.

### Tests

We have a lot of java classes written to test the interpreter, you can find them in the `test` directory.
To run the tests, you can use the `scripts/run_tests.sh` script, it will build all the java classes and run the interpreter with them.
So we will require the `javac` and `java` commands to be available in the system.

### Building

We provide a `Makefile` to build the project, you can use the `make` command to build the project.
The output will be in the `bin` directory.
