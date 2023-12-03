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
        "/api/v1/workout": {
            "post": {
                "description": "This endpoint starts a new workout session for a player with the given details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Start a new workout session",
                "operationId": "start-workout",
                "parameters": [
                    {
                        "description": "Details of the workout to start",
                        "name": "workout",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/httphandler.StartWorkout"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully started workout session"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/distance": {
            "get": {
                "description": "This endpoint retrieves the distance covered in a workout session either by workout ID or by player ID within a date range.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Get distance covered in a workout",
                "operationId": "get-distance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the workout session",
                        "name": "workoutID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "ID of the player",
                        "name": "playerID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date for the range (RFC3339 format)",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date for the range (RFC3339 format)",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully retrieved distance"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/escapes": {
            "get": {
                "description": "This endpoint retrieves the number of escapes made either by workout ID or between dates for a player.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Get escapes made in a workout",
                "operationId": "get-escapes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Workout ID to fetch escapes",
                        "name": "workoutID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Player ID to fetch escapes between dates",
                        "name": "playerID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date for fetching escapes",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date for fetching escapes",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully retrieved escape count"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/fights": {
            "get": {
                "description": "This endpoint retrieves the number of fights fought either by workout ID or between dates for a player.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Get fights fought in a workout",
                "operationId": "get-fights",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Workout ID to fetch fights",
                        "name": "workoutID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Player ID to fetch fights between dates",
                        "name": "playerID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date for fetching fights",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date for fetching fights",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully retrieved fight count"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/shelters": {
            "get": {
                "description": "This endpoint retrieves the number of shelters taken either by workout ID or between dates for a player.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Get shelters taken in a workout",
                "operationId": "get-shelters",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Workout ID to fetch shelters",
                        "name": "workoutID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Player ID to fetch shelters between dates",
                        "name": "playerID",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date for fetching shelters",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date for fetching shelters",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully retrieved shelter count"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/{workoutId}": {
            "put": {
                "description": "This endpoint stops the workout session for a player based on the provided workout ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Stop an ongoing workout session",
                "operationId": "stop-workout",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the workout session to stop",
                        "name": "workoutId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Successfully stopped workout session"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            },
            "delete": {
                "description": "This endpoint deletes a workout session based on the provided workout ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Delete a workout session",
                "operationId": "delete-workout",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the workout session to delete",
                        "name": "workoutId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully deleted workout session"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    },
                    "404": {
                        "description": "Workout session not found"
                    }
                }
            }
        },
        "/api/v1/workout/{workoutId}/options": {
            "get": {
                "description": "This endpoint retrieves the available options for a workout session based on the workout ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Get workout session options",
                "operationId": "get-workout-options",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the workout session",
                        "name": "workoutId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved workout options"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/{workoutId}/options/start": {
            "post": {
                "description": "This endpoint starts a specific option for an ongoing workout session based on the workout ID and option details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Start a workout option",
                "operationId": "start-workout-option",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the workout session",
                        "name": "workoutId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Details of the workout option to start",
                        "name": "option",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/httphandler.StartWorkoutOption"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully started workout option"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        },
        "/api/v1/workout/{workoutId}/options/stop": {
            "patch": {
                "description": "This endpoint stops a specific option of an ongoing workout session based on the workout ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workout"
                ],
                "summary": "Stop a workout option",
                "operationId": "stop-workout-option",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the workout session",
                        "name": "workoutId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully stopped workout option"
                    },
                    "400": {
                        "description": "Bad Request with error details"
                    }
                }
            }
        }
    },
    "definitions": {
        "httphandler.StartWorkout": {
            "type": "object",
            "properties": {
                "hardcore_mode": {
                    "description": "HardCore Mode of User",
                    "type": "boolean"
                },
                "hrm_connected": {
                    "description": "Whether HRM is connected or not",
                    "type": "boolean"
                },
                "hrm_id": {
                    "description": "If HRM is connected then HRM ID otherwise garbage",
                    "type": "string"
                },
                "player_id": {
                    "description": "PlayerID of the player starting the workout session",
                    "type": "string"
                },
                "trail_id": {
                    "description": "TrailID chosen by the Player",
                    "type": "string"
                }
            }
        },
        "httphandler.StartWorkoutOption": {
            "type": "object",
            "properties": {
                "option": {
                    "description": "WorkoutID for which the workout option is to be stopped",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
