# track-wpp-cloud

The idea of the project is to create a Chat Bot for Whatsapp that is able to track package codes from the Brazilian postal system. This way, the user can send the code he wants to track to the BOT and will automatically start receiving the tracking updates.

The step-by-step of the project is as follows: the user registers a package by sending the command:

/track CODE (in Portuguese)

The WhatsApp Cloud API is configured to send the webhooks to an API Gateway that calls a lambda function. This function will interpret the user's message and make a call to an external package tracking API (https://linketrack.com).
After this a message is sent to the user indicating the current status of the order. In addition, the code, the user number, and the date of the last order event are saved in DynamoDB to send recurring notifications.

An event is triggered every 5 minutes by CloudWatch events, which triggers the execution of another lambda function, this function fetches all items from DynamoDB and makes another call to the external tracking API. If there are updates (according to the date of the last event), the notification is sent to the users that have already searched for this package.

Some handling was also implemented to give feedback to the user in case they type an invalid message. As well as help for the basic /help and /track commands.

There are still some areas for improvement, such as identifying the packages that have already been delivered and will no longer receive updates, removing the need to always scan DynamoDB, and other specific improvements.
