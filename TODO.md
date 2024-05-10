## Todos

Here we keep track of the things that need to be done in the project so we have a minimal viable JVM implementation, later we should use GitHub issues in order to keep track of the project.

### Documentation

- [x] Add a `TODO.md` file to the project root.
- [x] Add a `README.md` file to the project root.
- [ ] Add a `CONTRIBUTING.md` file to the project root.
- [ ] Add a `LICENSE` file to the project root.

### Features

- [x] Add a simple CLI to run the interpreter (clone the `java` one).
- [x] Implement a Big Endian byte reader
- [x] Be able to read a `.class` file
- [ ] Be able to read a `.jar` file

- [ ] Implement constant pool validations (we just assume it is correct)
- [ ] Validate the class file object
- [ ] Implement all constant types

### Quality of Life

- [x] Add a `Makefile` to the project root.
- [ ] Group test fixtures by utility
- [ ] Add GitHub actions to test
- [ ] Add linting to our CI
- [ ] Think of a better way to implement the tests inside `tests/class_file`

### Bugs

### Tests

- [ ] Add a `record` type class to test
- [ ] Add an `abstract class` type class to test
- [ ] Write a simple LinkedList class and use it in a test
