// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/in-memory": {
            "get": {
                "description": "fetches cache entries in-memory",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "In-Memory"
                ],
                "summary": "fetches a cache entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Tuple"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "creates a cache entry in-memory",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "In-Memory"
                ],
                "summary": "creates a cache entry",
                "parameters": [
                    {
                        "description": "Key Value data",
                        "name": "tuple",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Tuple"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Tuple"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/mdb": {
            "post": {
                "description": "fetch records by filtering",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MongoDB"
                ],
                "summary": "Fetch records",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.MDBRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.MDBResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.MDBResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.MDBRequest": {
            "type": "object",
            "properties": {
                "endDate": {
                    "type": "string"
                },
                "maxCount": {
                    "type": "integer"
                },
                "minCount": {
                    "type": "integer"
                },
                "startDate": {
                    "type": "string"
                }
            }
        },
        "model.MDBResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "msg": {
                    "type": "string"
                },
                "records": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Record"
                    }
                }
            }
        },
        "model.Record": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "key": {
                    "type": "string"
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "model.Tuple": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Rest API Documentation",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
