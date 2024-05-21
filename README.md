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

![Kingdom CLI](https://placehold.co/600x400)

<!--Global:
Ctrl + U - go to the teller menu
Ctrl + Y - go to the input menu
Ctrl + C - exit

Input menu:
Left arrow - Move cursor left
Right arrow - Move cursor right
Up arrow - previous command
Down arrow - next command
Ctrl + J - scroll history down
Ctrl + K - scroll history up
Enter - submit command

Teller menu:
Up arrow - scroll one teller up
Down arrow - scroll one teller down
Enter -  select teller-->

### Controls

- **Global**:
  - `Ctrl + U`: Go to the teller menu
  - `Ctrl + Y`: Go to the input menu
  - `Ctrl + C`: Exit

- **Input Menu**:
    - `Left Arrow`: Move cursor left
    - `Right Arrow`: Move cursor right
    - `Up Arrow`: Previous command
    - `Down Arrow`: Next command
    - `Ctrl + J`: Scroll history down
    - `Ctrl + K`: Scroll history up
    - `Enter`: Submit command

- **Teller Menu**:
    - `Up Arrow`: Scroll one teller up
    - `Down Arrow`: Scroll one teller down
    - `Enter`: Select teller

### Examples

- **Create a teller**:

    ```shell
    tl <ip>:<port>
    ```

- **Creaye a client**:

    ```shell
    cl <ip>:<port>
    ```

- **List tellers**:

    ```shell
    lst
    ```

- **List clients**:

    ```shell
    lsc
    ```

- **Remove teller**:

    ```shell
    rmt <teller_id>
    ```

- **Remove client**:

    ```shell
    rmc <client_id>
    ```

- **Lock a client to a teller**:

    ```shell
    lock
    ```

- **Unlock a client from a teller**:

    ```shell
    unlock
    ```

- **Execute a command on a client**:

    ```shell
    send <command/macro>
    ```

## Contributing

There is no support for contributing :(.

## Contact

For any questions or feedback, feel free to reach out to us at [your-email@example.com](mailto:your-email@example.com).
