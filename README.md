
# Go Git Swap 🔄

[![Go Report Card](https://goreportcard.com/badge/github.com/frivas/go-git-swap)](https://goreportcard.com/report/github.com/frivas/go-git-swap)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Note**: This is my first Go project! I'm actively learning and improving it. Feedback and contributions are very welcome!

Go Git Swap is a command-line tool that helps developers manage multiple Git configurations easily. Switch between different Git profiles seamlessly, perfect for developers who work with multiple Git accounts (personal, work, client projects, etc.).

## 🗺️ Roadmap

### Coming Soon (Q2 2024)
- [ ] Modern TUI interface using [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [ ] Enhanced SSH signing key support
- [ ] Profile templates and quick switching
- [ ] Auto-detection of Git repositories

### Future Plans (Q3-Q4 2024)
- [ ] Profile groups for different contexts (work, personal, client)
- [ ] Git hooks management
- [ ] Profile sharing and import/export
- [ ] Integration with popular Git hosting platforms
- [ ] Configuration backup and sync

## Features ✨

Current features:
- Manage multiple Git profiles
- Easy switching between profiles
- Basic SSH signing key support
- Automatic configuration of Git global settings
- Simple CLI interface
- Cross-platform support (macOS, Linux)
- Secure local configuration storage

## Installation 🚀

### Using Go

```bash
go install github.com/frivas/go-git-swap@latest
```

### From Releases

Download the latest release for your platform from the [releases page](https://github.com/frivas/go-git-swap/releases).

### Building from Source

```bash
# Clone the repository
git clone https://github.com/frivas/go-git-swap.git
cd go-git-swap

# Build for your platform
make build

# Or build for all supported platforms
make build-all
```

## Usage 💡

1. Start the application:
```bash
go-git-swap
```

2. Create a new profile:
   - Select option 2
   - Enter profile details (name, email, username, etc.)
   - Optionally add an SSH signing key

3. Switch between profiles:
   - Select option 1
   - Choose the profile you want to activate

4. Delete a profile:
   - Select option 3
   - Choose the profile you want to remove

## Configuration 🔧

The application stores its configuration in `~/.go-git-swap.json`. Each profile contains:

- Profile name (identifier)
- Git username
- Git email
- Full name
- SSH signing key path (optional)

## Development 🛠️

### Prerequisites

- Go 1.20 or higher
- Make (for building)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Project Structure

```
.
├── cmd/
│   └── go-git-swap/       # Main application
├── internal/
│   ├── config/           # Configuration management
│   ├── git/             # Git operations
│   ├── model/           # Data models
│   └── validator/       # Input validation
├── Makefile
└── README.md
```

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Code Style

Please follow these guidelines:
- Use `gofmt` to format your code
- Add comments for exported functions
- Write tests for new functionality
- Update documentation as needed

## License 📝

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments 🙏

- Thanks to all contributors who have helped shape Go Git Swap
- Built with ❤️ using Go

## Support 💪

If you find this project useful, please consider giving it a ⭐️ on GitHub!

For bugs, feature requests, or questions, please [open an issue](https://github.com/frivas/go-git-swap/issues).
# Go Git Swap 🔄

[![Go Report Card](https://goreportcard.com/badge/github.com/frivas/go-git-swap)](https://goreportcard.com/report/github.com/frivas/go-git-swap)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go Git Swap is a command-line tool that helps developers manage multiple Git configurations easily. Switch between different Git profiles seamlessly, perfect for developers who work with multiple Git accounts (personal, work, client projects, etc.).

## Features ✨

- Manage multiple Git profiles
- Easy switching between profiles
- Support for SSH signing keys
- Automatic configuration of Git global settings
- Interactive CLI interface
- Cross-platform support (macOS, Linux)
- Secure local configuration storage

## Installation 🚀

### Using Go

```bash
go install github.com/frivas/go-git-swap@latest
```

### From Releases

Download the latest release for your platform from the [releases page](https://github.com/frivas/go-git-swap/releases).

### Building from Source

```bash
# Clone the repository
git clone https://github.com/frivas/go-git-swap.git
cd go-git-swap

# Build for your platform
make build

# Or build for all supported platforms
make build-all
```

## Usage 💡

1. Start the application:
```bash
go-git-swap
```

2. Create a new profile:
   - Select option 2
   - Enter profile details (name, email, username, etc.)
   - Optionally add an SSH signing key

3. Switch between profiles:
   - Select option 1
   - Choose the profile you want to activate

4. Delete a profile:
   - Select option 3
   - Choose the profile you want to remove

## Configuration 🔧

The application stores its configuration in `~/.go-git-swap.json`. Each profile contains:

- Profile name (identifier)
- Git username
- Git email
- Full name
- SSH signing key path (optional)

## Development 🛠️

### Prerequisites

- Go 1.20 or higher
- Make (for building)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Project Structure

```
.
├── cmd/
│   └── go-git-swap/       # Main application
├── internal/
│   ├── config/           # Configuration management
│   ├── git/             # Git operations
│   ├── model/           # Data models
│   └── validator/       # Input validation
├── Makefile
└── README.md
```

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Code Style

Please follow these guidelines:
- Use `gofmt` to format your code
- Add comments for exported functions
- Write tests for new functionality
- Update documentation as needed

## License 📝

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments 🙏

- Thanks to all contributors who have helped shape Go Git Swap
- Built with ❤️ using Go

## Support 💪

If you find this project useful, please consider giving it a ⭐️ on GitHub!

For bugs, feature requests, or questions, please [open an issue](https://github.com/frivas/go-git-swap/issues).
