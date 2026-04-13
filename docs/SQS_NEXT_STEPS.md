# SQS Next Steps

This is the recommended shape for adding SQS listening while keeping responsibilities clean and scalable.

## Infrastructure

### `internal/infrastructure/sqs/client.go`

- builds the AWS SQS SDK client

### `internal/infrastructure/sqs/repo/listener.go`

- long-polls SQS
- receives messages
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
- handler owns use case

That is the cleanest version of this design.
