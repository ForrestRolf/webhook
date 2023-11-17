# Webhook

**Webhook** is a lightweight tool written in Go that can be used to create HTTP hook services, which you can use to execute configured commands.

For example, if you are using Harbor, you can use webhook to configure a hook service to automatically update Docker and restart the service on the server when the image is updated.

# Features

* Beautiful management UI
* Supports multiple data formats: JSON, XML, Form
* Supports multiple actions: Shell, http，email，slack，SMS
* Hook supports authentication
* Detailed execution logs
* Debug mode
* Backup / Restore hooks
* Script templates

> NOTE: SMS is not yet tested

# Dependencies
* MongoDB

# Usage

```shell
./webhook -uri mongodb://user:pass@localhost:27017
```

Open http://localhost:9000 in your browser.

## Enable Basic Auth

```shell
./webhook -uri mongodb://user:pass@localhost:27017 -u foo -p bar
```

# Using Variables

## Variables Scope

* Shell
* Http payload
* Email subject and body
* Slack message

## Using variables in shell

```shell
echo $WEBHOOK_FOO
```

**Note:** You need to set the name of the ENV in the Argument tab before you can use it.

## Using variables in http and other actions

For example, in a http action, configure the payload as follows:

```json
{
	"Var": "${foo}"
}
```

# Use Cases

webhook is a very flexible tool that can be used in a variety of scenarios. For example, you can use webhook to:

* Automatically update Docker and restart the service when the image is updated
* Send email notifications when a user registers
* Automatically send SMS notifications when an order is created

webhook is a tool worth trying.
