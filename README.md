# 🧑‍💻 lcode - Leetcode CLI Tool

![GitHub last commit](https://img.shields.io/github/last-commit/shadowmkj/lcode)
![GitHub issues](https://img.shields.io/github/issues/shadowmkj/lcode)
![GitHub forks](https://img.shields.io/github/forks/shadowmkj/lcode)
![GitHub stars](https://img.shields.io/github/stars/shadowmkj/lcode)
![License](https://img.shields.io/github/license/shadowmkj/lcode)

---

## 🚀 Overview

`lcode` is a powerful command-line interface (CLI) tool designed to streamline your Leetcode experience directly from your terminal. Say goodbye to constantly switching between your browser and editor. With `lcode`, you can effortlessly fetch Leetcode problems, work on them in your preferred terminal-based editor, and submit your solutions, all without leaving the comfort of your command line.

This tool is crafted for developers who love staying in their terminal environment, enhancing productivity, and making the Leetcode grind more efficient and enjoyable.

---

## ✨ Features

-   **Fetch Problems**: Quickly retrieve problem descriptions and starter code directly to your local machine. (coming soon)
-   **Submit Solutions**: Submit your code and get instant feedback without opening a browser.
-   **Test Locally**: (Planned) Integrate with local test cases for rapid iteration.
-   **Language Agnostic**: Work with your preferred programming language.

---

## 🛠 Installation

**(Installation instructions will go here once the tool is ready for distribution.)**

For now, you might need to build it from source:

** bat ** must be installed

```bash
git clone https://github.com/username/lcode.git
cd lcode
go build -o lcode main.go
# Move the executable to your PATH, e.g., /usr/local/bin
sudo mv lcode /usr/local/bin/
```

---

## 📖 Usage

**(Usage examples and commands will be detailed here.)**

```bash
# Example: Fetch a Leetcode problem by ID
lcode fetch 123
# Example: Fetch a Leetcode problem by name
lcode fetch two-sum

# Example: Submit your solution for a problem
lcode submit solution.go
```

---

## 💡 Recommendations for a Seamless Experience

`lcode` truly shines when integrated into a terminal-centric workflow. We highly recommend using it alongside:

-   **Neovim (or Vim)**: A highly configurable text editor that allows you to stay in the terminal for coding, testing, and debugging.
-   **Tmux (or other terminal multiplexers like Byobu, Zellij)**: Manage multiple terminal sessions within a single window, allowing you to have `lcode`, your editor, and a testing environment open simultaneously for maximum efficiency.

This setup creates a powerful and distraction-free environment for tackling Leetcode problems.

---

## 🤝 Contributing

We welcome contributions! If you have ideas for new features, bug fixes, or improvements, please open an issue or submit a pull request.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Developed with ❤️ by shadowmkj**
