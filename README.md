# Warp Code Injector (WCI)

**Warp Code Injector (WCI)** is a streamlined CLI tool crafted to inject custom Lua scripts into Factorio savegames.
With WCI, seamlessly enhance your savegame files while retaining achievements, making it the perfect utility for players
and modders alike.

---

## 🚀 Features

- 🔄 **Inject Lua Scripts**: Modify savegame files by adding predefined Lua scripts, such as the **Biter Killer**.
- 📂 **Savegame Management**: Effortlessly list and manage savegames sorted by creation date.
- ⚡ **Cross-Platform Support**: Designed for **Windows** and **macOS**, ensuring a smooth experience for all users.

---

## 🌟 Scripts Currently Available

### **🪓 Biter Killer**

- **Purpose**: Completely removes biters from the savegame world *without deactivating achievements*!
- **Command in Game**: `/cleanup_biters`
- **Script Location**: [biter-killer.lua](embedded/lua_injections/biter_killer.lua)

---

## 📚 Table of Contents

1. [✨ Features](#features)
2. [🌟 Scripts Currently Available](#scripts-currently-available)
3. [💻 Usage](#usage)
4. [🛠️ Technologies Used](#technologies-used)
5. [🚀 Future Enhancements](#future-enhancements)
6. [📜 License](#license)
7. [🤝 Contributors](#contributors)

---

## 💻 Usage

### **Basic Commands**

#### **1. List Savegames**

```bash
wci list
```

Displays a list of all savegames in the Factorio save directory, sorted by creation date.

#### **2. Inject Lua Script**

```bash
wci add-biter-killer [number-of-save-from-list-command]
```

Injects the **Biter Killer** Lua script into the specified savegame.

#### **3. Clean Temporary Files**

```bash
wci clean
```

Deletes WCI-generated files (`savegames.json`) from the executable's directory.

---

## 🛠️ Technologies Used

- **Go (Golang)**: Core programming language.
- **spf13 Cobra**: Robust CLI framework for command management.
- **Zerolog**: Lightweight, high-performance logging library.

---

## 🚀 Future Enhancements

1. **Advanced Lua Features**:
    - [ ] Validate Lua scripts before injection.
    - [ ] Enable template-based script creation.
    - [ ] Allow injecting custom scripts into user-defined locations.
2. **Savegame Enhancements**:
    - [ ] Add features for backup and restore.
3. **More Predefined Scripts**:
    - [ ] Add scripts for other gameplay modifications.

---

## 📜 License

This project is licensed under the **GNU Affero General Public License v3.0**.  
For more details, see the [LICENSE](LICENSE.md) file.

---

## 🤝 Contributors

We welcome contributions from the community! Whether it's:

- 🛠️ Improving the codebase.
- 📝 Suggesting new features or predefined Lua scripts.
- 🐛 Reporting issues or bugs.

Feel free to [create a pull request](https://github.com/KnightRider2070/Warp-Code-Injector/pulls)
or [open an issue](https://github.com/KnightRider2070/Warp-Code-Injector/issues).  
Together, let's make Factorio modding even better! 💡
