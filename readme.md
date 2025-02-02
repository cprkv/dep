# dep

Dependency management through git cloning.

## Usage

In root of project add file `repository.xml`

```xml
<repository>
  <feature name="main">
    <dependency name="sdl" url="https://github.com/libsdl-org/SDL.git" revision="main" />
    <dependency name="glm" url="https://github.com/g-truc/glm.git" />
  </feature>
</repository>
```

Run `dep` command in terminal. It will sync dependencies from `repository.xml` into `deps` folder with their respective names as directory names.

The tool does sync recursive. If it meets any `repository.xml` in root of dependency, it will clone all its dependencies into topmost `deps` folder.

## Build

### prepare environment:

- install [go](https://go.dev/doc/install)
- install [inno setup](https://jrsoftware.org/isdl.php)

### actual build:

- run `./build-installer.ps1`

## Install

Run installer: `./release/installer.exe`
