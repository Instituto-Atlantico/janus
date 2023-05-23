// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/agents": {
            "get": {
                "description": "Get a list with of all agents",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "summary": "Get details of all agents",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/agents/{IpAddress}": {
            "delete": {
                "description": "Delete the agent corresponding to the input IP address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "summary": "Delete the agent identified by the given IP address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "IP address of the device to be deleted",
                        "name": "IpAddress",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/provision": {
            "post": {
                "description": "Create a new provision with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "provision"
                ],
                "summary": "Create a new provision",
                "parameters": [
                    {
                        "description": "Create provision",
                        "name": "provision",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller_handlers.ProvisionBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller_handlers.ProvisionBody"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller_handlers.ProvisionBody": {
            "type": "object",
            "properties": {
                "brokerIp": {
                    "description": "127.0.0.1",
                    "type": "string",
                    "example": "192.168.0.2"
                },
                "brokerPassword": {
                    "description": "brokerPass",
                    "type": "string",
                    "example": "admin"
                },
                "brokerUsername": {
                    "description": "brokerUser:deviceID",
                    "type": "string",
                    "example": "admin:a1998e"
                },
                "deviceHostName": {
                    "description": "user@ip",
                    "type": "string",
                    "example": "rasp@192.168.0.5"
                },
                "permissions": {
                    "description": "[\"temperature\", \"humidity\"]",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "temperature",
                        " humidity"
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "janus-issuer",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}