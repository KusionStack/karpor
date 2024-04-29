---
title: Test Conventions
---

## Testing Principles

In Karpor, we primarily focus on the following three types of tests:

- Unit tests: Tests targeting the **smallest testable units** (such as functions, methods, utility classes, etc.)
- Integration tests: Tests targeting the interaction and integration between **multiple units (or modules)**
- End-to-End tests (e2e tests): Tests targeting the **entire system's behavior**, usually requiring the simulation of real user scenarios

Each has its strengths, weaknesses, and suitable scenarios. To achieve a better development experience, we should adhere to the following principles when writing tests.

**Testing principles**：

- A case should only cover one scenario
- Follow the **7-2-1 principle**, i.e., 70% unit tests, 20% integration tests, and 10% end-to-end tests
- **Avoid using Mock frameworks in unit tests unless necessary** (e.g., `golang/mock`). If you feel the need to use a Mock framework in unit tests, what you actually need might be integration tests or even end-to-end tests

## Technology Selection

At the current point in time, the most popular testing frameworks in the Go language ecosystem are [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/) and [Testify](https://github.com/stretchr/testify). Therefore, this section mainly discusses the characteristics, pros and cons, and the final selection of these two testing frameworks.

### Ginkgo/Gomega

**Advantages**:

1. **BDD Support**: Ginkgo is favored by many developers for its support of Behavior-Driven Development (BDD) style. It offers a rich DSL syntax, making test cases more descriptive and readable through keywords like `Describe`, `Context`, `It`, etc.
2. **Parallel Execution**: Ginkgo can execute tests in parallel across different processes, improving the efficiency of test execution.
3. **Rich Matchers**: Used in conjunction with the Gomega matchers library, it provides a wealth of assertion capabilities, making tests more intuitive and convenient.
4. **Asynchronous Support**: Ginkgo provides native support for handling complex asynchronous scenarios, reducing the risk of deadlocks and timeouts.
5. **Test Case Organization**: Supports organizing test cases into Suites for easy management and expansion.
6. **Documentation**: Ginkgo's [official documentation](http://onsi.github.io/ginkgo/) is very detailed, offering guides from beginner to advanced usage.

**Disadvantages**:

1. **Learning Curve:** For developers not familiar with BDD, Ginkgo's DSL syntax may take some time to get used to.
2. **Complexity in Parallel Testing:** Although Ginkgo supports parallel execution, managing resources and environment for parallel tests can introduce additional complexity in some cases.

### Testify

**Advantages**:

1. **Simplified API**: Testify provides a simple and intuitive API, easy to get started with, especially for developers accustomed to the `testing` package.
2. **Mock Support**: Testify offers powerful Mock functionalities, facilitating the simulation of dependencies and interfaces for unit testing.
3. **Table-Driven Tests**: Supports table-driven testing, allowing for easy provision of various inputs and expected outputs for the same test function, enhancing test case reusability.
4. **Compatibility with `testing` Package**: Testify is highly compatible with the Go standard library's `testing` package, allowing for seamless integration into existing testing workflows.
5. Documentation: Testify's [official documentation](https://pkg.go.dev/github.com/stretchr/testify) also provides rich introductions on how to use its assertion and mocking functionalities.

**Disadvantages**:

1. **Lack of BDD Support**: Testify does not support the BDD style, potentially less intuitive for developers looking to improve test case readability.
2. **Relatively Simple Features**: Compared to Ginkgo, Testify's features are relatively simple and may not meet some complex testing scenarios' requirements.

### Summary

In short, Ginkgo/Gomega offers better readability and maintainability, producing clean and clear tests, but with a higher learning curve requiring familiarity with BDD style. Testify is simpler, more practical, with a lower learning curve, but as time progresses, the testing code style may become more varied, lowering maintainability.

Considering the actual situation of Karpor and the pros and cons of both frameworks, we decide to use these two frameworks in combination:
- Use Testify for unit testing, adhering to [Table-Driven Testing](https://go.dev/wiki/TableDrivenTests) to constrain the code style and prevent decay;
- Utilize Ginkgo's BDD features for writing higher-level integration and end-to-end tests;

This combination fully leverages the strengths of both frameworks, improving the overall efficiency, readability, and quality of testing.

## Writing Specifications

### Test Style

[Table-Driven Testing](https://go.dev/wiki/TableDrivenTests) is a best practice for writing test cases, akin to design patterns in programming, and it is also the style recommended by the official Go language. Table-Driven Testing uses tables to provide a variety of inputs and expected outputs, allowing the same test function to verify different scenarios. The advantages of this method are that it increases the reusability of test cases, reduces repetitive code, and makes tests clearer and easier to maintain.

While there is no direct syntax support for Table-Driven Testing in Go's `testing` package, it can be emulated by writing helper functions and using anonymous functions.

Here is an example of Table-Driven Testing implemented in the Go standard library's `fmt` package:

```go
var flagtests = []struct {
    in  string
    out string
}{
    {"%a", "[%a]"},
    {"%-a", "[%-a]"},
    {"%+a", "[%+a]"},
    {"%#a", "[%#a]"},
    {"% a", "[% a]"},
    {"%0a", "[%0a]"},
    {"%1.2a", "[%1.2a]"},
    {"%-1.2a", "[%-1.2a]"},
    {"%+1.2a", "[%+1.2a]"},
    {"%-+1.2a", "[%+-1.2a]"},
    {"%-+1.2abc", "[%+-1.2a]bc"},
    {"%-1.2abc", "[%-1.2a]bc"},
}
func TestFlagParser(t *testing.T) {
    var flagprinter flagPrinter
    for _, tt := range flagtests {
        t.Run(tt.in, func(t *testing.T) {
            s := Sprintf(tt.in, &flagprinter)
            if s != tt.out {
                t.Errorf("got %q, want %q", s, tt.out)
            }
        })
    }
}
```

It is worth noting that most mainstream IDEs have already integrated [gotests](https://github.com/cweill/gotests), enabling the automatic generation of table-driven style Go unit tests, which I believe can enhance the efficiency of writing your unit tests:

- [GoLand](https://blog.jetbrains.com/go/2020/03/13/test-driven-development-with-goland/)
- [Visual Studio Code](https://betterprogramming.pub/a-quick-way-to-generate-go-tests-in-visual-studio-code-b7c675b88dac)

### File Naming

- **Specification Content**：Test files should end with `_test.go` to distinguish between test code and production code.
- **Positive Example**：`xxx_test.go`
- **Negative Example**：`testFile.go`、`test_xxx.go`

### Test Function Naming

- **Specification**: The name of the test function should start with `Test`, followed by the name of the function being tested, using camel case notation.
- **Positive Example**：
  ```go
  func TestAdd(t *testing.T) {
      // Test logic ...
  }
  ```
- **Negative Example**：
  ```go
  func TestAddWrong(t *testing.T) {
      // Test logic ...
  }
  ```

### Test Function Signature

- **Specification Content**: The signature of the test function should be `func TestXxx(t *testing.T)`, where `t` is the test object, of type `*testing.T`, and there should be no other parameters or return values.
- **Positive Example**：
  ```go
  func TestSubtraction(t *testing.T) {
      // Test logic ...
  }
  ```
- **Negative Example**：
  ```go
  func TestSubtraction(value int) {
      // Test logic ...
  }
  ```

### Test Organization

- **Specification Content**：Test cases should be independent of each other to avoid mutual influence between tests; use sub-tests (`t.Run`) to organize complex test scenarios.
- **Positive Example**：
  ```go
  func TestMathOperations(t *testing.T) {
      t.Run("Addition", func(t *testing.T) {
          // Test addition logic ...
      })
      t.Run("Subtraction", func(t *testing.T) {
          // Test subtraction logic ...
      })
  }
  ```
- **Negative Example**：
  ```go
  func TestMathOperations(t *testing.T) {
      // Mixed addition and subtraction logic...
  }
  ```

### Test Coverage

- **Specification Content**：Attention should be paid to test coverage, use the `go test -cover` command to examine the test coverage of the code.

- **Positive Example**：
  ```shell
  $ go test -cover
  ```
- **Negative Example**：
  ```shell
  $ go test # Without checking test coverage
  ```
- **Note**: Karpor has wrapped this command as `make cover`, which will output the coverage for each package and total coverage in the command line. If you would like to view the coverage report in the browser, please execute `make cover-html`.

### Benchmark Tests

- **Specification Content**：Benchmark test functions should start with `Benchmark` and accept an argument of type `*testing.B`, focusing on performance testing.
- **Positive Example**：
  ```go
  func BenchmarkAdd(b *testing.B) {
      for i := 0; i < b.N; i++ {
          add(1, 1)
      }
  }
  ```
- **Negative Example**：
  ```go
  func BenchmarkAddWrong(b *testing.B) {
      for i := 0; i < 1000; i++ {
          add(1, 1)
      }
  }
  ```

### Concurrency Testing

- **Specification Content**：For concurrent code, appropriate test cases should be written to ensure the correctness of the concurrency logic.
- **Positive Example**：
  ```go
  func TestConcurrentAccess(t *testing.T) {
      // Set up concurrent environment ... 
      // Test logic for concurrent access ...
  }
  ```
- **Negative Example**：
  ```go
  func TestConcurrentAccess(t *testing.T) {
      // Only test single-thread logic...
  }
  ```

### Test Helper Functions

- **Specification Content**：Helper functions can be defined within the test files to help set up the test environment or clean up resources.
- **Positive Example**：
  ```go
  func setupTest(t *testing.T) {
      // Set up test environment ...
  }

  func tearDownTest(t *testing.T) {
      // Clean up resources ...
  }

  func TestMyFunction(t *testing.T) {
      t.Run("TestSetup", func(t *testing.T) {
          setupTest(t)
          // Test logic ...
      })
  }
  ```
- **Negative Example**：
  ```go
  // Directly setting up and cleaning up resources in the test
  func TestMyFunction(t *testing.T) {
      // Set up test environment... 
      // Test logic... 
      // Clean up resources...
  }
  ```

### Avoid Using Global Variables

- **Specification Content**: Try to avoid using global variables in tests to ensure test independence.
- **Positive Example**: Declare and use the necessary variables inside the test function.
- **Negative Example**: Declare global variables at the top of the test file.

### Clear Error Messages

- **Specification Content**: When a test fails, output clear and understandable error messages to help developers locate the problem.
- **Positive Example**: 
	- `t.Errorf("Expected value %d, but got %d", expected, real)`
- **Negative Example**: 
	- `t.Errorf("Error occurred")`
	- `fmt.Println("Error occurred")`
	- `panic("Error occurred")`

When a test fails, output clear and understandable error messages to help developers locate the problem.
