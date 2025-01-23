# txtnest-cli  

Welcome to `txtnest-cli`! This project is all about giving you a slick, terminal-based way to interact with the TxtNest backend. Whether you're running commands locally through the CLI or connecting remotely via SSH, TxtNest has you covered.  

The idea for this type of terminal-to-backend interaction was inspired by the awesome [terminal.shop](https://terminal.shop) project. It showed just how cool and efficient terminal-based workflows can be, and TxtNest takes that concept and runs with it—focused on managing and interacting with text entries.  

## Features  

- **CLI Tool**: Manage your text entries directly from your terminal.  
- **SSH Server**: Spin up an SSH server to interact remotely.  
- **Smooth Workflow**: Easily toggle between local CLI use and SSH-based management.  

## Getting Started  

### Prerequisites  
First, make sure you’ve got [Go](https://golang.org/dl/) installed.  

### Installation  
1. Clone the repository:  
   ```bash  
   git clone https://github.com/NesterovYehor/txtnest-cli.git  
   cd txtnest-cli  
   ```  
2. Build the project:  
   ```bash  
   make build  
   ```  

## Usage  

You’ve got two main ways to work with TxtNest:  

- **Run the CLI Tool**:  
   ```bash  
   make run-cli  
   ```  
   This runs the CLI tool so you can work with TxtNest directly in your terminal.  

- **Run the SSH Server**:  
   ```bash  
   make run-ssh  
   ```  
   This starts the SSH server, letting you connect remotely via an SSH client and manage text entries.  

## Makefile Commands  

Here’s a quick rundown of the `make` commands you’ll use:  

- `make build` – Build the project and put the binaries in the `build/` folder.  
- `make run-cli` – Run the CLI tool.  
- `make run-ssh` – Start the SSH server.  
- `make clean` – Clean up the build directory.  

## Why TxtNest?  

TxtNest is built for developers who love the terminal and want an efficient way to manage their backend. If you’ve checked out [terminal.shop](https://terminal.shop), you already know how effective this kind of interaction can be. TxtNest takes that inspiration and fine-tunes it for managing text data.  

## Contributing  

Have an idea or found a bug? Feel free to fork the repo, make your changes, and submit a pull request. Every bit helps, and contributions are always welcome!  

## License  

This project is licensed under the MIT License. Check out the [LICENSE](./LICENSE) file for more details.  

---  

For updates and more info, visit the [txtnest-cli GitHub repository](https://github.com/NesterovYehor/txtnest-cli).  
