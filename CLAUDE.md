# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Text Processors is a Go desktop application built with Fyne that provides text processing tools with an extensible custom processor system. The application features built-in tools (JSON Formatter, Text Splitter, Text Joiner) and allows users to create custom JavaScript processors for advanced text transformations.

## Build Commands

### Primary Build Command
```bash
# Windows batch script (recommended)
./build.bat

# Release mode (optimized binary)
./build.bat release

# Direct Go build
go build -o build/text_processors.exe .

# Linux/macOS
go build -o build/text_processors .
```

### Development Commands
```bash
# Run without building
go run .

# Update dependencies
go mod tidy

# Get dependencies
go mod download
```

## Architecture

### Core Structure
- **`main.go`**: Application entry point, initializes Fyne app and loads custom processors
- **`ui/`**: All UI components and business logic
- **`ui/processors/`**: Built-in text processors and processor interfaces
- **`build/`**: Compiled binaries output directory

### Key Packages
- **Fyne v2.6.1**: GUI framework
- **Chroma v2.14.0**: Syntax highlighting for JavaScript editor
- **Goja**: JavaScript engine for custom processors

### UI Architecture
The application uses a navigation-based UI structure:

1. **Main Grid (`tools_grid.go`)**: Tool selection interface
2. **Individual Tool UIs**: Each processor has its own UI implementation
3. **Pipeline Builder (`pipeline_builder.go`)**: Allows chaining multiple processors
4. **Custom Processor System**: JavaScript-based extensible processors

### Processor System
Processors implement a common interface defined in `ui/processors/processor.go`. The system supports:
- Built-in processors (JSON Formatter, Text Splitter, Text Joiner)
- Custom JavaScript processors with full JS runtime
- Pipeline composition of multiple processors

### Custom Processors
Custom processors are:
- Written in JavaScript with a `process(input)` function
- Executed using the Goja JavaScript engine
- Automatically saved/loaded from the `conf/` directory
- Integrated into pipelines alongside built-in processors

### File Organization Patterns
- Each processor has its own UI file: `ui/processors/[processor_name].go`
- UI constructors follow pattern: `Make[ProcessorName]UI() fyne.CanvasObject`
- Custom processors use a manager pattern with persistence

### Navigation Pattern
All tool UIs follow a consistent pattern:
1. Input area (text entry)
2. Configuration options
3. Process button
4. Output area (with copy functionality)
5. Back button for navigation

## Development Guidelines

### Adding New Built-in Processors
1. Create `ui/processors/new_processor.go` with `MakeNewProcessorUI()` function
2. Add processor to the tools grid in `ui/tools_grid.go`
3. Follow the input→process→output→copy pattern
4. Use proper error handling and validation

### Code Conventions
- Public functions start with capital letters (Go convention)
- UI constructors return `fyne.CanvasObject`
- Error handling should be graceful (no crashes on user input)
- Use `container.NewBorder()` for responsive layouts

### Custom Processor Integration
Custom processors automatically integrate with:
- Pipeline Builder (appear with "Custom:" prefix)
- Export/import system (included in pipeline configurations)
- Main processor grid