package config

const deploymentConfigSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"title": "Deployment configuration",
	"description": "URLs and instructions for how to install and run the application.",
    "type": "object",
	"definitions": {
		"TargetPlatformsArray": {
			"type": "array",
			"items": {
				"type": "string",
				"pattern": "^(((windows|darwin|linux|\\{\\{\\.OS\\}\\})(-(386|amd64|\\{\\{\\.Arch\\}\\}))?)|(386|amd64|\\{\\{\\.Arch\\}\\}))$"
			},
			"uniqueItems": true
		},
		"URL": {
			"type": "string",
			"pattern": "^(https?|file)://.*$"
		}
    },
	"properties": {
		"Timestamp": {
			"type": "string",
			"pattern": "^([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})|(<TIMESTAMP>)$"
		},
		"LauncherUpdate": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"BundleInfoURL": {
						"$ref": "#/definitions/URL"
					},
					"BundleURL": {
						"$ref": "#/definitions/URL"
					},
					"TargetPlatforms": {
						"$ref": "#/definitions/TargetPlatformsArray"
					}
				},
                "required": [ "BundleInfoURL" ]
			}
		},
		"Bundles": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"BundleInfoURL": {
						"$ref": "#/definitions/URL"
					},
					"BundleURL": {
						"$ref": "#/definitions/URL"
					},
					"TargetPlatforms": {
						"$ref": "#/definitions/TargetPlatformsArray"
					},
					"LocalDirectory": {
						"type": "string",
						"minLength": 1
					},
					"Tags": {
						"type": "array",
						"items": {
							"type": "string"
						}
					}
				},
				"required": [ "BundleInfoURL", "LocalDirectory" ]
			},
			"minItems": 1,
			"uniqueItems": true
		},
		"Execution": {
			"type": "object",
			"properties": {
				"Commands": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"Name": {
								"type": "string",
								"minLength": 1
							},
							"WorkingDirectoryBundleName": {
								"type": "string"
							},
							"Arguments": {
								"type": "array",
								"items": {
									"type": "string"
								}
							},
							"Env": {
								"type": "object",
								"additionalProperties": {
									"oneOf": [ { "type": "string" }, { "type": "null" } ]
								}
							},
							"TargetPlatforms": {
								"$ref": "#/definitions/TargetPlatformsArray"
							}
						},
						"required": [ "Name" ]
					}
				},
				"LingerTimeMilliseconds": {
					"type": "integer",
					"minimum": 0
				}
			}
		}
	},
	"required": [ "Timestamp", "Bundles", "Execution" ]
}`
