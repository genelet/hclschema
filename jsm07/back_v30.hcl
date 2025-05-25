
  patternProperties = {
    "^x-" = {}
  }
  description = "The description of OpenAPI v3.0.x documents, as defined by https://spec.openapis.org/oas/v3.0.3"
  type = "object"
  required = ["openapi", "info", "paths"]
  properties = {
    info = definitions.Info,
    externalDocs = definitions.ExternalDocumentation,
    servers = array(definitions.Server),
    security = array(definitions.SecurityRequirement),
    tags = array(definitions.Tag, uniqueItems(true)),
    paths = definitions.Paths,
    components = definitions.Components,
    openapi = string(pattern("^3\\.0\\.\\d(-.+)?$"))
  }
  additionalProperties = false
  schema = "http://json-schema.org/draft-04/schema#"
  id = "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
  definitions {
    Discriminator = object({
      propertyName = string(),
      mapping = map(string())
    }, required(["propertyName"]))
    SecurityRequirement = map(array(string()))
    Encoding {
      type = "object"
      properties = {
        explode = boolean(),
        allowReserved = boolean(default(false)),
        contentType = string(),
        headers = map({
          oneOf = [definitions.Header, definitions.Reference]
        }),
        style = string(enum("form", "spaceDelimited", "pipeDelimited", "deepObject"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Reference {
      type = "object"
      required = ["$ref"]
      patternProperties = {
        "^\\$ref$" = string(format("uri-reference"))
      }
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
    Tag {
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["name"]
      properties = {
        name = string(),
        description = string(),
        externalDocs = definitions.ExternalDocumentation
      }
      additionalProperties = false
    }
    ImplicitOAuthFlow {
      type = "object"
      required = ["authorizationUrl", "scopes"]
      properties = {
        refreshUrl = string(format("uri-reference")),
        scopes = map(string()),
        authorizationUrl = string(format("uri-reference"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    ServerVariable {
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["default"]
      properties = {
        enum = array(string()),
        default = string(),
        description = string()
      }
      additionalProperties = false
    }
    Parameter {
      oneOf = [definitions.PathParameter, definitions.QueryParameter, definitions.HeaderParameter, definitions.CookieParameter]
      type = "object"
      required = ["name", "in"]
      properties = {
        examples = map({
          oneOf = [definitions.Example, definitions.Reference]
        }),
        required = boolean(default(false)),
        explode = boolean(),
        deprecated = boolean(default(false)),
        allowEmptyValue = boolean(default(false)),
        style = string(),
        name = string(),
        description = string(),
        allowReserved = boolean(default(false)),
        in = string(),
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
    SecurityScheme {
      oneOf = [definitions.APIKeySecurityScheme, definitions.HTTPSecurityScheme, definitions.OAuth2SecurityScheme, definitions.OpenIdConnectSecurityScheme]
    }
    PathParameter {
      required = ["required"]
      properties = {
        style = {
          enum = ["matrix", "label", "simple"],
          default = "simple"
        },
        required = {
          enum = [true]
        },
        in = {
          enum = ["path"]
        }
      }
      description = "Parameter in path"
    }
    OAuthFlows {
      properties = {
        authorizationCode = definitions.AuthorizationCodeOAuthFlow,
        implicit = definitions.ImplicitOAuthFlow,
        password = definitions.PasswordOAuthFlow,
        clientCredentials = definitions.ClientCredentialsFlow
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
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
    Example {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      properties = {
        summary = string(),
        description = string(),
        externalValue = string(format("uri-reference")),
        value = {}
      }
    }
    License {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["name"]
      properties = {
        url = string(format("uri-reference")),
        name = string()
      }
    }
    AuthorizationCodeOAuthFlow {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["authorizationUrl", "tokenUrl", "scopes"]
      properties = {
        refreshUrl = string(format("uri-reference")),
        scopes = map(string()),
        authorizationUrl = string(format("uri-reference")),
        tokenUrl = string(format("uri-reference"))
      }
    }
    Header {
      type = "object"
      properties = {
        examples = map({
          oneOf = [definitions.Example, definitions.Reference]
        }),
        allowEmptyValue = boolean(default(false)),
        style = string(default("null"), enum("simple")),
        allowReserved = boolean(default(false)),
        description = string(),
        explode = boolean(),
        required = boolean(default(false)),
        deprecated = boolean(default(false)),
        example = {},
        schema = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        content = {
          additionalProperties = definitions.MediaType,
          type = "object",
          maxProperties = 1,
          minProperties = 1
        }
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      allOf = [definitions.ExampleXORExamples, definitions.SchemaXORContent]
    }
    CookieParameter {
      properties = {
        in = {
          enum = ["cookie"]
        },
        style = {
          default = "form",
          enum = ["form"]
        }
      }
      description = "Parameter in cookie"
    }
    Server {
      properties = {
        description = string(),
        variables = map(definitions.ServerVariable),
        url = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["url"]
    }
    ExampleXORExamples {
      not = {
        required = ["example", "examples"]
      }
      description = "Example and examples are mutually exclusive"
    }
    Link {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      not = {
        required = ["operationId", "operationRef"],
        description = "Operation Id and Operation Ref are mutually exclusive"
      }
      type = "object"
      properties = {
        operationId = string(),
        operationRef = string(),
        parameters = map({}),
        description = string(),
        server = definitions.Server,
        requestBody = {}
      }
    }
    Contact {
      properties = {
        name = string(),
        url = string(format("uri-reference")),
        email = string(format("email"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
    }
    ExternalDocumentation {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      required = ["url"]
      properties = {
        url = string(format("uri-reference")),
        description = string()
      }
    }
    OpenIdConnectSecurityScheme {
      type = "object"
      required = ["type", "openIdConnectUrl"]
      properties = {
        type = string(enum("openIdConnect")),
        openIdConnectUrl = string(format("uri-reference")),
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Callback {
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      additionalProperties = definitions.PathItem
    }
    Components {
      type = "object"
      properties = {
        responses = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Response]
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
        requestBodies = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.RequestBody]
            }
          }
        },
        schemas = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Schema, definitions.Reference]
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
        parameters = {
          type = "object",
          patternProperties = {
            "^[a-zA-Z0-9\\.\\-_]+$" = {
              oneOf = [definitions.Reference, definitions.Parameter]
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
    APIKeySecurityScheme {
      required = ["type", "name", "in"]
      properties = {
        in = string(enum("header", "query", "cookie")),
        description = string(),
        type = string(enum("apiKey")),
        name = string()
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
    HTTPSecurityScheme {
      type = "object"
      required = ["scheme", "type"]
      properties = {
        scheme = string(),
        bearerFormat = string(),
        description = string(),
        type = string(enum("http"))
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
    Schema {
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      type = "object"
      properties = {
        exclusiveMaximum = boolean(default(false)),
        nullable = boolean(default(false)),
        exclusiveMinimum = boolean(default(false)),
        xml = definitions.XML,
        minItems = integer(default(0), minimum(0)),
        discriminator = definitions.Discriminator,
        enum = array({}, minItems(1), uniqueItems(false)),
        allOf = array({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        writeOnly = boolean(default(false)),
        minLength = integer(default(0), minimum(0)),
        minProperties = integer(default(0), minimum(0)),
        required = array(string(), minItems(1), uniqueItems(true)),
        multipleOf = number(minimum(0), exclusiveMinimum(true)),
        properties = map({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        title = string(),
        pattern = string(format("regex")),
        oneOf = array({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        anyOf = array({
          oneOf = [definitions.Schema, definitions.Reference]
        }),
        minimum = number(),
        maxProperties = integer(minimum(0)),
        maxLength = integer(minimum(0)),
        format = string(),
        maxItems = integer(minimum(0)),
        externalDocs = definitions.ExternalDocumentation,
        type = string(enum("array", "boolean", "integer", "number", "object", "string")),
        readOnly = boolean(default(false)),
        maximum = number(),
        deprecated = boolean(default(false)),
        description = string(),
        uniqueItems = boolean(default(false)),
        default = {},
        example = {},
        not = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        items = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        additionalProperties = {
          default = true,
          oneOf = [definitions.Schema, definitions.Reference, boolean()]
        }
      }
    }
    MediaType {
      type = "object"
      properties = {
        encoding = map(definitions.Encoding),
        examples = map({
          oneOf = [definitions.Example, definitions.Reference]
        }),
        schema = {
          oneOf = [definitions.Schema, definitions.Reference]
        },
        example = {}
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
      allOf = [definitions.ExampleXORExamples]
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
    ClientCredentialsFlow {
      type = "object"
      required = ["tokenUrl", "scopes"]
      properties = {
        scopes = map(string()),
        tokenUrl = string(format("uri-reference")),
        refreshUrl = string(format("uri-reference"))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
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
      patternProperties = {
        "^x-" = {},
        "^[1-5](?:\\d{2}|XX)$" = {
          oneOf = [definitions.Response, definitions.Reference]
        }
      }
      type = "object"
      minProperties = 1
      properties = {
        default = {
          oneOf = [definitions.Response, definitions.Reference]
        }
      }
      additionalProperties = false
    }
    OAuth2SecurityScheme {
      type = "object"
      required = ["type", "flows"]
      properties = {
        type = string(enum("oauth2")),
        flows = definitions.OAuthFlows,
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Info {
      type = "object"
      required = ["title", "version"]
      properties = {
        description = string(),
        termsOfService = string(format("uri-reference")),
        contact = definitions.Contact,
        license = definitions.License,
        version = string(),
        title = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    Operation {
      type = "object"
      required = ["responses"]
      properties = {
        description = string(),
        security = array(definitions.SecurityRequirement),
        parameters = array({
          oneOf = [definitions.Parameter, definitions.Reference]
        }, uniqueItems(true)),
        tags = array(string()),
        callbacks = map({
          oneOf = [definitions.Callback, definitions.Reference]
        }),
        summary = string(),
        deprecated = boolean(default(false)),
        servers = array(definitions.Server),
        externalDocs = definitions.ExternalDocumentation,
        operationId = string(),
        responses = definitions.Responses,
        requestBody = {
          oneOf = [definitions.RequestBody, definitions.Reference]
        }
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
    PathItem {
      type = "object"
      properties = {
        get = definitions.Operation,
        patch = definitions.Operation,
        parameters = array({
          oneOf = [definitions.Parameter, definitions.Reference]
        }, uniqueItems(true)),
        delete = definitions.Operation,
        servers = array(definitions.Server),
        head = definitions.Operation,
        "$ref" = string(),
        description = string(),
        post = definitions.Operation,
        put = definitions.Operation,
        trace = definitions.Operation,
        summary = string(),
        options = definitions.Operation
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
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
    RequestBody {
      type = "object"
      required = ["content"]
      properties = {
        description = string(),
        content = map(definitions.MediaType),
        required = boolean(default(false))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = {}
      }
    }
  }
