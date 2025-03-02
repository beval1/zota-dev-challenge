// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/deposit": {
            "post": {
                "description": "handle deposit",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "deposit"
                ],
                "summary": "deposit example",
                "parameters": [
                    {
                        "description": "Deposit ClientRequest",
                        "name": "depositRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/deposit.ClientRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Deposit Successful",
                        "schema": {
                            "$ref": "#/definitions/zota.DepositResponse"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "description": "handle status check",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status check"
                ],
                "summary": "status check example",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "orderId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Merchant Order ID",
                        "name": "merchantOrderId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Status Check successful",
                        "schema": {
                            "$ref": "#/definitions/zota.StatusResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "deposit.ClientRequest": {
            "type": "object",
            "properties": {
                "customerAddress": {
                    "type": "string"
                },
                "customerBankCode": {
                    "type": "string"
                },
                "customerCity": {
                    "type": "string"
                },
                "customerCountryCode": {
                    "type": "string"
                },
                "customerEmail": {
                    "type": "string"
                },
                "customerFirstName": {
                    "type": "string"
                },
                "customerLastName": {
                    "type": "string"
                },
                "customerPhone": {
                    "type": "string"
                },
                "customerState": {
                    "type": "string"
                },
                "customerZipCode": {
                    "type": "string"
                },
                "orderAmount": {
                    "type": "string"
                },
                "orderCurrency": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "zota.Data": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "customParam": {
                    "type": "string"
                },
                "customerEmail": {
                    "type": "string"
                },
                "endpointID": {
                    "type": "string"
                },
                "errorMessage": {
                    "type": "string"
                },
                "extraData": {
                    "$ref": "#/definitions/zota.ExtraData"
                },
                "merchantOrderID": {
                    "type": "string"
                },
                "orderID": {
                    "type": "string"
                },
                "processorTransactionID": {
                    "type": "string"
                },
                "request": {
                    "$ref": "#/definitions/zota.StatusRequest"
                },
                "status": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "zota.DepositResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {
                    "$ref": "#/definitions/zota.DepositResponseData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "zota.DepositResponseData": {
            "type": "object",
            "properties": {
                "depositUrl": {
                    "type": "string"
                },
                "merchantOrderID": {
                    "type": "string"
                },
                "orderID": {
                    "type": "string"
                }
            }
        },
        "zota.ExtraData": {
            "type": "object",
            "properties": {
                "amountChanged": {
                    "type": "boolean"
                },
                "amountManipulated": {
                    "type": "boolean"
                },
                "amountRounded": {
                    "type": "boolean"
                },
                "dcc": {
                    "type": "boolean"
                },
                "originalAmount": {
                    "type": "string"
                },
                "paymentMethod": {
                    "type": "string"
                },
                "selectedBankCode": {
                    "type": "string"
                },
                "selectedBankName": {
                    "type": "string"
                }
            }
        },
        "zota.StatusRequest": {
            "type": "object",
            "properties": {
                "merchantID": {
                    "type": "string"
                },
                "merchantOrderID": {
                    "type": "string"
                },
                "orderID": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "zota.StatusResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {
                    "$ref": "#/definitions/zota.Data"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Merchant Server",
	Description:      "This is a simple merchant's server implementing Zota payment gateway.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
