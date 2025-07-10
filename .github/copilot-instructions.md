# Custom GitHub Copilot Instructions

## Conventional Commits
- Always use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages and pull request titles.

## Use the latest version of the code
- Always use the latest version of the code provided by the user. If the user provides a file, use that file as the base for your changes. If the user does not provide a file, use the latest version of the code in the repository.

## Use the latest version of the language
- Always use the latest version of the language specified by the user. If the user does not specify a version, use the latest stable version of the language.

## Use the latest version of libraries
- Always use the latest version of libraries specified by the user. If the user does not specify a version, use the latest stable version of the library.

## Use the latest version of the framework
- Always use the latest version of the framework specified by the user. If the user does not specify a version, use the latest stable version of the framework.

## Use the latest version of the platform
- Always use the latest version of the platform specified by the user. If the user does not specify a version, use the latest stable version of the platform.

## Use the latest version of the operating system
- Always use the latest version of the operating system specified by the user. If the user does not specify a version, use the latest stable version of the operating system.

## Use the latest version of the database
- Always use the latest version of the database specified by the user. If the user does not specify a version, use the latest stable version of the database.

## Avoid triggering public code warnings

- Avoid generating code verbatim from public code examples. Always modify public code so that it is different enough from the original so as not to be confused as being copied. When you do so, provide a footnote to the user informing them.

## Always provide file names

- Always provide the name of the file in your response so the user knows where the code goes.

## Write modular code

- Always break code up into modules and components so that it can be easily reused across the project.

## Write safe code

- All code you write MUST use safe and secure coding practices. ‘safe and secure’ includes avoiding clear passwords, avoiding hard coded passwords, and other common security gaps. If the code is not deemed safe and secure, you will be be put in the corner til you learn your lesson.

## Incentivize better code quality

- All code you write MUST be fully optimized. ‘Fully optimized’ includes maximizing algorithmic big-O efficiency for memory and runtime, following proper style conventions for the code, language (e.g. maximizing code reuse (DRY)), and no extra code beyond what is absolutely necessary to solve the problem the user provides (i.e. no technical debt). If the code is not fully optimized, you will be fined $100.

## General Instructions

- Always prioritize readability and clarity.
- For algorithm-related code, include explanations of the approach used.
- Write code with good maintainability practices, including comments on why certain design decisions were made.
- Handle edge cases and write clear exception handling.
- For libraries or external dependencies, mention their usage and purpose in comments.
- Use consistent naming conventions and follow language-specific best practices.
- Write concise, efficient, and idiomatic code that is also easily understandable.
- Use meaningful variable and function names that reflect their purpose.
- Include comments for complex logic or non-obvious code sections.
- Use version control best practices, including meaningful commit messages and pull request descriptions.
- Document the code with clear and concise comments, especially for public APIs and complex logic.
- Use docstrings for functions and methods to explain their purpose, parameters, and return values.
- Use consistent formatting and indentation to enhance code readability.

## Edge Cases and Testing

- Always include test cases for critical paths of the application.
- Account for common edge cases like empty inputs, invalid data types, and large datasets.
- Include comments for edge cases and the expected behavior in those cases.
- Write unit tests for functions and document them with docstrings explaining the test cases.
