// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Server to store articles and search similar articles.",
    "title": "Article similarity",
    "version": "1.0.0"
  },
  "basePath": "/",
  "paths": {
    "/articles": {
      "get": {
        "summary": "Get unique articles.",
        "responses": {
          "200": {
            "description": "OK.",
            "schema": {
              "type": "object",
              "required": [
                "articles"
              ],
              "properties": {
                "articles": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Article"
                  }
                }
              }
            }
          },
          "500": {
            "$ref": "#/responses/ServerError"
          }
        }
      },
      "post": {
        "summary": "Add an article.",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "content"
              ],
              "properties": {
                "content": {
                  "description": "Article content",
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Article added.",
            "schema": {
              "$ref": "#/definitions/Article"
            }
          },
          "400": {
            "$ref": "#/responses/InvalidArgument"
          },
          "500": {
            "$ref": "#/responses/ServerError"
          }
        }
      }
    },
    "/articles/{id}": {
      "get": {
        "summary": "Get article by id.",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "Article id",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ArticleOk"
          },
          "400": {
            "$ref": "#/responses/InvalidArgument"
          },
          "404": {
            "description": "Article not found.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "$ref": "#/responses/ServerError"
          }
        }
      }
    },
    "/duplicate_groups": {
      "get": {
        "summary": "Get duplicate groups ids.",
        "responses": {
          "200": {
            "description": "OK.",
            "schema": {
              "type": "object",
              "required": [
                "duplicate_groups"
              ],
              "properties": {
                "duplicate_groups": {
                  "type": "array",
                  "items": {
                    "type": "array",
                    "items": {
                      "$ref": "#/definitions/ArticleId"
                    }
                  }
                }
              }
            }
          },
          "500": {
            "$ref": "#/responses/ServerError"
          }
        }
      }
    }
  },
  "definitions": {
    "Article": {
      "type": "object",
      "required": [
        "id",
        "content",
        "duplicate_article_ids"
      ],
      "properties": {
        "content": {
          "description": "Article content",
          "type": "string"
        },
        "duplicate_article_ids": {
          "description": "Duplicated articles",
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "id": {
          "$ref": "#/definitions/ArticleId"
        }
      }
    },
    "ArticleId": {
      "description": "Article id",
      "type": "integer",
      "format": "int64"
    },
    "Error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "description": "Error code for machine parsing",
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "description": "Human-readable error message",
          "type": "string"
        }
      }
    }
  },
  "responses": {
    "ArticleOk": {
      "description": "OK",
      "schema": {
        "$ref": "#/definitions/Article"
      }
    },
    "InvalidArgument": {
      "description": "Invalid arguments",
      "schema": {
        "$ref": "#/definitions/Error"
      }
    },
    "ServerError": {
      "description": "Internal server error",
      "schema": {
        "$ref": "#/definitions/Error"
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Server to store articles and search similar articles.",
    "title": "Article similarity",
    "version": "1.0.0"
  },
  "basePath": "/",
  "paths": {
    "/articles": {
      "get": {
        "summary": "Get unique articles.",
        "responses": {
          "200": {
            "description": "OK.",
            "schema": {
              "type": "object",
              "required": [
                "articles"
              ],
              "properties": {
                "articles": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Article"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "summary": "Add an article.",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "content"
              ],
              "properties": {
                "content": {
                  "description": "Article content",
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Article added.",
            "schema": {
              "$ref": "#/definitions/Article"
            }
          },
          "400": {
            "description": "Invalid arguments",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/articles/{id}": {
      "get": {
        "summary": "Get article by id.",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "Article id",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/Article"
            }
          },
          "400": {
            "description": "Invalid arguments",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "404": {
            "description": "Article not found.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/duplicate_groups": {
      "get": {
        "summary": "Get duplicate groups ids.",
        "responses": {
          "200": {
            "description": "OK.",
            "schema": {
              "type": "object",
              "required": [
                "duplicate_groups"
              ],
              "properties": {
                "duplicate_groups": {
                  "type": "array",
                  "items": {
                    "type": "array",
                    "items": {
                      "$ref": "#/definitions/ArticleId"
                    }
                  }
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Article": {
      "type": "object",
      "required": [
        "id",
        "content",
        "duplicate_article_ids"
      ],
      "properties": {
        "content": {
          "description": "Article content",
          "type": "string"
        },
        "duplicate_article_ids": {
          "description": "Duplicated articles",
          "type": "array",
          "items": {
            "type": "integer"
          }
        },
        "id": {
          "$ref": "#/definitions/ArticleId"
        }
      }
    },
    "ArticleId": {
      "description": "Article id",
      "type": "integer",
      "format": "int64"
    },
    "Error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "description": "Error code for machine parsing",
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "description": "Human-readable error message",
          "type": "string"
        }
      }
    }
  },
  "responses": {
    "ArticleOk": {
      "description": "OK",
      "schema": {
        "$ref": "#/definitions/Article"
      }
    },
    "InvalidArgument": {
      "description": "Invalid arguments",
      "schema": {
        "$ref": "#/definitions/Error"
      }
    },
    "ServerError": {
      "description": "Internal server error",
      "schema": {
        "$ref": "#/definitions/Error"
      }
    }
  }
}`))
}
