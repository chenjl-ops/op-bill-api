// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package api

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/bill/v1/get_all_bill_data": {
            "get": {
                "description": "查询决算全量数据",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Compute API"
                ],
                "summary": "Select billing data",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "get bill of share or source all data",
                        "name": "isShare",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/bill/v1/get_all_prediction_data": {
            "get": {
                "description": "查询预测全量数据",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Compute API"
                ],
                "summary": "Select prediction all data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "get all prediction data",
                        "name": "date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/bill/v1/get_bill_data": {
            "get": {
                "description": "查询决算数据",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Compute API"
                ],
                "summary": "Select billing data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "get bill of month",
                        "name": "month",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "get bill of share or source",
                        "name": "isShare",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/bill/v1/get_prediction_data": {
            "get": {
                "description": "查询预测数据",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Compute API"
                ],
                "summary": "Select prediction data of someday",
                "parameters": [
                    {
                        "type": "string",
                        "description": "get prediction of date, default today",
                        "name": "date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/billing/v1/create_table": {
            "get": {
                "description": "创建损益和资金口径账单数据表，对应账单状态表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Create Table",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/billing/v1/get_month_data": {
            "get": {
                "description": "插入账单数据 资金和损益口径",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Select Month Data",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/billing/v1/init_tex_data": {
            "get": {
                "description": "初始化折扣率数据",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Create Table",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/billing/v1/insert_bill_data": {
            "get": {
                "description": "插入账单数据 资金和损益口径",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Insert Data",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/billing/v1/tex": {
            "get": {
                "description": "获取资源折扣率",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Get Data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "select name of tex",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "更新资源折扣率",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Get Data",
                "parameters": [
                    {
                        "description": "new SourceBillTex",
                        "name": "tex",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/billing.SourceBillTex"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "新增资源折扣率",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Get Data",
                "parameters": [
                    {
                        "description": "new SourceBillTex",
                        "name": "tex",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/billing.SourceBillTex"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "删除资源折扣率",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Billing API"
                ],
                "summary": "Get Data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "delete name of tex",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/prediction/v1/create_table": {
            "get": {
                "description": "create BillData table",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Prediction API"
                ],
                "summary": "Create BillData Table",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/prediction/v1/insert_baidu_bill_data": {
            "get": {
                "description": "Insert Bill Data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Prediction API"
                ],
                "summary": "Insert Bill Data",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.ResponseData"
                        },
                        "headers": {
                            "config.ResponseData": {
                                "type": "object"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "billing.SourceBillTex": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "tex": {
                    "type": "number"
                }
            }
        },
        "config.ResponseData": {
            "type": "object",
            "properties": {
                "columns": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "additionalProperties": {
                            "type": "string"
                        }
                    }
                },
                "data": {
                    "type": "object",
                    "additionalProperties": true
                },
                "error": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Op-bill-api API",
	Description: "This is op-bill-api api server.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
