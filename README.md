![gqlgen](https://user-images.githubusercontent.com/980499/133180111-d064b38c-6eb9-444b-a60f-7005a6e68222.png)


# gqlgen [![Integration](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/actions/workflows/integration.yml/badge.svg)](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/actions) [![Coverage Status](https://coveralls.io/repos/gitlab/nsx-allspark_users/nexus-sdk/gqlgen/badge.svg?branch=master)](https://coveralls.io/gitlab/nsx-allspark_users/nexus-sdk/gqlgen?branch=master) [![Go Report Card](https://goreportcard.com/badge/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git)](https://goreportcard.com/report/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git) [![Go Reference](https://pkg.go.dev/badge/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git.svg)](https://pkg.go.dev/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git) [![Read the Docs](https://badgen.net/badge/docs/available/green)](http://gqlgen.com/)

## What is gqlgen?

[gqlgen](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git) is a Go library for building GraphQL servers without any fuss.<br/>

- **gqlgen is based on a Schema first approach** — You get to Define your API using the GraphQL [Schema Definition Language](http://graphql.org/learn/schema/).
- **gqlgen prioritizes Type safety** — You should never see `map[string]interface{}` here.
- **gqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.

Still not convinced enough to use **gqlgen**? Compare **gqlgen** with other Go graphql [implementations](https://gqlgen.com/feature-comparison/)

## Quick start
1. [Initialise a new go module](https://golang.org/doc/tutorial/create-module)

       mkdir example
       cd example
       go mod init example

2. Add `gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git` to your [project's tools.go](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)

       printf '// +build tools\npackage tools\nimport _ "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git"' | gofmt > tools.go
       go mod tidy

3. Initialise gqlgen config and generate models

       go run gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git init

4. Start the graphql server

       go run server.go

More help to get started:
 - [Getting started tutorial](https://gqlgen.com/getting-started/) - a comprehensive guide to help you get started
 - [Real-world examples](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/tree/master/_examples) show how to create GraphQL applications
 - [Reference docs](https://pkg.go.dev/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git) for the APIs

## Reporting Issues

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/issues) on GitHub.

## Contributing

We welcome contributions, Read our [Contribution Guidelines](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/blob/master/CONTRIBUTING.md) to learn more about contributing to **gqlgen**
## Frequently asked questions

### How do I prevent fetching child objects that might not be used?

When you have nested or recursive schema like this:

```graphql
type User {
  id: ID!
  name: String!
  friends: [User!]!
}
```

You need to tell gqlgen that it should only fetch friends if the user requested it. There are two ways to do this;

- #### Using Custom Models

Write a custom model that omits the friends field:

```go
type User struct {
  ID int
  Name string
}
```

And reference the model in `gqlgen.yml`:

```yaml
# gqlgen.yml
models:
  User:
    model: github.com/you/pkg/model.User # go import path to the User struct above
```

- #### Using Explicit Resolvers

If you want to Keep using the generated model, mark the field as requiring a resolver explicitly in `gqlgen.yml` like this:

```yaml
# gqlgen.yml
models:
  User:
    fields:
      friends:
        resolver: true # force a resolver to be generated
```

After doing either of the above and running generate we will need to provide a resolver for friends:

```go
func (r *userResolver) Friends(ctx context.Context, obj *User) ([]*User, error) {
  // select * from user where friendid = obj.ID
  return friends,  nil
}
```

You can also use inline config with directives to achieve the same result

```graphql
directive @goModel(model: String, models: [String!]) on OBJECT
    | INPUT_OBJECT
    | SCALAR
    | ENUM
    | INTERFACE
    | UNION

directive @goField(forceResolver: Boolean, name: String) on INPUT_FIELD_DEFINITION
    | FIELD_DEFINITION

type User @goModel(model: "github.com/you/pkg/model.User") {
    id: ID!         @goField(name: "todoId")
    friends: [User!]!   @goField(forceResolver: true)
}
```

### Can I change the type of the ID from type String to Type Int?

Yes! You can by remapping it in config as seen below:

```yaml
models:
  ID: # The GraphQL type ID is backed by
    model:
      - gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql.IntID # a go integer
      - gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql.ID # or a go string
```

This means gqlgen will be able to automatically bind to strings or ints for models you have written yourself, but the
first model in this list is used as the default type and it will always be used when:

- Generating models based on schema
- As arguments in resolvers

There isn't any way around this, gqlgen has no way to know what you want in a given context.

## Other Resources

- [Christopher Biscardi @ Gophercon UK 2018](https://youtu.be/FdURVezcdcw)
- [Introducing gqlgen: a GraphQL Server Generator for Go](https://99designs.com.au/blog/engineering/gqlgen-a-graphql-server-generator-for-go/)
- [Dive into GraphQL by Iván Corrales Solera](https://medium.com/@ivan.corrales.solera/dive-into-graphql-9bfedf22e1a)
- [Sample Project built on gqlgen with Postgres by Oleg Shalygin](https://github.com/oshalygin/gqlgen-pg-todo-example)
- [Hackernews GraphQL Server with gqlgen by Shayegan Hooshyari](https://www.howtographql.com/graphql-go/0-introduction/)
