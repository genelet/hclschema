
  schema = "http://json-schema.org/draft-04/schema#"
  id = "http://openapis.org/v3/schema.json#"
  patternProperties = {
    "^x-" = definitions.specificationExtension
  }
  description = "This is the root document object of the OpenAPI document."
  type = "object"
  required = ["openapi", "info", "paths"]
  properties = {
    security = array(definitions.securityRequirement, uniqueItems(true)),
    tags = array(definitions.tag, uniqueItems(true)),
    externalDocs = definitions.externalDocs,
    openapi = string(),
    info = definitions.info,
    servers = array(definitions.server, uniqueItems(true)),
    paths = definitions.paths,
    components = definitions.components
  }
  additionalProperties = false
  title = "A JSON Schema for OpenAPI 3.0."
  definitions {
    parametersOrReferences = map(definitions.parameterOrReference)
    expression = map(true)
    responsesOrReferences = map(definitions.responseOrReference)
    mediaTypes = map(definitions.mediaType)
    object = map(true)
    linksOrReferences = map(definitions.linkOrReference)
    headersOrReferences = map(definitions.headerOrReference)
    schemasOrReferences = map(definitions.schemaOrReference)
    strings = map(string())
    examplesOrReferences = map(definitions.exampleOrReference)
    callbacksOrReferences = map(definitions.callbackOrReference)
    encodings = map(definitions.encoding)
    securitySchemesOrReferences = map(definitions.securitySchemeOrReference)
    requestBodiesOrReferences = map(definitions.requestBodyOrReference)
    serverVariables = map(definitions.serverVariable)
    specificationExtension {
      description = "Any property starting with x- is valid."
      oneOf = [{
        type = "null"
      }, number(), boolean(), string(), object(), array()]
    }
    info {
      type = "object"
      required = ["title", "version"]
      properties = {
        contact = definitions.contact,
        license = definitions.license,
        version = string(),
        summary = string(),
        title = string(),
        description = string(),
        termsOfService = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "The object provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience."
    }
    pathItem {
      description = "Describes the operations available on a single path. A Path Item MAY be empty, due to ACL constraints. The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available."
      type = "object"
      properties = {
        description = string(),
        get = definitions.operation,
        put = definitions.operation,
        delete = definitions.operation,
        servers = array(definitions.server, uniqueItems(true)),
        "$ref" = string(),
        post = definitions.operation,
        parameters = array(definitions.parameterOrReference, uniqueItems(true)),
        summary = string(),
        options = definitions.operation,
        trace = definitions.operation,
        head = definitions.operation,
        patch = definitions.operation
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
    }
    encoding {
      type = "object"
      properties = {
        contentType = string(),
        headers = definitions.headersOrReferences,
        style = string(),
        explode = boolean(),
        allowReserved = boolean()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "A single encoding definition applied to a single schema property."
    }
    callback {
      description = "A map of possible out-of band callbacks related to the parent operation. Each value in the map is a Path Item Object that describes a set of requests that may be initiated by the API provider and the expected responses. The key value used to identify the callback object is an expression, evaluated at runtime, that identifies a URL to use for the callback operation."
      type = "object"
      additionalProperties = false
      patternProperties = {
        "^" = definitions.pathItem,
        "^x-" = definitions.specificationExtension
      }
    }
    operation {
      type = "object"
      required = ["responses"]
      properties = {
        responses = definitions.responses,
        callbacks = definitions.callbacksOrReferences,
        security = array(definitions.securityRequirement, uniqueItems(true)),
        summary = string(),
        parameters = array(definitions.parameterOrReference, uniqueItems(true)),
        servers = array(definitions.server, uniqueItems(true)),
        tags = array(string(), uniqueItems(true)),
        description = string(),
        operationId = string(),
        deprecated = boolean(),
        externalDocs = definitions.externalDocs,
        requestBody = definitions.requestBodyOrReference
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Describes a single API operation on a path."
    }
    mediaType {
      properties = {
        encoding = definitions.encodings,
        schema = definitions.schemaOrReference,
        example = definitions.any,
        examples = definitions.examplesOrReferences
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Each Media Type Object provides schema and examples for the media type identified by its key."
      type = "object"
    }
    link {
      description = "The `Link object` represents a possible design-time link for a response. The presence of a link does not guarantee the caller's ability to successfully invoke it, rather it provides a known relationship and traversal mechanism between responses and other operations.  Unlike _dynamic_ links (i.e. links provided **in** the response payload), the OAS linking mechanism does not require link information in the runtime response.  For computing links, and providing instructions to execute them, a runtime expression is used for accessing values in an operation and using them as parameters while invoking the linked operation."
      type = "object"
      properties = {
        operationRef = string(),
        operationId = string(),
        parameters = definitions.anyOrExpression,
        requestBody = definitions.anyOrExpression,
        description = string(),
        server = definitions.server
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
    }
    discriminator {
      description = "When request bodies or response payloads may be one of a number of different schemas, a `discriminator` object can be used to aid in serialization, deserialization, and validation.  The discriminator is a specific object in a schema which is used to inform the consumer of the specification of an alternative schema based on the value associated with it.  When using the discriminator, _inline_ schemas will not be considered."
      type = "object"
      required = ["propertyName"]
      properties = {
        mapping = definitions.strings,
        propertyName = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
    }
    securitySchemeOrReference {
      oneOf = [definitions.securityScheme, definitions.reference]
    }
    defaultType {
      oneOf = [{
        type = "null"
      }, array(), object(), number(), boolean(), string()]
    }
    server {
      type = "object"
      required = ["url"]
      properties = {
        url = string(),
        description = string(),
        variables = definitions.serverVariables
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "An object representing a Server."
    }
    serverVariable {
      type = "object"
      required = ["default"]
      properties = {
        default = string(),
        description = string(),
        enum = array(string(), uniqueItems(true))
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "An object representing a Server Variable for server URL template substitution."
    }
    callbackOrReference {
      oneOf = [definitions.callback, definitions.reference]
    }
    requestBody {
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Describes a single request body."
      type = "object"
      required = ["content"]
      properties = {
        required = boolean(),
        description = string(),
        content = definitions.mediaTypes
      }
    }
    tag {
      description = "Adds metadata to a single tag that is used by the Operation Object. It is not mandatory to have a Tag Object per tag defined in the Operation Object instances."
      type = "object"
      required = ["name"]
      properties = {
        externalDocs = definitions.externalDocs,
        name = string(),
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
    }
    anyOrExpression {
      oneOf = [definitions.any, definitions.expression]
    }
    parameter {
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Describes a single operation parameter.  A unique parameter is defined by a combination of a name and location."
      type = "object"
      required = ["name", "in"]
      properties = {
        examples = definitions.examplesOrReferences,
        style = string(),
        explode = boolean(),
        allowReserved = boolean(),
        schema = definitions.schemaOrReference,
        content = definitions.mediaTypes,
        required = boolean(),
        example = definitions.any,
        name = string(),
        in = string(),
        deprecated = boolean(),
        description = string(),
        allowEmptyValue = boolean()
      }
    }
    contact {
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Contact information for the exposed API."
      type = "object"
      properties = {
        email = string(format("email")),
        name = string(),
        url = string(format("uri"))
      }
      additionalProperties = false
    }
    paths {
      patternProperties = {
        "^/" = definitions.pathItem,
        "^x-" = definitions.specificationExtension
      }
      description = "Holds the relative paths to the individual endpoints and their operations. The path is appended to the URL from the `Server Object` in order to construct the full URL.  The Paths MAY be empty, due to ACL constraints."
      type = "object"
      additionalProperties = false
    }
    linkOrReference {
      oneOf = [definitions.link, definitions.reference]
    }
    response {
      type = "object"
      required = ["description"]
      properties = {
        headers = definitions.headersOrReferences,
        content = definitions.mediaTypes,
        links = definitions.linksOrReferences,
        description = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Describes a single response from an API Operation, including design-time, static  `links` to operations based on the response."
    }
    xml {
      type = "object"
      properties = {
        wrapped = boolean(),
        name = string(),
        namespace = string(),
        prefix = string(),
        attribute = boolean()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "A metadata object that allows for more fine-tuned XML model definitions.  When using arrays, XML element names are *not* inferred (for singular/plural forms) and the `name` property SHOULD be used to add that information. See examples for expected behavior."
    }
    oauthFlow {
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Configuration details for a supported OAuth Flow"
      type = "object"
      properties = {
        authorizationUrl = string(),
        tokenUrl = string(),
        refreshUrl = string(),
        scopes = definitions.strings
      }
    }
    requestBodyOrReference {
      oneOf = [definitions.requestBody, definitions.reference]
    }
    license {
      properties = {
        name = string(),
        url = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "License information for the exposed API."
      type = "object"
      required = ["name"]
    }
    schema {
      type = "object"
      properties = {
        multipleOf = properties.multipleOf,
        uniqueItems = properties.uniqueItems,
        title = properties.title,
        minimum = properties.minimum,
        oneOf = array(definitions.schemaOrReference, minItems(1)),
        minItems = properties.minItems,
        format = string(),
        minProperties = properties.minProperties,
        maximum = properties.maximum,
        maxProperties = properties.maxProperties,
        readOnly = boolean(),
        externalDocs = definitions.externalDocs,
        enum = properties.enum,
        not = definitions.schema,
        example = definitions.any,
        deprecated = boolean(),
        nullable = boolean(),
        description = string(),
        maxLength = properties.maxLength,
        pattern = properties.pattern,
        writeOnly = boolean(),
        type = string(),
        xml = definitions.xml,
        default = definitions.defaultType,
        exclusiveMinimum = properties.exclusiveMinimum,
        maxItems = properties.maxItems,
        required = properties.required,
        anyOf = array(definitions.schemaOrReference, minItems(1)),
        minLength = properties.minLength,
        exclusiveMaximum = properties.exclusiveMaximum,
        properties = map(definitions.schemaOrReference),
        allOf = array(definitions.schemaOrReference, minItems(1)),
        discriminator = definitions.discriminator,
        additionalProperties = {
          oneOf = [definitions.schemaOrReference, boolean()]
        },
        items = {
          anyOf = [definitions.schemaOrReference, array(definitions.schemaOrReference, minItems(1))]
        }
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "The Schema Object allows the definition of input and output data types. These types can be objects, but also primitives and arrays. This object is an extended subset of the JSON Schema Specification Wright Draft 00.  For more information about the properties, see JSON Schema Core and JSON Schema Validation. Unless stated otherwise, the property definitions follow the JSON Schema."
    }
    components {
      type = "object"
      properties = {
        examples = definitions.examplesOrReferences,
        requestBodies = definitions.requestBodiesOrReferences,
        headers = definitions.headersOrReferences,
        schemas = definitions.schemasOrReferences,
        links = definitions.linksOrReferences,
        callbacks = definitions.callbacksOrReferences,
        parameters = definitions.parametersOrReferences,
        securitySchemes = definitions.securitySchemesOrReferences,
        responses = definitions.responsesOrReferences
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Holds a set of reusable objects for different aspects of the OAS. All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object."
    }
    externalDocs {
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Allows referencing an external resource for extended documentation."
      type = "object"
      required = ["url"]
      properties = {
        description = string(),
        url = string()
      }
    }
    reference {
      type = "object"
      required = ["$ref"]
      properties = {
        description = string(),
        "$ref" = string(),
        summary = string()
      }
      additionalProperties = false
      description = "A simple object to allow referencing other components in the specification, internally and externally.  The Reference Object is defined by JSON Reference and follows the same structure, behavior and rules.   For this specification, reference resolution is accomplished as defined by the JSON Reference specification and not by the JSON Schema specification."
    }
    oauthFlows {
      description = "Allows configuration of the supported OAuth Flows."
      type = "object"
      properties = {
        implicit = definitions.oauthFlow,
        password = definitions.oauthFlow,
        clientCredentials = definitions.oauthFlow,
        authorizationCode = definitions.oauthFlow
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
    }
    responseOrReference {
      oneOf = [definitions.response, definitions.reference]
    }
    schemaOrReference {
      oneOf = [definitions.schema, definitions.reference]
    }
    securityScheme {
      required = ["type"]
      properties = {
        in = string(),
        scheme = string(),
        bearerFormat = string(),
        flows = definitions.oauthFlows,
        openIdConnectUrl = string(),
        type = string(),
        description = string(),
        name = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "Defines a security scheme that can be used by the operations. Supported schemes are HTTP authentication, an API key (either as a header, a cookie parameter or as a query parameter), mutual TLS (use of a client certificate), OAuth2's common flows (implicit, password, application and access code) as defined in RFC6749, and OpenID Connect.   Please note that currently (2019) the implicit flow is about to be deprecated OAuth 2.0 Security Best Current Practice. Recommended for most use case is Authorization Code Grant flow with PKCE."
      type = "object"
    }
    exampleOrReference {
      oneOf = [definitions.example, definitions.reference]
    }
    headerOrReference {
      oneOf = [definitions.header, definitions.reference]
    }
    any {
      additionalProperties = true
    }
    responses {
      description = "A container for the expected responses of an operation. The container maps a HTTP response code to the expected response.  The documentation is not necessarily expected to cover all possible HTTP response codes because they may not be known in advance. However, documentation is expected to cover a successful operation response and any known errors.  The `default` MAY be used as a default response object for all HTTP codes  that are not covered individually by the specification.  The `Responses Object` MUST contain at least one response code, and it  SHOULD be the response for a successful operation call."
      type = "object"
      properties = {
        default = definitions.responseOrReference
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension,
        "^([0-9X]{3})$" = definitions.responseOrReference
      }
    }
    example {
      type = "object"
      properties = {
        summary = string(),
        description = string(),
        value = definitions.any,
        externalValue = string()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = ""
    }
    header {
      type = "object"
      properties = {
        allowReserved = boolean(),
        example = definitions.any,
        content = definitions.mediaTypes,
        description = string(),
        required = boolean(),
        deprecated = boolean(),
        examples = definitions.examplesOrReferences,
        schema = definitions.schemaOrReference,
        allowEmptyValue = boolean(),
        style = string(),
        explode = boolean()
      }
      additionalProperties = false
      patternProperties = {
        "^x-" = definitions.specificationExtension
      }
      description = "The Header Object follows the structure of the Parameter Object with the following changes:  1. `name` MUST NOT be specified, it is given in the corresponding `headers` map. 1. `in` MUST NOT be specified, it is implicitly in `header`. 1. All traits that are affected by the location MUST be applicable to a location of `header` (for example, `style`)."
    }
    securityRequirement {
      description = "Lists the required security schemes to execute this operation. The name used for each property MUST correspond to a security scheme declared in the Security Schemes under the Components Object.  Security Requirement Objects that contain multiple schemes require that all schemes MUST be satisfied for a request to be authorized. This enables support for scenarios where multiple query parameters or HTTP headers are required to convey security information.  When a list of Security Requirement Objects is defined on the OpenAPI Object or Operation Object, only one of the Security Requirement Objects in the list needs to be satisfied to authorize the request."
      type = "object"
      additionalProperties = array(string(), uniqueItems(true))
    }
    parameterOrReference {
      oneOf = [definitions.parameter, definitions.reference]
    }
  }
