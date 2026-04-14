# SQS Next Steps

This is the recommended shape for adding SQS listening while keeping responsibilities clean and scalable.

## What We Implemented

Today we moved from the design stage into a working first implementation.

- `internal/infrastructure/sqs/client.go` builds the AWS SQS SDK client.
- `internal/infrastructure/sqs/repo/listener.go` now owns the full SQS message lifecycle.
- `internal/application/events/sqsDispatcher.go` routes incoming events by `DetailType`.
- `internal/application/user/mapper.go` maps `appEvent.PublishingEvent` into typed application events like `LoginEvent`.
- `internal/application/user/handlers/handle_user_login.go` handles the actual business reaction for login events.
- `cmd/server/main.go` wires the listener as a background worker, separate from HTTP routing.

This means the app now has two inbound entry points:

- HTTP requests through the router
- SQS messages through the background listener

## Current Runtime Flow

1. The application publishes a login event to EventBridge.
2. Terraform routes the matching event into `dev-calculator-main-queue`.
3. `SqsListener.Listen(ctx)` long-polls the queue.
4. The listener receives an SQS message.
5. The listener decodes the EventBridge-delivered SQS body into `appEvent.PublishingEvent`.
6. The listener calls `SqsDispatcher.Dispatch(event)`.
7. The dispatcher checks `event.DetailType`.
8. For login events, the dispatcher uses `UserMapper.FromPublishingEventToLoginEvent(...)`.
9. The mapper unmarshals `event.Detail` JSON into `user/events/LoginEvent`.
10. The dispatcher calls `UserLoginHandler.HandleLogin(...)`.
11. If handling succeeds, the listener deletes the SQS message.
12. If handling fails, the listener leaves the message in the queue so SQS retry/DLQ behavior can take over.

## What The Listener Does Now

`internal/infrastructure/sqs/repo/listener.go` currently does the following:

- accepts `ctx context.Context`
- runs an infinite loop until the context is canceled
- uses `ReceiveMessage` with long polling
- processes one message at a time
- decodes the incoming message body
- dispatches the event to the application layer
- calls a private `deleteMessage(...)` helper only after successful dispatch
- does not delete failed messages

This keeps transport concerns in infrastructure and business concerns in application handlers.

## EventBridge And SQS Message Shape

One important detail we clarified today:

- `PublishingEvent.Detail` is stored as a JSON string
- when consuming from SQS, the listener first unwraps the outer EventBridge envelope
- then the application mapper uses `json.Unmarshal([]byte(event.Detail), &loginEvent)` to decode the typed event payload

So the decoding responsibilities are split like this:

- listener: decode SQS/EventBridge transport envelope
- mapper: decode typed application event payload

## Main Wiring

The SQS listener is not part of the HTTP router.

Instead, `cmd/server/main.go` now:

- builds the SQS client
- builds the SQS dispatcher
- builds the SQS listener
- starts `Listen(ctx)` in a goroutine
- starts the HTTP server normally

This is the correct pattern because the listener is a background worker, not an HTTP endpoint.

## Current Boundary Decision

The boundary we landed on is:

- listener owns transport
- dispatcher owns routing
- mapper owns payload conversion
- handler owns use case

That separation is now implemented and working end to end.

## Infrastructure

### `internal/infrastructure/sqs/client.go`

- builds the AWS SQS SDK client

### `internal/infrastructure/sqs/repo/listener.go`

- long-polls SQS
- receives messages
- decodes the EventBridge-delivered SQS body into `appEvent.PublishingEvent`
- forwards events to the application dispatcher
- deletes messages only after successful handling
- sends failed ones back by not deleting, letting redrive policy work

### Optional: `internal/infrastructure/sqs/repo/message_decoder.go`

- extracts the EventBridge envelope from the SQS body if needed

## Application

### `internal/application/events/handlers/...`

- one handler per event/use case, like `HandleUserLogin`

### `internal/application/events/dispatcher.go`

- routes by `source` + `detailType`
- unmarshals `detail` into the right app event
- calls the right handler

In the current implementation this is `internal/application/events/sqsDispatcher.go`.

## Flow

1. SQS listener polls messages
2. listener parses outer SQS message
3. listener extracts EventBridge payload
4. listener passes a generic envelope to an application dispatcher
5. dispatcher picks the correct use case handler
6. on success, listener deletes message
7. on failure, listener logs and leaves message for retry/DLQ

## Best-Practice Boundaries

### `SqsClient`

- only AWS transport
- `ReceiveMessage`, `DeleteMessage`, maybe `ChangeMessageVisibility`

### `SqsListener`

- polling loop and message lifecycle
- decoding the SQS/EventBridge transport envelope
- no business logic

### `Application handlers`

- business reaction to events
- no SQS/AWS SDK knowledge

### `Dispatcher`

- small orchestration layer between incoming event contract and handlers

## Handling Different Events

The listener should stay generic and should not know all event types.

Instead, introduce a dispatcher interface conceptually like:

- `Handle(ctx, envelope) error`

Then the dispatcher can do:

- if `detailType == login`, decode login event, call login handler
- if `detailType == calculation`, decode calculation event, call calculation handler

That scales much better than `switch` logic inside the SQS repo itself.

## Practical Best Practices

- use long polling with `WaitTimeSeconds`
- process messages one by one first; batch later if needed
- only delete after successful handler execution
- log message id, source, detail type, receive count
- treat unknown event types explicitly
- keep handlers idempotent because SQS can deliver more than once
- use `context.Context` through the whole chain
- let DLQs handle poison messages instead of custom retry loops in app code

## Recommended Starting Point

- `internal/infrastructure/sqs/client.go`
- `internal/infrastructure/sqs/repo/listener.go`
- `internal/application/appEvent/dispatcher.go` or similar
- `internal/application/user/handlers/handle_user_login.go`

That keeps one listener, one dispatcher, and many handlers as the app grows.

## Key Rule

- listener owns transport
- dispatcher owns routing
- mapper owns payload conversion
- handler owns use case

That is the cleanest version of this design.
