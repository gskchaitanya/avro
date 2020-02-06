package example

import (
	avroPkg "github.com/heetch/cue-schema/avro"
)

avro :: avroPkg
avro :: Metadata :: heetchmeta?: {
	commentary: string
	status: string | *"active"
	partitions: int | *1
	topickey: string
}

DomainName :: "{{.DomainName}}"
EventName ::  "{{.EventName}}"
Version :: "{{.Version}}"

cloudEvent: avro.Schema & {
	type:      "record"
	name:      "com.heetch.\(DomainName).\(EventName)"
	heetchmeta: {
		commentary: "This Schema describes version \(Version) of the event \(EventName) from the domain \(DomainName)."
		topickey:   "\(DomainName).\(EventName).\(Version)"
	}
	fields: [{
		name: "Metadata"
		type: {
			type: "record"
			name: "Metadata"
			fields: [{
				name: "CloudEvent"
				type: {
					type: "record"
					name: "CloudEvent"
					fields: [{
						name: "id"
						type: "string"
					}, {
						name: "source"
						type: "string"
					}, {
						name: "specversion"
						type: "string"
					}, {
						name: "time"
						type: {
							type:        "long"
							logicalType: "timestamp-micros"
						}
					}]
				}
			}]
		}
	}]
}

fixedType: avro.Schema & {
	type: "fixed"
	name: "five"
	size: 5
}

locicalType: avro.Schema & {
	type:        "long"
	logicalType: "timestamp-micros"
}

anyType: avro.Schema & [
		"bytes",
		"null",
		"boolean",
		"double",
		"string",
		{
		type: "map"
		values: ["anyType"]
	},
	{
		type: "array"
		items: ["anyType"]
	},
]

heetchEventSchema: avro.Schema & {
	type: "record"
	name: "CloudEvent"
	doc:  "Avro Event Format for Heetch Events"
	fields: [
		{
			name: "attribute"
			type: {
				type: "map"
				values: ["null", "boolean", "int", "string", "bytes"]
			}
		},
		{
			name: "data"
			type: anyType
		},
	]
}

exampleHeetchEvent: avro.Protocol & {
	namespace: "com.heetch"
	protocol:  "heetch"
	types: [heetchEventSchema]
}

exampleProtocol: avro.Protocol & {
	namespace: "com.acme"
	protocol:  "HelloWorld"
	doc:       "Protocol Greetings"
	types: [{
		name: "Greeting"
		type: "record"
		fields: [{name: "message", type: "string"}]
	}, {
		name: "Curse"
		type: "record"
		fields: [{name: "message", type: "string"}]
	}]
	messages: hello: {
		doc: "Say hello."
		request: [{name: "greeting", type: "Greeting"}]
		response: "Greeting"
		errors: ["Curse"]
	}
}

exampleTODOApp : avro.Schema & {
	type:      "record"
	name:      "User"
	namespace: "com.example.avro"
	doc: """
		This is a user record in a fictitious to-do-list
		management app. It supports arbitrary grouping and
		nesting of items, and allows you to add items by email
		or by tweeting.

		Note this app doesn't actually exist. The schema is just a demo for [Avrodoc](https://github.com/ept/avrodoc)!
		"""
	fields: [{
		name: "id"
		doc:  "System-assigned numeric user ID. Cannot be changed by the user."
		type: "int"
	}, {
		name: "username"
		doc:  "The username chosen by the user. Can be changed by the user."
		type: "string"
	}, {
		name: "passwordHash"
		doc:  "The user's password, hashed using [scrypt](http://www.tarsnap.com/scrypt.html)."
		type: "string"
	}, {
		name: "signupDate"
		doc:  "Timestamp (milliseconds since epoch) when the user signed up"
		type: "long"
	}, {
		name: "emailAddresses"
		doc:  "All email addresses on the user's account"
		type: {
			type: "array"
			items: {
				type: "record"
				name: "EmailAddress"
				doc:  "Stores details about an email address that a user has associated with their account."
				fields: [{
					name: "address"
					doc:  "The email address, e.g. `foo@example.com`"
					type: "string"
				}, {
					name:    "verified"
					doc:     "true if the user has clicked the link in a confirmation email to this address."
					type:    "boolean"
					default: false
				}, {
					name: "dateAdded"
					doc:  "Timestamp (milliseconds since epoch) when the email address was added to the account."
					type: "long"
				}, {
					name: "dateBounced"
					doc:  "Timestamp (milliseconds since epoch) when an email sent to this address last bounced. Reset to null when the address no longer bounces."
					type: ["null", "long"]
				}]
			}
		}
	}, {
		name: "twitterAccounts"
		doc:  "All Twitter accounts that the user has OAuthed"
		type: {
			type: "array"
			items: {
				type: "record"
				name: "TwitterAccount"
				doc:  "Stores access credentials for one Twitter account, as granted to us by the user by OAuth."
				fields: [{
					name: "status"
					doc:  "Indicator of whether this authorization is currently active, or has been revoked"
					type: {
						type: "enum"
						name: "OAuthStatus"
						doc: """
							* `PENDING`: the user has started authorizing, but not yet finished
							* `ACTIVE`: the token should work
							* `DENIED`: the user declined the authorization
							* `EXPIRED`: the token used to work, but now it doesn't
							* `REVOKED`: the user has explicitly revoked the token
							"""
						symbols: ["PENDING", "ACTIVE", "DENIED", "EXPIRED", "REVOKED"]
					}
				}, {
					name: "userId"
					doc:  "Twitter's numeric ID for this user"
					type: "long"
				}, {
					name: "screenName"
					doc:  "The twitter username for this account (can be changed by the user)"
					type: "string"
				}, {
					name: "oauthToken"
					doc:  "The OAuth token for this Twitter account"
					type: "string"
				}, {
					name: "oauthTokenSecret"
					doc:  "The OAuth secret, used for signing requests on behalf of this Twitter account. `null` whilst the OAuth flow is not yet complete."
					type: ["null", "string"]
				}, {
					name: "dateAuthorized"
					doc:  "Timestamp (milliseconds since epoch) when the user last authorized this Twitter account"
					type: "long"
				}]
			}
		}
	}, {
		name: "toDoItems"
		doc:  "The top-level items in the user's to-do list"
		type: {
			type: "array"
			items: {
				type: "record"
				name: "ToDoItem"
				doc:  "A record is one node in a To-Do item tree (every record can contain nested sub-records)."
				fields: [{
					name: "status"
					doc:  "User-selected state for this item (e.g. whether or not it is marked as done)"
					type: {
						type: "enum"
						name: "ToDoStatus"
						doc: """
							* `HIDDEN`: not currently visible, e.g. because it becomes actionable in future
							* `ACTIONABLE`: appears in the current to-do list
							* `DONE`: marked as done, but still appears in the list
							* `ARCHIVED`: marked as done and no longer visible
							* `DELETED`: not done and removed from list (preserved for undo purposes)
							"""
						symbols: ["HIDDEN", "ACTIONABLE", "DONE", "ARCHIVED", "DELETED"]
					}
				}, {
					name: "title"
					doc:  "One-line summary of the item"
					type: "string"
				}, {
					name: "description"
					doc:  "Detailed description (may contain HTML markup)"
					type: ["null", "string"]
				}, {
					name: "snoozeDate"
					doc:  "Timestamp (milliseconds since epoch) at which the item should go from `HIDDEN` to `ACTIONABLE` status"
					type: ["null", "long"]
				}, {
					name: "subItems"
					doc:  "List of children of this to-do tree node"
					type: {
						type:  "array"
						items: "ToDoItem"
					}
				}]
			}
		}
	}]
}
