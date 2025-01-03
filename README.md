# copybooktogo

`copybooktogo` is a command-line tool that converts COBOL copybooks to Go struct definitions. It simplifies the process of migrating COBOL data structures to Go by handling the normalization, parsing, and generation of equivalent Go code.

## Installation

```bash
go install github.com/yasv98/copybooktogo@latest
```

## Usage

Basic usage:
```bash
copybooktogo -c path/to/copybook.cpy
```

### Command Line Options

- `-c, --copybook` (required): Path to the COBOL copybook file to convert
- `-p, --package` (optional): Package name for the generated Go code (default: "main")
- `-t, --typeOverrides` (optional): Custom type mapping overrides in from=to format
- `-o, --output` (optional): Path to the output file or directory

### Type Overrides

The `-t, --typeOverrides` flag allows you to customize how COBOL PIC types are mapped to Go types. Use a comma-separated list of mappings in the format `cobolType=goType`. For example:

```bash
copybooktogo -c data.cpy -t "unsigned=int,decimal=custom.DecimalType"
```

This will override the default type mappings for unsigned and decimal types in the generated code.

### Examples

Convert a copybook using default settings:
```bash
copybooktogo -c data.cpy
```

Specify a custom package name and output location:
```bash
copybooktogo -c data.cpy -p models -o ./generated/
```

Override specific type conversions:
```bash
copybooktogo -c data.cpy -t "unsigned=int,decimal=custom.Type"
```

## Notes

- The tool will automatically handle COBOL copybook normalization and Go code generation
- Generated structs will follow Go naming conventions
- Type mappings can be customized to match your specific requirements