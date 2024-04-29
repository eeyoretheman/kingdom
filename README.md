# Kingdom

## Description

Kingdom is a powerful Go-based tool inspired by Empire PowerShell. It provides a platform for post-exploitation activities and offensive security operations.

## Features

- **Command and Control (C2) Framework**: Kingdom offers a flexible and extensible C2 framework for managing compromised systems.
- **Post-Exploitation Modules**: It includes a wide range of post-exploitation modules to perform various tasks on compromised systems.
<!-- - **Multi-Platform Support**: Kingdom supports multiple platforms, including Windows, Linux, and macOS. -->
<!-- - **Stealthy Operations**: It focuses on stealthy operations to evade detection and maintain persistence on compromised systems. -->
- **Extensibility**: Kingdom is designed to be easily extensible, allowing you to create your own modules and plugins.

## Installation

To install Kingdom, follow these steps:

1. Clone the repository:

    ```shell
    git clone https://github.com/eeyoretheman/kingdom.git
    ```

2. Build the project:

    ```shell
    cd kingdom
    go build
    ```

3. Run Kingdom:

    ```shell
    ./bin/kingdom
    ```

For detailed installation instructions and usage examples, please refer to the [documentation](https://github.com/eeyoretheman/kingdom/wiki).

## Usage

After starting Kingdom, you can interact with the C2 server using the command-line interface (CLI). Here are some examples of the available commands:

- **Teller**: Create a new listener:

    ```shell
    ! ! tl 192.168.88.185:8080
    ```
- **Client**: Create a client:

    ```shell
    ! ! cl localhost:1337
    ```
- **List Tellers**: List all listeners:

    ```shell
    ! ! lst !
    ```
- **List Clients**: List all clients:

    ```shell
    ! ! lsc !
    ```
- **Remove Teller**: Remove a listener:

    ```shell
    ! ! rmt <teller_id>
    ```
- **Remove Client**: Remove a client:

    ```shell
    ! ! rmc <client_id>
    ```
- **Lock**: Lock a teller to a client:

    ```shell
    <client_id> <teller_id> lock !
    ```
- **Unlock**: Unlock a teller from a client:

    ```shell
    <client_id> <teller_id> unlock !
    ```
- **Execute**: Execute a command on a client:

    ```shell
    <client_id> <teller_id> send <command>
    ```

## Contributing

There is no support for contributing :(.

## License

Kingdom is licensed under the [MIT License](https://github.com/your-username/kingdom/LICENSE).

## Contact

For any questions or feedback, feel free to reach out to us at [your-email@example.com](mailto:your-email@example.com).
