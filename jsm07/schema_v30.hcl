
  schema = "http://json-schema.org/draft-04/schema#"
  id = "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
  patternProperties = {
    "^x-" = {}
  }
  description = "The description of OpenAPI v3.0.x documents, as defined by https://spec.openapis.org/oas/v3.0.3"
  type = "object"
  required = ["openapi", "info", "paths"]
  properties = {
    servers = array(definitions.Server),
    security = array(definitions.SecurityRequirement),
    tags = array(definitions.Tag, uniqueItems(true)),
    paths = definitions.Paths,
    components = definitions.Components,
    openapi = string(pattern("^3\\.0\\.\\d(-.+)?$")),
    info = definitions.Info,
    externalDocs = definitions.ExternalDocumentation
  }
  additionalProperties = false
  definitions {
    SecurityRequirement = map(array(string()))
    Discriminator = object({
      propertyName = string(),
      mapping = map(string())
    }, required(["propertyName"]))
    Schema {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      properties = {
        minLength = integer(default(0), minimum(0)),
        maximum = number(),
        type = string(enum("array", "boolean", "integer", "number", "object", "string")),
        multipleOf = number(minimum(0), exclusiveMinimum(true)),
        maxItems = integer(minimum(0)),
        exclusiveMaximum = boolean(default(false)),
        oneOf = array({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        nullable = boolean(default(false)),
        format = string(),
        pattern = string(format("regex")),
        deprecated = boolean(default(false)),
        allOf = array({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        externalDocs = definitions.ExternalDocumentation,
        maxProperties = integer(minimum(0)),
        maxLength = integer(minimum(0)),
        anyOf = array({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        title = string(),
        minProperties = integer(default(0), minimum(0)),
        description = string(),
        readOnly = boolean(default(false)),
        minItems = integer(default(0), minimum(0)),
        uniqueItems = boolean(default(false)),
        properties = map({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        xml = definitions.XML,
        discriminator = definitions.Discriminator,
        writeOnly = boolean(default(false)),
        exclusiveMinimum = boolean(default(false)),
        enum = array({}, minItems(1), uniqueItems(false)),
        minimum = number(),
        required = array(string(), minItems(1), uniqueItems(true)),
        not = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        items = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        additionalProperties = {
          default = true,
          oneOf = [definitions.Schema, definitions.Reference, boolean()]
        },
        example = {},
        default = {}
      }
    }
    MediaType {
      patternProperties = {
        "^x-" = {}
      }
      allOf = [definitions.ExampleXORExamples]
      type = "object"
      properties = {
        examples = map({
          oneOf = [definitions.Example, definitions.Reference]
        }),
        encoding = map(definitions.Encoding),
        schema = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        example = {}
      }
      additionalProperties = false
    }
    Header {
      allOf = [definitions.ExampleXORExamples, definitions.SchemaXORContent]
      type = "object"
      properties = {
        deprecated = boolean(default(false)),
        examples = map({
          oneOf = [definitions.Example, definitions.Reference]
        }),
        allowEmptyValue = boolean(default(false)),
        style = string(default("simple"), enum("simple")),
        explode = boolean(),
        allowReserved = boolean(default(false)),
        description = string(),
        required = boolean(default(false)),
        schema = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        content = {
          type = "object",
          maxProperties = 1,
          minProperties = 1,
          additionalProperties = definitions.MediaType
        },
        example = {}
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Contact {
      type = "object"
      properties = {
        name = string(),
        url = string(format("uri-reference")),
        email = string(format("email"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Response {
      type = "object"
      required = ["description"]
      properties = {
        content = map(definitions.MediaType),
        links = map({
          oneOf = [definitions.Link, definitions.Reference]
        }),
        description = string(),
        headers = map({
          oneOf = [definitions.Header, definitions.Reference]
        })
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    ExternalDocumentation {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["url"]
      properties = {
        description = string(),
        url = string(format("uri-reference"))
      }
    }
    SecurityScheme {
      oneOf = [definitions.APIKeySecurityScheme, definitions.HTTPSecurityScheme, definitions.OAuth2SecurityScheme, definitions.OpenIdConnectSecurityScheme]
    }
    Operation {
      type = "object"
      required = ["responses"]
      properties = {
        summary = string(),
        description = string(),
        externalDocs = definitions.ExternalDocumentation,
        operationId = string(),
        security = array(definitions.SecurityRequirement),
        responses = definitions.Responses,
        deprecated = boolean(default(false)),
        parameters = array({
          oneOf = [definitions.Parameter, definitions.Reference]
        }, uniqueItems(true)),
        callbacks = map({
          oneOf = [definitions.Callback, definitions.Reference]
        }),
        tags = array(string()),
        servers = array(definitions.Server),
        requestBody = {
          oneOf = [definitions.RequestBody, definitions.Reference]
        }
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    PathParameter {
      required = ["required"]
      properties = {
        in = {
          enum = ["path"]
        },
        style = {
          default = "simple",
          enum = ["matrix", "label", "simple"]
        },
        required = {
          enum = [true]
        }
      }
      description = "Parameter in path"
    }
    QueryParameter {
      properties = {
        in = {
          enum = ["query"]
        },
        style = {
          default = "form",
          enum = ["form", "spaceDelimited", "pipeDelimited", "deepObject"]
        }
      }
      description = "Parameter in query"
    }
    CookieParameter {
      properties = {
        in = {
          enum = ["cookie"]
        },
        style = {
          enum = ["form"],
          default = "form"
        }
      }
      description = "Parameter in cookie"
    }
    OpenIdConnectSecurityScheme {
      properties = {
        type = string(enum("openIdConnect")),
        openIdConnectUrl = string(format("uri-reference")),
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["type", "openIdConnectUrl"]
    }
    ClientCredentialsFlow {
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["tokenUrl", "scopes"]
      properties = {
        refreshUrl = string(format("uri-reference")),
        scopes = map(string()),
        tokenUrl = string(format("uri-reference"))
      }
      additionalProperties = false
    }
    Paths {
      type = "object"
      additionalProperties = false
      patternProperties = {
        "^/" = definitions.PathItem,
        "^x-" = {}
      }
    }
    Responses {
      type = "object"
      minProperties = 1
      properties = {
        default = {
          oneOf = [definitions.Response, definitions.Reference]
        }
      }
      additionalProperties = false
      patternProperties = {
        "^[1-5](?:\\d{2}|XX)$" = {
          oneOf = [definitions.Response, definitions.Reference]
        },
        "^x-" = {}
      }
    }
    OAuthFlows {
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      properties = {
        implicit = definitions.ImplicitOAuthFlow,
        password = definitions.PasswordOAuthFlow,
        clientCredentials = definitions.ClientCredentialsFlow,
        authorizationCode = definitions.AuthorizationCodeOAuthFlow
      }
      additionalProperties = false
    }
    Callback {
      type = "object"
      additionalProperties = definitions.PathItem
      patternProperties = {
        "^x-" = {}
      }
    }
    Server {
      type = "object"
      required = ["url"]
      properties = {
        url = string(),
        description = string(),
        variables = map(definitions.ServerVariable)
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Components {
      type = "object"
      properties = {
        schemas = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Schema, definitions.Reference]
            }
          }
        },
        responses = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Response]
            }
          }
        },
        examples = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Example]
            }
          }
        },
        headers = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Header]
            }
          }
        },
        securitySchemes = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.SecurityScheme]
            }
          }
        },
        callbacks = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Callback]
            }
          }
        },
        parameters = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Parameter]
            }
          }
        },
        requestBodies = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.RequestBody]
            }
          }
        },
        links = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Link]
            }
          }
        }
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    XML {
      type = "object"
      properties = {
        prefix = string(),
        attribute = boolean(default(false)),
        wrapped = boolean(default(false)),
        name = string(),
        namespace = string(format("uri"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    ExampleXORExamples {
      not = {
        required = ["example", "examples"]
      }
      description = "Example and examples are mutually exclusive"
    }
    HTTPSecurityScheme {
      type = "object"
      required = ["scheme", "type"]
      properties = {
        type = string(enum("http")),
        scheme = string(),
        bearerFormat = string(),
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      oneOf = [{
        properties = {
          scheme = string(pattern("^[Bb][Ee][Aa][Rr][Ee][Rr]$"))
        },
        description = "Bearer"
      }, {
        properties = {
          scheme = {
            not = string(pattern("^[Bb][Ee][Aa][Rr][Ee][Rr]$"))
          }
        },
        not = {
          required = ["bearerFormat"]
        },
        description = "Non Bearer"
      }]
    }
    OAuth2SecurityScheme {
      properties = {
        type = string(enum("oauth2")),
        flows = definitions.OAuthFlows,
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["type", "flows"]
    }
    Info {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["title", "version"]
      properties = {
        contact = definitions.Contact,
        license = definitions.License,
        version = string(),
        title = string(),
        description = string(),
        termsOfService = string(format("uri-reference"))
      }
    }
    PathItem {
      properties = {
        trace = definitions.Operation,
        servers = array(definitions.Server),
        delete = definitions.Operation,
        "$ref" = string(),
        description = string(),
        post = definitions.Operation,
        options = definitions.Operation,
        parameters = array({
          oneOf = [definitions.Parameter, definitions.Reference]
        }, uniqueItems(true)),
        summary = string(),
        get = definitions.Operation,
        put = definitions.Operation,
        head = definitions.Operation,
        patch = definitions.Operation
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
    }
    APIKeySecurityScheme {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["type", "name", "in"]
      properties = {
        in = string(enum("header", "query", "cookie")),
        description = string(),
        type = string(enum("apiKey")),
        name = string()
      }
    }
    PasswordOAuthFlow {
      type = "object"
      required = ["tokenUrl", "scopes"]
      properties = {
        refreshUrl = string(format("uri-reference")),
        scopes = map(string()),
        tokenUrl = string(format("uri-reference"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Example {
      type = "object"
      properties = {
        externalValue = string(format("uri-reference")),
        summary = string(),
        description = string(),
        value = {}
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Tag {
      required = ["name"]
      properties = {
        name = string(),
        description = string(),
        externalDocs = definitions.ExternalDocumentation
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
    }
    SchemaXORContent {
      oneOf = [{
        required = ["schema"]
      }, {
        allOf = [{
          not = {
            required = ["style"]
          }
        }, {
          not = {
            required = ["explode"]
          }
        }, {
          not = {
            required = ["allowReserved"]
          }
        }, {
          not = {
            required = ["example"]
          }
        }, {
          not = {
            required = ["examples"]
          }
        }],
        description = "Some properties are not allowed if content is present",
        required = ["content"]
      }]
      not = {
        required = ["schema", "content"]
      }
      description = "Schema and content are mutually exclusive, at least one is required"
    }
    RequestBody {
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["content"]
      properties = {
        description = string(),
        content = map(definitions.MediaType),
        required = boolean(default(false))
      }
      additionalProperties = false
    }
    ImplicitOAuthFlow {
      required = ["authorizationUrl", "scopes"]
      properties = {
        scopes = map(string()),
        authorizationUrl = string(format("uri-reference")),
        refreshUrl = string(format("uri-reference"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
    }
    Link {
      not = {
        required = ["operationId", "operationRef"],
        description = "Operation Id and Operation Ref are mutually exclusive"
      }
      type = "object"
      properties = {
        server = definitions.Server,
        operationId = string(),
        operationRef = string(),
        parameters = map({}),
        description = string(),
        requestBody = {}
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Encoding {
      type = "object"
      properties = {
        contentType = string(),
        headers = map({
          oneOf = [definitions.Header, definitions.Reference]
        }),
        style = string(enum("form", "spaceDelimited", "pipeDelimited", "deepObject")),
        explode = boolean(),
        allowReserved = boolean(default(false))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Reference {
      patternProperties = {
        "^\\$ref$" = string(format("uri-reference"))
      }
      type = "object"
      required = ["$ref"]
    }
    License {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["name"]
      properties = {
        name = string(),
        url = string(format("uri-reference"))
      }
    }
    ServerVariable {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["default"]
      properties = {
        description = string(),
        enum = array(string()),
        default = string()
      }
    }
    Parameter {
      properties = {
        description = string(),
        deprecated = boolean(default(false)),
        allowEmptyValue = boolean(default(false)),
        style = string(),
        name = string(),
        required = boolean(default(false)),
        allowReserved = boolean(default(false)),
        in = string(),
        explode = boolean(),
        examples = map({
          oneOf = [definitions.Example, definitions.Reference]
        }),
        schema = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        content = {
          type = "object",
          maxProperties = 1,
          minProperties = 1,
          additionalProperties = definitions.MediaType
        },
        example = {}
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      allOf = [definitions.ExampleXORExamples, definitions.SchemaXORContent]
      oneOf = [definitions.PathParameter, definitions.QueryParameter, definitions.HeaderParameter, definitions.CookieParameter]
      type = "object"
      required = ["name", "in"]
    }
    HeaderParameter {
      properties = {
        in = {
          enum = ["header"]
        },
        style = {
          default = "simple",
          enum = ["simple"]
        }
      }
      description = "Parameter in header"
    }
    AuthorizationCodeOAuthFlow {
      type = "object"
      required = ["authorizationUrl", "tokenUrl", "scopes"]
      properties = {
        authorizationUrl = string(format("uri-reference")),
        tokenUrl = string(format("uri-reference")),
        refreshUrl = string(format("uri-reference")),
        scopes = map(string())
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
  }
