# Cisconf

Cisconf is a Go-based library for unmarshalling, marshalling, and comparing Cisco network configurations. It supports various Cisco configuration components such as VLANs, interfaces, OSPF, and EIGRP.

## Features

- Parse Cisco configuration files into Go structs.
- Compare two Cisco configurations and generate a diff.
- Supported components:
    - VLANs
    - Interfaces (Layer 2 and Layer 3)
    - OSPF
    - EIGRP

## Usage

### Parsing a Configuration

Use the `cisconf.Unmarshal` function to parse a Cisco configuration string into a Go struct.

```go
var vlan cisconf.Vlan
config := "vlan 300\n name office\n!"
err := cisconf.Unmarshal(config, &vlan)
if err != nil {
    panic(err)
}
```

### Generating a Configuration

Use the `cisconf.Marshal` function to generate a Cisco configuration string from a Go struct.

```go
generated, err := cisconf.Marshal(vlan)
if err != nil {
    panic(err)
}
```

### Comparing Configurations

Use the `cisconf.Diff` function to compare two Cisco configuration strings and generate a diff.

```go
src := cisconf.Vlan{Id: 300, Name: "office"}
dest := cisconf.Vlan{Id: 300, Name: "new_office"}
diff, err := cisconf.Diff(src, dest)
if err != nil {
    panic(err)
}
fmt.Println(diff)
```

## Running Tests

Run the test suite to validate the library's functionality:

```bash
go test ./test
```

## Credits
This project is built upon the work of Johannes Bindriem in [Jap](https://github.com/Letsu/jap).

## Contributing

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a detailed description of your changes.
