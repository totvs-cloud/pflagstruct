# pflagstruct

[![License](https://img.shields.io/github/license/saltstack/salt)](https://opensource.org/license/apache-2-0/)

A code generation tool that simplifies the process of registering command line flags in Go applications

## Installation

To install you need to have Go installed and set up on your machine. Then, you can use the following command to install
the tool:

```shell
go install github.com/totvs-cloud/pflagstruct@latest
```

## Usage

The CLI Code Generation Tool provides the following flags:

- `--destination string`: Specifies the path where the generated code will be saved. If not provided, the current path
  is used.
- `--directory string`: Specifies the path where the source file containing the struct definition is located. This flag
  is required if `--package` is not informed.
- `--package string`: Specifies the package path of the struct definition. This flag is required if `--directory` is not
  informed.
- `--struct-name string`: Specifies the name of the struct. This flag is required.

## Examples

1. Generate code using a struct definition in a directory:
   ```shell
   pflagstruct --directory /path/to/source --struct-name MyStruct
   ```

2. Generate code using a struct definition in a package:
   ```shell
   pflagstruct --package github.com/example/package --struct-name MyStruct
   ```

3. Generate code with a custom destination path:
   ```shell
   pflagstruct --destination /path/to/destination --package github.com/example/package --struct-name MyStruct
   ```

4. Automate the code generation process using the "//go:generate" comment in your Go source files:
    ```go
    //go:generate pflagstruct --struct-name=User --package=github.com/example/model
    
    package main
    
    import (
        "context"
    
        "github.com/spf13/cobra"
        "github.com/example/client"
    )
    
    func NewCommand() *cobra.Command {
        cmd := &cobra.Command{
            Use:   "create",
            Short: "Create a new user.",
            RunE: func(cmd *cobra.Command, args []string) error {
                // Get user from command flags
                user, err := GetUserFromFlags(cmd.Flags()) // generated method
                if err != nil {
                    return err
                }
    
                // Create a new user
                return client.CreateUser(user)
            },
        }
        SetUpUserToFlags(cmd.Flags()) // generated method
    
        return cmd
    }
    ```
   For the referenced struct, as shown below:
    ```go
    package model
    
    type User struct {
        // Unique identifier of the user.
        ID string
        // Name of the user.
        Name string
        // Email address of the user.
        Email string
    }
    ```
   The generated code appears as follows:
    ```go
    package main
    
    import (
        "fmt"
        "github.com/spf13/pflag"
        "github.com/example/model"
    )
    
    func SetUpUserToFlags(flags *pflag.FlagSet) {
        (&userFlagsBuilder{flags: flags}).setUpUser()
    }

    func GetUserFromFlags(flags *pflag.FlagSet) (*model.User, error) {
        return (&userFlagsBuilder{flags: flags}).getUser()
    }
    
    type userFlagsBuilder struct { flags *pflag.FlagSet }
    
    func (cf *userFlagsBuilder) setUpUser() {
        cf.flags.String("id", "", "Unique identifier of the user.")
        cf.flags.String("name", "", "Name of the user.")
        cf.flags.String("email", "", "Email address of the user.")
    }

    func (cf *userFlagsBuilder) getUser() (user *model.User, err error) {
        user = new(model.User)
        if user.ID, err = cf.flags.GetString("id"); err != nil {
            return nil, fmt.Errorf("error retrieving \"id\" from command flags: %w", err)
        }
        if user.Name, err = cf.flags.GetString("name"); err != nil {
            return nil, fmt.Errorf("error retrieving \"name\" from command flags: %w", err)
        }
        if user.Email, err = cf.flags.GetString("email"); err != nil {
            return nil, fmt.Errorf("error retrieving \"email\" from command flags: %w", err)
        }
        return user, nil
    }
    ```

Feel free to explore the available flags and experiment with different options to generate code based on your struct
definitions.

## Contributing

Contributions to the CLI Code Generation Tool are welcome! If you find any issues or have suggestions for improvement,
please open an issue or submit a pull request.
